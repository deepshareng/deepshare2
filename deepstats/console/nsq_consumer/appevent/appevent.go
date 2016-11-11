package appevent

import (
	"encoding/json"
	"strings"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepstats"
	"github.com/MISingularity/deepshare2/deepstats/appevent"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/nsqio/go-nsq"
	"gopkg.in/mgo.v2"
)

type AppeventConsumer struct {
	appeventService appevent.AppEventService
}

func (c *AppeventConsumer) Consume(msg *nsq.Message) error {
	var dp messaging.Event
	log.Info("Receive Msg : ", string(msg.Body))

	err := json.Unmarshal(msg.Body, &dp)
	if err != nil {
		log.Errorf("Marshal message body failed! Err Msg=%v", err)
		return err
	}
	if strings.HasPrefix(dp.EventType, api.CounterPrefix) {
		_, err = c.appeventService.InsertEvent(dp.AppID, dp.EventType)
		if err != nil {
			log.Fatalf("Query appevent service failed! Err Msg=%v", err)
			return err
		}
	}

	return nil
}
func NewAppeventConsumer(c *mgo.Collection) *AppeventConsumer {
	return &AppeventConsumer{appevent.NewMongoAppEventService(c)}
}

func ApplyNSQConsumer(c *mgo.Collection, nsqsel, nsqlookupdAddr, nsqdAddr, nsqtopic, nsqchannel string, maxinflight int) (*nsq.Consumer, error) {
	var err error
	appeventConsumer := NewAppeventConsumer(c)
	appeventNSQConsumer := deepstats.MustCreateNSQConsumerObj(nsqtopic, nsqchannel)
	appeventNSQConsumer.ChangeMaxInFlight(4)
	if nsqsel == "nsqlookupd" {
		err = messaging.NsqlookupdConsumeMessage(appeventNSQConsumer, nsqlookupdAddr, appeventConsumer)
	} else {
		err = messaging.NsqdConsumeMessage(appeventNSQConsumer, nsqdAddr, appeventConsumer)
	}
	if err != nil {
		return nil, err
	}
	return appeventNSQConsumer, nil
}
