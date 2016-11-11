package main

import (
	"flag"
	"os"
	"strings"

	nsqConsumer "github.com/MISingularity/deepshare2/deepstats/console/nsq_consumer/retention"
	"github.com/MISingularity/deepshare2/deepstats/retention"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"github.com/nsqio/go-nsq"
)

func main() {
	fs := flag.NewFlagSet("deepstatsd", flag.ExitOnError)
	nsqsel := fs.String("nsqsel", "nsqlookupd", "Specify the way to get nsq message, nsqlookupd/nsqd")
	nsqdAddr := fs.String("nsqd-tcp-addr", "", "Specify the nsqd adress")
	nsqlookupdAddr := fs.String("nsqlookupd-http-addr", "", "Specify the nsqlookupd adress")

	urlNSQ := fs.String("nsq-url", "", "The nsqd url to produce message")
	nsqTopics := fs.String("topics", "inappdata", "Specify the NSQ topic for consume, flag format should be topic1, topic2...")
	nsqChannel := fs.String("channel", "deepstats_retention", "Specify the NSQ channel for consumer")
	retentionTopic := fs.String("retention-topic", "retention", "Specify the NSQ Topic for retention service produce its message")
	retentionName := fs.String("retention-service-name", "3-day-retention", "Specify the unique retention service")
	retentionDay := fs.Int("retention-day", 3, "Specify the duration of retention calculated")

	addrRedis := fs.String("redis-addr", "", "The redis url to use as DB")
	passwordRedis := fs.String("redis-password", "", "The redis password to use as DB")
	log.InitLog("[DEEPSTATS][RENTENTION]", "", log.LevelDebug)
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	if *nsqsel == "nsqd" && *nsqdAddr == "" || *nsqsel == "nsqlookupd" && *nsqlookupdAddr == "" {
		log.Fatal("nsq message genrator addr is not set! Please see nsqset/nsqd-tcp-addr/nsqlookupd-http-addr usage.")
	}

	log.Debug("Redis url:", *addrRedis)
	db := storage.NewRedisClient(*addrRedis, *passwordRedis)

	log.InitLog("[DEEPSTATS][RETENTION]", "", log.LevelDebug)
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	if *nsqlookupdAddr == "" {
		log.Fatal("nsqlookupd addr is needed")
	}
	topics := strings.Split(*nsqTopics, ",")

	log.Debug("Start retention retrive match user...")
	nsqConsumers := make([]*nsq.Consumer, len(topics))
	p, err := messaging.NewNSQProducer(*urlNSQ, log.GetInfoLogger(), nsq.LogLevelDebug)
	if err != nil {
		log.Fatal("Failed to new NSQ client:", err)
	}

	rs := nsqConsumer.NewRetentionConsumer(retention.NewRedisRetentionService(*retentionDay, *retentionName, db, []byte(*retentionTopic), p))

	for i, topic := range topics {
		var err error
		log.Infof("Listening topic=%s, channel=%s\n", topic, *nsqChannel)
		nsqConsumers[i], err = nsqConsumer.ApplyNSQConsumer(&rs, *nsqsel, *nsqlookupdAddr, *nsqdAddr, topic, *nsqChannel, 4)
		if err != nil {
			log.Fatalf("Construct the nsq consumer of retention service failed, error message = %v", err)
			return
		}
	}
	log.Debug("Start Retention Service...")

	// Wait nsq consumer stop
	// We need to wait all nsq consumers stop,
	// thus we could wait every consumer stop in particular order,
	// instead of waiting all consumer stop repeatly.
	for _, v := range nsqConsumers {
		<-v.StopChan
	}

}
