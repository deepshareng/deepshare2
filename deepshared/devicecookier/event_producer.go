package devicecookier

import "github.com/MISingularity/deepshare2/pkg/messaging"

type EventProducer interface {
}

func NewEventProducer(p messaging.Producer) EventProducer {
	return &eventProducer{
		producer: p,
	}
}

type eventProducer struct {
	producer messaging.Producer
}
