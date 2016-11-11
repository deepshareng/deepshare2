package retention

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepstats/aggregate"
	"github.com/MISingularity/deepshare2/deepstats/appchannel"
	"github.com/MISingularity/deepshare2/deepstats/console/job"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type RetentionJob struct {
	serviceName                                                      string
	appchannelCollection, lastTimeCollection, lastTimeMarkCollection *mgo.Collection
	deepstats_url                                                    string
	produceTopic                                                     []byte
	producer                                                         messaging.MultiProducer
	lastTime                                                         time.Time
	gran                                                             time.Duration
	convert                                                          func(time.Time) time.Time
}

func NewRetentionJob(serviceName string, appchannelCollection, lastTimeCollection, lastTimeChannelMarkCollection *mgo.Collection, deepstats_url string, produceTopic []byte, p messaging.MultiProducer, gran time.Duration, convert func(time.Time) time.Time) *RetentionJob {
	lastTime := time.Now().Add(-gran)
	res := job.LastTimeServiceTime{}
	err := lastTimeCollection.Find(bson.M{"service-name": serviceName}).One(&res)

	if err == nil {
		lastTime = res.LastTime
	} else {
		if strings.Contains(err.Error(), "not found") {
			err = lastTimeCollection.Insert(bson.M{"service-name": serviceName, "lasttime": lastTime})
			if err != nil {
				log.Fatalf("Lasttime collection insert failed! Err Msg=%v", err)
			}
		} else {
			log.Fatalf("Lasttime collection find failed! Err Msg=%v", err)
		}

	}
	return &RetentionJob{
		serviceName:            serviceName,
		lastTimeCollection:     lastTimeCollection,
		appchannelCollection:   appchannelCollection,
		lastTimeMarkCollection: lastTimeChannelMarkCollection,
		deepstats_url:          deepstats_url,
		produceTopic:           produceTopic,
		producer:               p,
		lastTime:               lastTime,
		gran:                   gran,
		convert:                convert,
	}
}

func (r *RetentionJob) Run() {
	limit := int(math.Floor(float64(r.convert(time.Now()).Sub(r.convert(r.lastTime))) / float64(r.gran)))
	r.Process(r.convert(r.lastTime), limit)
	// last time is first unprocessed time of job, signify the start point of duration if we processed
	r.lastTime = r.convert(time.Now())
	_, err := r.lastTimeCollection.Upsert(bson.M{"service-name": r.serviceName}, job.LastTimeServiceTime{ServiceName: r.serviceName, LastTime: r.lastTime})
	if err != nil {
		log.Fatal(err)
	}
}

// TODO:
// add os dimension
func (r *RetentionJob) Process(processStart time.Time, processLimit int) {
	result := make([]appchannel.AppChannel, 0)
	err := r.appchannelCollection.Find(bson.M{}).All(&result)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		log.Fatalf("Retention service access app channel mongo failed! Err Msg=%v!", err)
		return
	}
	end := processStart.Add(r.gran * time.Duration(processLimit))
	for _, v := range result {
		// Get channel's last time incompleted.
		// Decide whether calculate it or not on the premise of comprision between the newest time and channel's last time finished.
		if v.Channel == "" {
			continue
		}
		result := job.LastTimeChannelTime{}
		err := r.lastTimeMarkCollection.Find(bson.M{"service-name": r.serviceName, "channel": v.Channel}).One(&result)
		start := time.Now()
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				err := r.lastTimeMarkCollection.Insert(bson.M{"channel": v.Channel, "service-name": r.serviceName, "lasttime": processStart})
				if err != nil {
					log.Fatalf("Retention service access mongo failed! Err Msg=%v!", err)
					return
				}
				start = processStart
			} else {
				log.Fatal(err)
				return
			}
		} else {
			start = result.LastTime
		}
		limit := int(math.Floor(float64(end.Sub(r.convert(start))) / float64(r.gran)))
		log.Infof("Processing channel %s, start = %v, limit = %d, gran = %s,\n", v.Channel, start, limit, r.gran)
		res, err := r.GetChannelCounters(v.AppID, v.Channel, []string{api.RetentionPrefix + r.serviceName, api.RetentionPrefix + r.serviceName + "_install"}, start, "d", strconv.Itoa(limit))
		if err != nil {
			log.Error(err)
			continue
		}
		for dateNum := 0; dateNum < limit; dateNum++ {
			install, retention := 0, 0
			for _, s := range res {
				switch s.Event {
				case api.RetentionPrefix + r.serviceName:
					if len(s.Counts) > 0 {
						retention = s.Counts[dateNum].Count
					}
				case api.RetentionPrefix + r.serviceName + "_install":
					if len(s.Counts) > 0 {
						install = s.Counts[dateNum].Count
					}
				}
			}
			if install != 0 {
				rate := int(retention * 10000 / install)
				e1 := &messaging.Event{
					AppID:     v.AppID,
					EventType: api.RetentionAmountPrefix + r.serviceName + "_retention-rate",
					UniqueID:  "",
					Channels:  []string{v.Channel},
					SenderID:  "",
					CookieID:  "",
					Count:     rate,
					KVs:       map[string]interface{}{},
					TimeStamp: start.Add(r.gran * time.Duration(dateNum)).Unix(),
				}

				e2 := &messaging.Event{
					AppID:     v.AppID,
					EventType: api.RetentionAmountPrefix + r.serviceName + "_retention-day",
					UniqueID:  "",
					Channels:  []string{v.Channel},
					SenderID:  "",
					CookieID:  "",
					Count:     1,
					KVs:       map[string]interface{}{},
					TimeStamp: start.Add(r.gran * time.Duration(dateNum)).Unix(),
				}
				r.producer.MultiProduce(r.produceTopic, []*messaging.Event{e1, e2})
			}
			// After sending messages to NSQ, we finished this time's all jobs. As the time
			// varies from earlier to later, we could change the newest time of this channel
			// signify we have processed all jobs until this time.
			_, err := r.lastTimeMarkCollection.Upsert(bson.M{"channel": v.Channel, "service-name": r.serviceName}, bson.M{"$set": bson.M{"lasttime": start.Add(r.gran * time.Duration(dateNum+1))}})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func eventSerialize(events []string, appid string, start time.Time, gran string, limit string) string {
	serialize := "?appid=" + url.QueryEscape(appid) + "&gran=" + gran + "&start=" + strconv.FormatInt(start.Unix(), 10) + "&limit=" + limit
	for _, v := range events {
		serialize += "&event=" + url.QueryEscape(v)
	}
	return serialize
}

func (r *RetentionJob) GetChannelCounters(appid string, channel_id string, event []string, start time.Time, gran string, limit string) ([]*aggregate.AggregateResult, error) {
	log.Infof("Request GetChannelCounters, URL=%s", (r.deepstats_url + "/v2/channels/" + channel_id + "/counters" + eventSerialize(event, appid, start, gran, limit)))
	res, err := http.Get(r.deepstats_url + "/v2/channels/" + channel_id + "/counters" + eventSerialize(event, appid, start, gran, limit))
	if err != nil {
		return []*aggregate.AggregateResult{}, fmt.Errorf("Request Deepstats failed%v", err)
	}
	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return []*aggregate.AggregateResult{}, err
	}

	channelInfo := struct {
		Counters []*aggregate.AggregateResult
	}{}

	err = json.Unmarshal(result, &channelInfo)
	if err != nil {
		return []*aggregate.AggregateResult{}, err
	}
	return channelInfo.Counters, nil
}
