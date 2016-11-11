package messaging

import (
	"encoding/json"
	"io"

	"github.com/MISingularity/deepshare2/pkg/log"
)

// simpleProducer is used to mock producer.
// It's mostly used in tests or when we don't need write anything to data pipeline.
type simpleProducer struct {
	w io.Writer
}

type SimpleProducerEvent struct {
	Topic []byte
	Msg   []byte
}

// NewSimpleProducer returns a Producer
// that writes json marshalled SimpleProducerEvents into the given io.writer
// If the given io.writer is nil, Produce() will be a noop.
func NewSimpleProducer(w io.Writer) Producer {
	return &simpleProducer{w}
}

func (p *simpleProducer) Produce(topic []byte, event *Event) {
	message, err := json.Marshal(event)
	if err != nil {
		log.Errorf("[simpleProducer]Failed to marshal event:%v err:%v\n", event, err)
	}
	log.Debugf("Produce an event (simpleProducer), topic:%s, message:%v", topic, string(message))
	if p.w == nil {
		log.Debug("p.w is nil, ignore the event")
		return
	}
	en := json.NewEncoder(p.w)
	e := &SimpleProducerEvent{
		Topic: topic,
		Msg:   message,
	}
	if err := en.Encode(e); err != nil {
		log.Error("failed to encode event, err:", err)
	}
}
