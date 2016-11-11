package messaging

import (
	"encoding/json"
	"log"
	"time"

	dslog "github.com/MISingularity/deepshare2/pkg/log"
	"github.com/nsqio/go-nsq"
)

type nsqProducerAdatpor struct {
	w *nsq.Producer
}

func NewNSQProducer(nsqdAddr string, logger *log.Logger, logLevel nsq.LogLevel) (Producer, error) {
	config := nsq.NewConfig()
	w, err := nsq.NewProducer(nsqdAddr, config)
	if err != nil {
		return nil, err
	}
	if err := w.Ping(); err != nil {
		logger.Fatal(err)
	}
	if logger != nil {
		w.SetLogger(logger, logLevel)
	}
	return &nsqProducerAdatpor{w}, nil
}

func NewNSQMultiProducer(nsqdAddr string, logger *log.Logger, logLevel nsq.LogLevel) (MultiProducer, error) {
	config := nsq.NewConfig()
	w, err := nsq.NewProducer(nsqdAddr, config)
	if err != nil {
		return nil, err
	}
	if err := w.Ping(); err != nil {
		logger.Fatal(err)
	}
	if logger != nil {
		w.SetLogger(logger, logLevel)
	}
	return &nsqProducerAdatpor{w}, nil
}

func NewTestingNSQProducer(w *nsq.Producer) Producer {
	return &nsqProducerAdatpor{w}
}

func (p *nsqProducerAdatpor) Produce(topic []byte, event *Event) {
	if event.TimeStamp == 0 {
		event.TimeStamp = time.Now().Unix()
	}
	message, err := json.Marshal(event)
	if err != nil {
		dslog.Errorf("[nsqProducerAdatpor]Failed to marshal event:%v err:%v\n", event, err)
	}
	dslog.Infof("Produce an event to NSQ, topic:%s, message:%v", topic, string(message))
	if err := p.w.PublishAsync(string(topic), message, nil); err != nil {
		dslog.Fatalf("nsqProducerAdatpor PublishAsync failed: %v", err)
	}
}

func (p *nsqProducerAdatpor) MultiProduce(topic []byte, events []*Event) {
	messages := make([][]byte, len(events))
	for i, event := range events {
		if event.TimeStamp == 0 {
			event.TimeStamp = time.Now().Unix()
		}
		message, err := json.Marshal(event)
		if err != nil {
			dslog.Errorf("[nsqProducerAdatpor]Failed to marshal event:%v err:%v\n", event, err)
		}
		messages[i] = message
		dslog.Infof("Ready to produce an event to NSQ, topic:%s, message:%v", topic, string(message))
	}
	if err := p.w.MultiPublishAsync(string(topic), messages, nil); err != nil {
		dslog.Fatalf("nsqProducerAdatpor MultiPublishAsync failed: %v", err)
	}
	dslog.Info("Produce Successful!")

}
