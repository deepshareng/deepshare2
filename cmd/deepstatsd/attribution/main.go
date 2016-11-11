package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/MISingularity/deepshare2/deepstats"
	"github.com/MISingularity/deepshare2/deepstats/attribution"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"github.com/nsqio/go-nsq"
)

func main() {
	fs := flag.NewFlagSet("deepstatsd", flag.ExitOnError)
	nsqlookupdAddr := fs.String("nsqlookupd-http-addr", "127.0.0.1:4161", "Specify the nsqlookupd adress")
	nsqdAddr := fs.String("nsqd-tcp-addr", "127.0.0.1:4150", "Specify the nsqd tcp adress")
	nsqTopics := fs.String("topics", "counter_raw,inappdata,dsaction", "Specify the NSQ topic for consume, flag format should be topic1, topic2...")
	nsqChannel := fs.String("channel", "deepstats_attribution", "Specify the NSQ channel for consumer")
	mongoDB := fs.String("mongodb", "deepstats", "Specify the Mongo database")
	mongoColl := fs.String("mongocoll", "attribution", "Specify the Mongo collection")
	mongoAddr := fs.String("mongo-addr", "127.0.0.1:27017", "Specify the raw data mongo database URL")
	urlRedis := fs.String("redis-url", "", "The redis url to use as DB")
	passwordRedis := fs.String("redis-password", "", "The redis password to use as DB")
	sentinelUrls := fs.String("redis-sentinel-urls", "", "The redis sentinel urls, splited by , if there are multiple")
	clusterNodeUrl := fs.String("redis-cluster-node-url", "", "redis 3.0 cluster nodes urls")
	redisMasterName := fs.String("redis-master-name", "mymaster", "The redis master name on sentinel deployment")
	poolSizeRedis := fs.Int("redis-pool-size", 100, "Specify the size of connection pool to redis")
	log.InitLog("[DEEPSTATS][ATTRIBUTION]", "", log.LevelInfo)
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	log.Info("Args from cmd line:")
	fs.VisitAll(func(f *flag.Flag) {
		fmt.Println("	", f.Name, ":", f.Value)
	})
	if *nsqlookupdAddr == "" {
		log.Fatal("nsqlookupd addr is needed")
	}
	if *mongoAddr == "" {
		log.Fatal("mongo-addr is not set!")
	}
	topics := strings.Split(*nsqTopics, ",")
	redisKV := storage.NewInMemSimpleKV()
	if *clusterNodeUrl != "" {
		redisClusterUrls := strings.Split(*clusterNodeUrl, ",")
		log.Info("RedisCluster urls:", redisClusterUrls, "pool size:", *poolSizeRedis)
		redisKV = storage.NewRedisClusterSimpleKV(redisClusterUrls, *passwordRedis, *poolSizeRedis)
	} else if *sentinelUrls != "" {
		log.Info("RedisSentinel urls:", *sentinelUrls, "master name:", *redisMasterName, "pool size:", *poolSizeRedis)
		urls := strings.Split(*sentinelUrls, ",")
		redisKV = storage.NewRedisSentinelSimpleKV(urls, *redisMasterName, *passwordRedis, *poolSizeRedis)
	} else if *urlRedis != "" {
		log.Info("Redis url:", *urlRedis, "pool size:", *poolSizeRedis)
		redisKV = storage.NewRedisSimpleKV(*urlRedis, *passwordRedis, *poolSizeRedis)
	}

	log.Debug("Start attribution retrive consumer...")
	nsqConsumers := make([]*nsq.Consumer, len(topics))
	for i, topic := range topics {
		var err error
		log.Infof("Listening topic=%s, channel=%s\n", topic, *nsqChannel)
		mgoSession := deepstats.MustCreateMongoSession(*mongoAddr)
		p, err := messaging.NewNSQProducer(*nsqdAddr, log.GetInfoLogger(), nsq.LogLevelDebug)
		if err != nil {
			log.Fatal("Failed to new NSQ producer, err:", err)
		}
		nsqConsumers[i], err = attribution.NewAttributionConsumer(*nsqlookupdAddr, topic, *nsqChannel, p, mgoSession, *mongoDB, *mongoColl, redisKV)
		if err != nil {
			log.Fatalf("Construct the nsq consumer of attribution service failed, error message = %v", err)
			return
		}
	}

	// Wait nsq consumer stop
	// We need to wait all nsq consumers stop,
	// thus we could wait every consumer stop in particular order,
	// instead of waiting all consumer stop repeatly.
	for _, v := range nsqConsumers {
		<-v.StopChan
	}
}
