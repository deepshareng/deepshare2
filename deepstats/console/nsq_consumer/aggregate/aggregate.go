package aggregate

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepstats"
	"github.com/MISingularity/deepshare2/deepstats/aggregate"
	"github.com/MISingularity/deepshare2/pkg/condition"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"github.com/nsqio/go-nsq"
	"gopkg.in/mgo.v2"
)

type AggregateConsumer struct {
	aggregateService aggregate.AggregateService
}

func NewAggregateConsumer(service string, redisDB storage.SimpleKV, session *mgo.Session, dbName, collNamePrefix string) AggregateConsumer {
	var aggService aggregate.AggregateService
	switch service {
	case "hour":
		aggService = aggregate.NewHourAggregateService(redisDB, session, dbName, collNamePrefix, "", "")
	case "day":
		aggService = aggregate.NewDayAggregateService(redisDB, session, dbName, collNamePrefix, "", "")
	case "total":
		aggService = aggregate.NewTotalAggregateService(redisDB, session, dbName, collNamePrefix, "", "")
	default:
		log.Fatalf("No such aggregate service, please use hour/day/total to specify!")
	}
	return AggregateConsumer{aggService}
}

func (c *AggregateConsumer) insert(dp *messaging.Event, i int, event string) error {
	os := "other"
	if dp.UAInfo.IsAndroid() {
		os = "android"
	}
	if dp.UAInfo.IsIos() {
		os = "ios"
	}
	countevent := aggregate.CounterEvent{
		AppID:     dp.AppID,
		Channel:   dp.Channels[i],
		Event:     event,
		Count:     dp.Count,
		Timestamp: time.Unix(dp.TimeStamp, 0),
		Os:        os,
	}

	err := c.aggregateService.Insert(dp.AppID, countevent)
	if err != nil {
		log.Fatalf("Query aggregate service failed! Err Msg=%v", err)
		return err
	}

	// This part convert channel from inappdata(i.e. v2/inappdata)
	// into our eventual demonstration column(match/)

	if condition.IsInstallEvent(dp) || condition.IsOpenEvent(dp) {
		event = condition.SubstitutePrefix(event)
	} else {
		return nil
	}

	countevent = aggregate.CounterEvent{
		AppID:     dp.AppID,
		Channel:   dp.Channels[i],
		Event:     event,
		Count:     dp.Count,
		Timestamp: time.Unix(dp.TimeStamp, 0),
		Os:        os,
	}

	err = c.aggregateService.Insert(dp.AppID, countevent)
	if err != nil {
		log.Fatalf("Query aggregate service failed! Err Msg=%v", err)
		return err
	}
	return nil
}

func (c *AggregateConsumer) Diffuse(dp messaging.Event) error {
	if condition.IsEmptyChannel(dp.Channels) {
		dp.Channels = []string{}
	}

	inappdata := condition.HasInappdata(&dp)
	bydeepshare := condition.ByDeepshare(&dp)
	// sharelink dstag examination
	// if dsg is not empty, we add suffix "_ds_tag", otherwise add "_no_ds_tag"
	// including two steps:
	//   - check if the event is attributed to sharelink.
	//   - check if the dstag is in Kvs.
	if condition.IsSharelink(&dp) {
		if condition.HasDstag(&dp) {
			dp.EventType += "_ds_tag"
		} else {
			dp.EventType += "_no_ds_tag"
		}
	}

	if !strings.HasPrefix(dp.EventType, api.RetentionAmountPrefix) {
		dp.Channels = append(dp.Channels, "all")
	}

	for i, _ := range dp.Channels {
		err := c.insert(&dp, i, dp.EventType)
		if err != nil {
			return err
		}
		if inappdata {
			err = c.insert(&dp, i, dp.EventType+"_with_inappdata")
			if err != nil {
				return err
			}
		}
		if bydeepshare {
			err := c.insert(&dp, i, dp.EventType+"_with_params")
			if err != nil {
				return err
			}
		}

		err = c.aggregateService.Aggregate(dp.AppID)
		if err != nil {
			log.Fatalf("Aggregate service aggregates failed! Err Msg=%v", err)
			return err
		}
	}

	return nil
}

func (c *AggregateConsumer) Consume(msg *nsq.Message) error {
	log.Info("Receive Msg : ", string(msg.Body))
	var dp messaging.Event
	err := json.Unmarshal(msg.Body, &dp)
	if err != nil {
		log.Errorf("Marshal message body failed! Err Msg=%v", err)
		return err
	}
	if v, ok := dp.KVs["short_seg"]; dp.EventType == api.GetInAppDataPrefix+"open" && ok && v != "" && v != nil {
		err := c.Diffuse(dp)
		if err != nil {
			return err
		}
		dp.EventType = api.GenerateUrlPrefix
	}
	return c.Diffuse(dp)
}

func ApplyNSQConsumer(service string, redisDB storage.SimpleKV, session *mgo.Session, dbName, collNamePrefix, nsqsel, nsqlookupdAddr, nsqdAddr, nsqtopic, nsqchannel string, maxinflight int) (*nsq.Consumer, aggregate.AggregateService, error) {
	var err error
	aggregateNSQConsumer := deepstats.MustCreateNSQConsumerObj(nsqtopic, nsqchannel)
	var aggService aggregate.AggregateService
	switch service {
	case "hour":
		aggService = aggregate.NewHourAggregateService(redisDB, session, dbName, collNamePrefix, nsqtopic, nsqchannel)
	case "day":
		aggService = aggregate.NewDayAggregateService(redisDB, session, dbName, collNamePrefix, nsqtopic, nsqchannel)
	case "total":
		aggService = aggregate.NewTotalAggregateService(redisDB, session, dbName, collNamePrefix, nsqtopic, nsqchannel)
	default:
		log.Fatalf("No such aggregate service, please use hour/day/total to specify!")
	}

	aggregateConsumer := &AggregateConsumer{aggService}
	aggregateNSQConsumer.ChangeMaxInFlight(maxinflight)
	if nsqsel == "nsqlookupd" {
		err = messaging.NsqlookupdConsumeMessage(aggregateNSQConsumer, nsqlookupdAddr, aggregateConsumer)
	} else {
		err = messaging.NsqdConsumeMessage(aggregateNSQConsumer, nsqdAddr, aggregateConsumer)
	}
	if err != nil {
		return nil, nil, err
	}
	return aggregateNSQConsumer, aggService, nil
}
