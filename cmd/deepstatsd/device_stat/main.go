package main

import (
	"flag"
	"net/http"
	"os"
	"strings"

	nsqConsumer "github.com/MISingularity/deepshare2/deepstats/console/nsq_consumer/device_stat"
	"github.com/MISingularity/deepshare2/deepstats/device_stat"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"github.com/nsqio/go-nsq"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	fs := flag.NewFlagSet("deepstatsd", flag.ExitOnError)
	hl := fs.String("http-listen", "0.0.0.0:16759", "HTTP/HTTPs Host and Port to listen on")
	nsqsel := fs.String("nsqsel", "nsqlookupd", "Specify the way to get nsq message, nsqlookupd/nsqd")
	nsqdAddr := fs.String("nsqd-tcp-addr", "", "Specify the nsqd adress")
	nsqlookupdAddr := fs.String("nsqlookupd-http-addr", "", "Specify the nsqlookupd adress")
	nsqTopics := fs.String("topics", "counter,sharelink,dsaction,genurl,inappdata,match", "Specify the NSQ topic for consume, flag format should be topic1, topic2...")
	nsqChannel := fs.String("channel", "deepstats_device_stat", "Specify the NSQ channel for consumer")
	urlRedis := fs.String("redis-url", "", "The redis url to use as DB")
	passwordRedis := fs.String("redis-password", "", "The redis password to use as DB")
	sentinelUrls := fs.String("redis-sentinel-urls", "", "The redis sentinel urls, splited by , if there are multiple")
	clusterNodeUrl := fs.String("redis-cluster-node-url", "", "redis 3.0 cluster nodes urls")
	redisMasterName := fs.String("redis-master-name", "mymaster", "The redis master name on sentinel deployment")
	poolSizeRedis := fs.Int("redis-pool-size", 100, "Specify the size of connection pool to redis")

	log.InitLog("[DEEPSTATS][DEVICE_STAT]", "", log.LevelDebug)
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	if *nsqsel == "nsqd" && *nsqdAddr == "" || *nsqsel == "nsqlookupd" && *nsqlookupdAddr == "" {
		log.Fatal("nsq message genrator addr is not set! Please see nsqset/nsqd-tcp-addr/nsqlookupd-http-addr usage.")
	}

	// simple kv storage initialization
	db := storage.NewInMemSimpleKV()

	if *clusterNodeUrl != "" {
		redisClusterUrls := strings.Split(*clusterNodeUrl, ",")
		log.Info("RedisCluster urls:", redisClusterUrls, "pool size:", *poolSizeRedis)
		db = storage.NewRedisClusterSimpleKV(redisClusterUrls, *passwordRedis, *poolSizeRedis)
	} else if *sentinelUrls != "" {
		log.Info("RedisSentinel urls:", *sentinelUrls, "master name:", *redisMasterName, "pool size:", *poolSizeRedis)
		urls := strings.Split(*sentinelUrls, ",")
		db = storage.NewRedisSentinelSimpleKV(urls, *redisMasterName, *passwordRedis, *poolSizeRedis)
	} else if *urlRedis != "" {
		log.Info("Redis url:", *urlRedis, "pool size:", *poolSizeRedis)
		db = storage.NewRedisSimpleKV(*urlRedis, *passwordRedis, *poolSizeRedis)
	}

	topics := strings.Split(*nsqTopics, ",")

	// This part waits any nsq topic we required
	nsqConsumers := make([]*nsq.Consumer, len(topics))
	for i, topic := range topics {
		var err error
		log.Info("Listening topic=%s, channel=%s\n", topic, *nsqChannel)
		nsqConsumers[i], err = nsqConsumer.ApplyNSQConsumer(db, device_stat.RedisPrefixDevice, *nsqsel, *nsqlookupdAddr, *nsqdAddr, topic, *nsqChannel, 4)
		if err != nil {
			log.Fatalf("Construct the nsq consumer of aggregate service failed, error message = %v", err)
			return
		}
	}

	go func() {
		http.Handle("/metrics", prometheus.Handler())
		log.Info("HTTP Listen on ", *hl)
		if err := http.ListenAndServe(*hl, nil); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait nsq consumer stop
	// We need to wait all nsq consumers stop,
	// thus we could wait every consumer stop in particular order,
	// instead of waiting all consumer stop repeatly.
	for _, v := range nsqConsumers {
		<-v.StopChan
	}

}
