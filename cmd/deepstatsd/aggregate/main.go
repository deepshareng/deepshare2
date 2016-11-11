package main

import (
	"flag"
	"net/http"
	"os"
	"strings"

	"github.com/MISingularity/deepshare2/deepstats"
	"github.com/MISingularity/deepshare2/deepstats/aggregate"
	nsqConsumer "github.com/MISingularity/deepshare2/deepstats/console/nsq_consumer/aggregate"
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

	mongoAddr := fs.String("mongo-addr", "", "Specify the raw data mongo database URL")
	nsqTopics := fs.String("topics", "counter,sharelink,dsaction,genurl,inappdata,match,retention", "Specify the NSQ topic for consume, flag format should be topic1, topic2...")
	nsqChannel := fs.String("channel", "deepstats_aggregate", "Specify the NSQ channel for consumer")
	mongoDB := fs.String("mongodb", "deepstats", "Specify the Mongo database")
	mongoColl := fs.String("mongocoll", "day", "Specify the Mongo collection")
	aggregateService := fs.String("agg-service", "day", "Specify which aggregate service is being used, \"hour\" means aggregating events by hour. Currently support hour/day/total.")
	urlRedis := fs.String("redis-url", "", "The redis url to use as DB")
	passwordRedis := fs.String("redis-password", "", "The redis password to use as DB")
	sentinelUrls := fs.String("redis-sentinel-urls", "", "The redis sentinel urls, splited by , if there are multiple")
	clusterNodeUrl := fs.String("redis-cluster-node-url", "", "The redis cluster urls, splited by , if there are multiple")
	redisMasterName := fs.String("redis-master-name", "mymaster", "The redis master name on sentinel deployment")
	poolSizeRedis := fs.Int("redis-pool-size", 100, "Specify the size of connection pool to redis")

	log.InitLog("[DEEPSTATS][AGGREGATE]", "", log.LevelInfo)
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	if *nsqsel == "nsqd" && *nsqdAddr == "" || *nsqsel == "nsqlookupd" && *nsqlookupdAddr == "" {
		log.Fatal("Nsq message genrator addr is not set! Please see nsqset/nsqd-tcp-addr/nsqlookupd-http-addr usage.")
	}
	if *mongoAddr == "" {
		log.Fatal("mongo-addr is not set!")
	}

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
	session := deepstats.MustCreateMongoSession(*mongoAddr)
	defer session.Close()
	collNamePrefix := *mongoColl
	nsqConsumers := make([]*nsq.Consumer, len(topics))
	for i, topic := range topics {
		var as aggregate.AggregateService
		var err error
		log.Infof("Listening topic=%s, channel=%s\n", topic, *nsqChannel)
		nsqConsumers[i], as, err = nsqConsumer.ApplyNSQConsumer(*aggregateService, db, session, *mongoDB, collNamePrefix, *nsqsel, *nsqlookupdAddr, *nsqdAddr, topic, *nsqChannel, 1000)
		if err != nil {
			log.Fatalf("Construct the nsq consumer of aggregate service failed, error message = %v", err)
		}
		as.StartRefreshLoop()
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

	// <-aggregateNSQConsumer.StopChan
	log.Infof("Processing bufferred message finished.")
}
