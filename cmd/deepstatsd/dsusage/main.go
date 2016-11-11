package main

import (
	"encoding/json"
	"flag"
	"os"
	"strings"

	"github.com/MISingularity/deepshare2/deepstats"
	"github.com/MISingularity/deepshare2/deepstats/dsusage"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/nsqio/go-nsq"
)

const (
	countPrefix = "match/"
)

type aggregateConsumer struct {
	aggregateService dsusage.AggregateSenderService
	nsqConsumer      *nsq.Consumer
}

func applyNSQConsumer(redisAddrs []string, password string, poolSize int, nsqsel, nsqlookupdAddr, nsqdAddr string, nsqtopic, nsqchannel string, maxinflight int) (*nsq.Consumer, error) {
	var err error
	aggregateNSQConsumer := deepstats.MustCreateNSQConsumerObj(nsqtopic, nsqchannel)
	aggService := dsusage.NewAggregateSenderService(redisAddrs, password, poolSize)

	aggregateConsumer := &aggregateConsumer{aggService, aggregateNSQConsumer}
	aggregateNSQConsumer.ChangeMaxInFlight(maxinflight)
	if nsqsel == "nsqlookupd" {
		err = messaging.NsqlookupdConsumeMessage(aggregateNSQConsumer, nsqlookupdAddr, aggregateConsumer)
	} else {
		err = messaging.NsqdConsumeMessage(aggregateNSQConsumer, nsqdAddr, aggregateConsumer)
	}
	if err != nil {
		return nil, err
	}
	return aggregateNSQConsumer, nil
}

func main() {
	fs := flag.NewFlagSet("deepstatsd", flag.ExitOnError)
	nsqsel := fs.String("nsqsel", "nsqlookupd", "Specify the way to get nsq message, nsqlookupd/nsqd")
	nsqdAddr := fs.String("nsqd-tcp-addr", "", "Specify the nsqd adress")
	nsqlookupdAddr := fs.String("nsqlookupd-http-addr", "", "Specify the nsqlookupd adress")

	nsqTopics := fs.String("topics", "counter,match,sharelink,dsaction,genurl", "Specify the NSQ topic for consume, flag format should be topic1, topic2...")
	nsqChannel := fs.String("channel", "deepstats_aggregate", "Specify the NSQ channel for consumer")
	clusterNodeUrl := fs.String("redis-cluster-node-url", "", "redis 3.0 cluster nodes urls")
	passwordRedis := fs.String("redis-password", "", "The redis password to use as DB")
	poolSizeRedis := fs.Int("redis-pool-size", 100, "Specify the size of connection pool to redis")

	log.InitLog("[DEEPSTATS][AGGREGATE]", "", log.LevelDebug)
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	if *nsqsel == "nsqd" && *nsqdAddr == "" || *nsqsel == "nsqlookupd" && *nsqlookupdAddr == "" {
		log.Fatal("Nsq message genrator addr is not set! Please see nsqset/nsqd-tcp-addr/nsqlookupd-http-addr usage.")
	}
	if *clusterNodeUrl == "" {
		log.Fatal("redis-cluster-node-url is not set!")
	}

	topics := strings.Split(*nsqTopics, ",")
	nsqConsumers := make([]*nsq.Consumer, len(topics))
	for i, topic := range topics {
		var err error
		log.Infof("Listening topic=%s, channel=%s\n", topic, *nsqChannel)
		redisClusterUrls := strings.Split(*clusterNodeUrl, ",")
		log.Info("RedisCluster urls:", redisClusterUrls, "pool size:", *poolSizeRedis)

		nsqConsumers[i], err = applyNSQConsumer(redisClusterUrls, *passwordRedis, *poolSizeRedis, *nsqsel, *nsqlookupdAddr, *nsqdAddr, topic, *nsqChannel, 4)
		if err != nil {
			log.Fatalf("Construct the nsq consumer of aggregate service failed, error message = %v", err)
		}
	}

	// Wait nsq consumer stop
	// We need to wait all nsq consumers stop,
	// thus we could wait every consumer stop in particular order,
	// instead of waiting all consumer stop repeatly.
	for _, v := range nsqConsumers {
		<-v.StopChan
	}

	// <-aggregateNSQConsumer.StopChan
	log.Infof("Processing bufferred message finished.")
}

func (c *aggregateConsumer) Consume(msg *nsq.Message) error {
	log.Info("Receive Msg : ", string(msg.Body))
	var dp messaging.Event
	err := json.Unmarshal(msg.Body, &dp)
	if err != nil {
		log.Errorf("Marshal message body failed! Err Msg=%v", err)
		return err
	}

	if len(dp.Channels) == 1 && dp.Channels[0] == "" {
		dp.Channels = []string{}
	}

	if dp.SenderID != "" && dp.UniqueID != "" && dp.SenderID != dp.UniqueID {
		err := c.aggregateService.Insert(&dp)
		if err != nil {
			return err
		}
	}

	return nil
}
