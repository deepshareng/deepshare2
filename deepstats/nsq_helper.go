package deepstats

import (
	"time"

	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/nsqio/go-nsq"
)

func MustCreateNSQConsumerObj(topic, channel string) *nsq.Consumer {
	config := nsq.NewConfig()
	config.DefaultRequeueDelay = 0
	config.MaxBackoffDuration = 50 * time.Millisecond
	c, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		log.Fatal(err)
	}
	return c
}
