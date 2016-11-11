package attribution

import (
	"encoding/json"

	"github.com/MISingularity/deepshare2/deepstats"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"github.com/nsqio/go-nsq"
	"gopkg.in/mgo.v2"
)

type AttributionConsumer struct {
	as AttributionService
}

func NewAttributionConsumer(nsqlookupAddr, topic, channel string, p messaging.Producer, mgoSession *mgo.Session, dbName, collNamePrefix string, db storage.SimpleKV) (*nsq.Consumer, error) {
	nsqConsumer := deepstats.MustCreateNSQConsumerObj(topic, channel)

	ar := NewAttributionRetriever(mgoSession, dbName, collNamePrefix, db)
	as := NewAttributionService(ar, p)

	ac := &AttributionConsumer{as: as}

	err := messaging.NsqlookupdConsumeMessage(nsqConsumer, nsqlookupAddr, ac)
	return nsqConsumer, err
}

func (ac *AttributionConsumer) Consume(msg *nsq.Message) error {
	log.Debug("AttributionService; Receive Msg : ", string(msg.Body))
	e := &messaging.Event{}
	err := json.Unmarshal(msg.Body, e)
	if err != nil {
		log.Errorf("AttributionService; Unmarshal message body failed! Err Msg=%v", err)
		return err
	}

	return ac.as.OnEvent(e)
}
