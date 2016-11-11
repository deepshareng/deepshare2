package aggregate

import (
	"reflect"
	"sort"
	"time"

	"fmt"

	"strings"

	"strconv"

	in "github.com/MISingularity/deepshare2/pkg/instrumentation"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const TimestampForm = "20140101 00:00:00 +0000 UTC"
const fakedAppIDForCompatibleColl = "all"

var (
	intervalRefreshToMongoSec int64 = 60 // Interval to refresh aggregate results to mongoDB (in seconds)
)

const (
	redisKeyFmt                       = "deepstats_buffer:%s:%s:%d"             // deepstats_buffer:<nsq_identifier>:<appid>:<interval_start>
	redisKeyBufferedAppIDsFmt         = "deepstats_buffered_appids:%s"          // all appids that buffered in redis, waiting for refresh to mongo
	redisKeyBufferedIntervalStartsFmt = "deepstats_buffered_interval_starts:%s" // all interval starts that buffered in redis, waiting for refresh to mongo
)

// AggregateService serves counter aggregation functionalities.
type AggregateService interface {
	// Store user data into a temporary storage
	//TODO Insert to different buffers for different apps
	Insert(appID string, aggregate CounterEvent) error

	// Aggregate returns the aggregated results on all records.
	// For example, if some counter records were updated in buffer,
	// and commited to persistent storage,
	// #Insert Counter(chan1, Install, 2006-01-02T10:00:00Z00:00, 1)
	// #Insert Counter(chan1, Install, 2006-01-02T12:00:00Z00:00, 1)
	// #Insert Counter(chan1, Install, 2006-01-01T12:00:00Z00:00, 1)
	// #Aggregate(chan1) => AggregateResult: [
	//    event: Install => [
	//      AggregateCount(2006-01-02T00:00:00Z00:00, 2)
	//      AggregateCount(2006-01-01T00:00:00Z00:00, 1)
	//    ]
	//  ]
	//
	// If eventFilters is empty, then all events will be returned.
	// Otherwise, it returns only the events in eventFilters.

	// QueryDuration query aggregate event in a period, and specified paritcular appid/channel/need event.
	QueryDuration(appid string, channel string, eventFilters []string, start time.Time, granularity time.Duration, limit int, os string) ([]*AggregateResult, error)

	// Aggregate data from temporarry storage to persistent one.
	Aggregate(appID string) error

	// StartRefreshLoop starts a loop to refresh aggregate results from redis to mongo periodically
	StartRefreshLoop()

	// On the basis of user concrete aggregate granularity,
	// ConvertTimeToGranularity will regulate the aggregation
	// data time
	// E.g.
	// If user want to aggregate data by day,
	// For any time in a day, like
	// 		20060102T10:00:00Z00:00
	// It might aggreate into this time
	// 		20060102T00:00:00Z00:00
	ConvertTimeToGranularity(time.Time) time.Time
}

type generalAggregateService struct {
	nsqIdentifier  string
	buffer         map[string]*CounterEvent
	redisDB        storage.SimpleKV
	colls          map[string]*mgo.Collection
	collNamePrefix string
	mgoSession     *mgo.Session
	mgoDBName      string
}

func getEventMapKey(aggregate *CounterEvent) string {
	return aggregate.Timestamp.String() + "#" + aggregate.AppID + "#" + aggregate.Channel + "#" + aggregate.Event + "#" + aggregate.Os
}

func retrieveCounterEventFromMapKey(key string) *CounterEvent {
	fields := strings.Split(key, "#")
	if len(fields) != 5 {
		log.Panic(fmt.Sprintf("invalid key: %s", key))
	}
	timestamp, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", fields[0])
	if err != nil {
		panic(err)
	}
	return &CounterEvent{
		Timestamp: timestamp,
		AppID:     fields[1],
		Channel:   fields[2],
		Event:     fields[3],
		Os:        fields[4],
	}
}

func (m *generalAggregateService) queryDurationWithEnd(appid string, channel string, eventFilters []string, start time.Time, end time.Time, os string) ([]*AggregateResultWithTimestamp, error) {
	results := []bson.M{}
	operations := buildMongoPipeOpsOfQueryAggregatedData(appid, channel, eventFilters, start, end, os)
	coll := GetColl(m.mgoSession, m.mgoDBName, m.colls, m.collNamePrefix, appid)
	err := coll.Pipe(operations).All(&results)
	if err != nil {
		return nil, err
	}
	//If results not found, the data may be in the old collection which contains data of all appIDs
	if len(results) == 0 {
		coll := GetColl(m.mgoSession, m.mgoDBName, m.colls, m.collNamePrefix, fakedAppIDForCompatibleColl)
		err := coll.Pipe(operations).All(&results)
		if err != nil {
			return nil, err
		}
	}
	aggrs := make([]*AggregateResultWithTimestamp, len(results))
	for i, result := range results {
		resCounts := result["counts"].([]interface{})
		counts := make([]*AggregateCountWithTimestamp, len(resCounts))
		for j, resCount := range resCounts {
			rc := resCount.(bson.M)
			ts := reflect.ValueOf(rc["timestamp"]).Interface().(time.Time)
			ts = ts.UTC()
			if err != nil {
				return nil, err
			}
			counts[j] = &AggregateCountWithTimestamp{
				Timestamp: ts,
				Count:     rc["count"].(int),
			}
		}
		sort.Sort(ByTimestamp(counts))
		aggrs[i] = &AggregateResultWithTimestamp{
			Event:  result["_id"].(string),
			Counts: counts,
		}
	}
	// We only record successful flow to reflect the actual workload

	return aggrs, err
}

func (m *generalAggregateService) QueryDuration(appid string, channel string, eventFilters []string, start time.Time, granulairty time.Duration, limit int, os string) ([]*AggregateResult, error) {
	log.Infof("[AGGREGATE SERVICE] Request aggregate by appid=%s, channel=%s, eventFilter=%v, start=%s, gran=%s, limit=%d, os=%s", appid, channel, eventFilters, start, granulairty, limit, os)
	starttime := time.Now()
	end := start.Add(granulairty * time.Duration(limit))
	results, err := m.queryDurationWithEnd(appid, channel, eventFilters, start, end, os)
	if err != nil {
		return []*AggregateResult{}, err
	}
	aggrs := make([]*AggregateResult, len(results))
	for i, result := range results {
		counts := make([]*AggregateCount, limit)
		lenCounts := len(result.Counts)
		interval := start
		cur := 0
		for i := 0; i < limit; i++ {
			interval = interval.Add(granulairty)
			counts[i] = &AggregateCount{}
			for cur < lenCounts && interval.After(result.Counts[cur].Timestamp) {
				counts[i].Count += result.Counts[cur].Count
				cur++
			}
		}

		aggrs[i] = &AggregateResult{
			Event:  result.Event,
			Counts: counts,
		}
	}
	// We only record successful flow to reflect the actual workload
	in.PromCounter.AggregateDuration(starttime)
	log.Debugf("[--Result][AGGREGATE SERVICE][QueryDuration] Aggregate Results=%s, Err=%v", AggregateResultsToString(aggrs), err)
	return aggrs, err
}

func (m *generalAggregateService) Aggregate(appID string) error {
	log.Debug("[AGGREGATE SERVICE] Aggregate current buffer events")
	startTime := time.Now()
	for k, aggEvent := range m.buffer {
		if aggEvent.AppID != appID {
			continue
		}
		///---------------------- store results to redis
		now := time.Now().Unix()
		intervalStart := now
		if intervalRefreshToMongoSec != 0 {
			intervalStart = (now / intervalRefreshToMongoSec) * intervalRefreshToMongoSec
		}

		key := fmt.Sprintf(redisKeyFmt, m.nsqIdentifier, appID, intervalStart)
		field := getEventMapKey(aggEvent)
		if m.redisDB == nil {
			log.Panic("redisDB is nil")
		}
		if err := m.redisDB.HIncrBy([]byte(key), field, aggEvent.Count); err != nil {
			log.Error("Redis HIncrBy failed, err:", err)
			return err
		}
		if err := m.redisDB.SAdd([]byte(fmt.Sprintf(redisKeyBufferedAppIDsFmt, m.nsqIdentifier)), appID); err != nil {
			log.Error("Redis SAdd appid failed, err:", err)
			return err
		}
		if err := m.redisDB.SAdd([]byte(fmt.Sprintf(redisKeyBufferedIntervalStartsFmt, m.nsqIdentifier)), strconv.Itoa(int(intervalStart))); err != nil {
			log.Error("Redis SAdd intervalStart failed, err:", err)
			return err
		}
		///---------------------- store results to redis END

		delete(m.buffer, k)
	}
	// We only record successful flow to reflect the actual workload
	in.PromCounter.AggregateDuration(startTime)

	// return aggrs, nil
	return nil
}

func (m *generalAggregateService) refreshToMongo() error {
	if m.redisDB == nil {
		log.Panic("redisDB is nil")
	}
	now := time.Now().Unix()
	keyIntervals := fmt.Sprintf(redisKeyBufferedIntervalStartsFmt, m.nsqIdentifier)
	keyAppids := fmt.Sprintf(redisKeyBufferedAppIDsFmt, m.nsqIdentifier)
	intervalStarts, err := m.redisDB.SMembers([]byte(keyIntervals))
	if err != nil {
		log.Panic(err)
	}
	appIDs, err := m.redisDB.SMembers([]byte(keyAppids))
	if err != nil {
		log.Panic(err)
	}
	log.Debug("refreshToMongo -------------------------", m.nsqIdentifier, now, intervalStarts, appIDs)
	for _, intervelStartStr := range intervalStarts {
		intervalStart, err := strconv.Atoi(intervelStartStr)
		if err != nil {
			panic(err)
		}
		// only process the expired data in redis
		if now-int64(intervalStart) < intervalRefreshToMongoSec*2 {
			log.Debug("not expire", now, intervalStart, intervalRefreshToMongoSec)
			continue
		}
		for _, appID := range appIDs {
			key := fmt.Sprintf(redisKeyFmt, m.nsqIdentifier, appID, intervalStart)
			kvs, err := m.redisDB.HGetAll([]byte(key))
			if err != nil {
				log.Error(err)
				continue
			}
			for k, v := range kvs {
				aggEvent := retrieveCounterEventFromMapKey(k)
				n, err := strconv.Atoi(v)
				if err != nil {
					panic(err)
				}
				aggEvent.Count = n
				///---------------------- upsert to mongo
				if appID != fakedAppIDForCompatibleColl && appID != aggEvent.AppID {
					continue
				}

				session := m.mgoSession.Copy()
				coll := GetColl(session, m.mgoDBName, m.colls, m.collNamePrefix, appID)

				if _, err := coll.Upsert(
					bson.M{"appid": aggEvent.AppID, "channel": aggEvent.Channel, "event": aggEvent.Event, "timestamp": aggEvent.Timestamp, "os": aggEvent.Os},
					bson.M{
						"$inc": bson.M{"count": aggEvent.Count},
					},
				); err != nil {
					return err
				}
				///----------------------- upsert to mongo END
			}
			m.redisDB.Delete([]byte(key)) //TODO err
			m.redisDB.SRem([]byte(keyIntervals), strconv.Itoa(int(intervalStart)))
			//m.redisCli.SRem([]byte(keyAppids), appID)
		}
	}
	return nil
}

func (m *generalAggregateService) StartRefreshLoop() {
	go func() {
		log.Debug("=====================StartRefreshLoop")
		for true {
			time.Sleep(time.Duration(intervalRefreshToMongoSec) * time.Second)
			m.refreshToMongo()
		}
	}()
}

func regexConstruct(events []string) []bson.RegEx {
	regexs := make([]bson.RegEx, len(events))
	for i, v := range events {
		regexs[i] = bson.RegEx{Pattern: `^` + v + `$`, Options: ""}
	}
	return regexs
}

func buildMongoPipeOpsOfQueryAggregatedData(appid string, channel string, eventFilters []string, start, end time.Time, os string) []bson.M {
	// Filter channel_id first. Assuming that it could reduce a lot of results.
	operations := []bson.M{}
	if len(eventFilters) != 0 {
		if os == "" {
			operations = append(operations, bson.M{
				"$match": bson.M{"appid": appid, "channel": channel, "timestamp": bson.M{"$gte": start, "$lt": end}, "event": bson.M{"$in": regexConstruct(eventFilters)}},
			})
		} else {
			operations = append(operations, bson.M{
				"$match": bson.M{"appid": appid, "channel": channel, "timestamp": bson.M{"$gte": start, "$lt": end}, "os": os, "event": bson.M{"$in": regexConstruct(eventFilters)}},
			})
		}
	} else {
		if os == "" {
			operations = append(operations, bson.M{
				"$match": bson.M{"appid": appid, "channel": channel, "timestamp": bson.M{"$gte": start, "$lt": end}},
			})
		} else {
			operations = append(operations, bson.M{
				"$match": bson.M{"appid": appid, "channel": channel, "timestamp": bson.M{"$gte": start, "$lt": end}, "os": os},
			})
		}
	}

	operations = append(operations, []bson.M{
		bson.M{
			"$group": bson.M{
				"_id": "$event",
				"counts": bson.M{
					"$push": bson.M{
						"timestamp": "$timestamp",
						"count":     "$count",
					},
				},
			},
		},
	}...)
	return operations
}

func GetColl(mgoSession *mgo.Session, dbName string, colls map[string]*mgo.Collection, collNamePrefix, appid string) *mgo.Collection {
	if collNamePrefix == "" {
		log.Fatal("generalAggregateService collNamePrefix is empty")
		return nil
	}
	if mgoSession == nil {
		log.Fatal("generalAggregateService mgoSession is nil")
		return nil
	}
	mgoDb := mgoSession.DB(dbName)
	if mgoDb == nil {
		log.Fatal("generalAggregateService mgoDb is nil")
		return nil
	}
	if colls == nil {
		colls = make(map[string]*mgo.Collection)
	}
	if coll := colls[appid]; coll != nil {
		return coll
	}

	//for old version data, use the collection name without appID
	if appid == fakedAppIDForCompatibleColl {
		coll := mgoDb.C(collNamePrefix)
		colls[appid] = coll
		return coll
	}
	coll := mgoDb.C(collNamePrefix + "_" + appid)
	colls[appid] = coll
	return coll
}
