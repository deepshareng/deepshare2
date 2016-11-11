package attribpush

import (
	"encoding/json"

	"github.com/MISingularity/deepshare2/deepstats"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/nsqio/go-nsq"
)

type AttributionConsumer struct {
	ar           AttributionRetriver
	produceTopic []byte
	producer     messaging.Producer
}

func NewAttributionConsumer(nsqlookupAddr, topic, channel, nsqdAddr, produceTopic string, mongoAddr, mongoDBName, collName string) (*nsq.Consumer, error) {
	nsqConsumer := deepstats.MustCreateNSQConsumerObj(topic, channel)
	mgoSession := deepstats.MustCreateMongoSession(mongoAddr)
	db := mgoSession.DB(mongoDBName)
	coll := db.C(collName)
	attribConsumer := &AttributionConsumer{ar: NewAttributionParser(coll)}

	if p, err := messaging.NewNSQProducer(nsqdAddr, log.GetInfoLogger(), nsq.LogLevelDebug); err != nil {
		return nil, err
	} else {
		attribConsumer.producer = p
		attribConsumer.produceTopic = []byte(produceTopic)
	}
	err := messaging.NsqlookupdConsumeMessage(nsqConsumer, nsqlookupAddr, attribConsumer)
	return nsqConsumer, err
}

func (ac *AttributionConsumer) Consume(msg *nsq.Message) error {
	log.Info("[AttributionConsumer] Receive Msg : ", string(msg.Body))
	var e messaging.Event
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		log.Errorf("Unmarshal message body failed! Err Msg=%v", err)
		return err
	}
	senderID, tag := ac.ar.Retrive(e)
	log.Debug("[AttributionConsumer] parsed attribution info:", senderID, tag)

	if senderID != "" && tag != "" {
		ac.produceBackToNSQ(&e, senderID, tag)
	}

	return nil
}

func (ac *AttributionConsumer) produceBackToNSQ(event *messaging.Event, senderID, tag string) {
	event.SenderID = senderID
	event.EventType = tag
	ac.producer.Produce(ac.produceTopic, event)
}
