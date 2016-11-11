package appchannel

import (
	"encoding/json"

	"github.com/MISingularity/deepshare2/deepstats"
	"github.com/MISingularity/deepshare2/deepstats/appchannel"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/nsqio/go-nsq"
	"gopkg.in/mgo.v2"
)

type AppchannelConsumer struct {
	appchannelService appchannel.AppChannelService
}

func ApplyNSQConsumer(c *mgo.Collection, nsqsel, nsqlookupdAddr, nsqdAddr, nsqtopic, nsqchannel string, maxinflight int) (*nsq.Consumer, error) {
	var err error
	appchannelConsumer := &AppchannelConsumer{appchannel.NewMongoAppChannelService(c)}
	appchannelNSQConsumer := deepstats.MustCreateNSQConsumerObj(nsqtopic, nsqchannel)
	appchannelNSQConsumer.ChangeMaxInFlight(4)
	if nsqsel == "nsqlookupd" {
		err = messaging.NsqlookupdConsumeMessage(appchannelNSQConsumer, nsqlookupdAddr, appchannelConsumer)
	} else {
		err = messaging.NsqdConsumeMessage(appchannelNSQConsumer, nsqdAddr, appchannelConsumer)
	}
	if err != nil {
		return nil, err
	}
	return appchannelNSQConsumer, nil
}

func (c *AppchannelConsumer) Consume(msg *nsq.Message) error {
	var dp messaging.Event
	log.Info("Receive Msg : ", string(msg.Body))
	if string(msg.Body) == "Aggregate" {
		return nil
	}
	err := json.Unmarshal(msg.Body, &dp)
	if err != nil {
		log.Errorf("Marshal message body failed! Err Msg=%v", err)
		return err
	}
	for i, _ := range dp.Channels {
		_, err = c.appchannelService.InsertChannel(dp.AppID, dp.Channels[i])
		if err != nil {
			log.Fatalf("Query appchannel service failed! Err Msg=%v", err)
			return err
		}
	}
	return nil
}
