/*
Package messaging implements the messaging layer of deepshare backend.
*/

package messaging

import "github.com/nsqio/go-nsq"

// Producer is the interface that wraps the basic Produce method.
//
// Produce sends the given message to the given topic. Produce is
// an async method, and it will never block.
// There is no guarantee about the delivery of the Produce method.
type Producer interface {
	Produce(topic []byte, e *Event)
}

type MultiProducer interface {
	MultiProduce(topic []byte, e []*Event)
}

// Consumer consumes messages from message queue which procuder sends message in.

// Different implementation consume messages differently for each channel.
type Consumer interface {
	// Fetch returns messages if succeeded.
	Consume(msg *nsq.Message) error
}

func NsqdConsumeMessage(nc *nsq.Consumer, nsqdAddr string, c Consumer) error {
	hf := func(message *nsq.Message) error {
		return c.Consume(message)
	}
	nc.AddHandler(nsq.HandlerFunc(hf))
	return nc.ConnectToNSQD(nsqdAddr)
}

func NsqlookupdConsumeMessage(nc *nsq.Consumer, nsqlookupdAddr string, c Consumer) error {
	hf := func(message *nsq.Message) error {
		return c.Consume(message)
	}
	nc.AddHandler(nsq.HandlerFunc(hf))
	return nc.ConnectToNSQLookupd(nsqlookupdAddr)
}
