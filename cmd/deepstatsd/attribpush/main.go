package main

import (
	"flag"
	"net/http"
	"os"
	"strings"

	"github.com/MISingularity/deepshare2/deepstats/attribpush"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"github.com/nsqio/go-nsq"
)

func main() {
	fs := flag.NewFlagSet("deepstatsd", flag.ExitOnError)
	hl := fs.String("http-listen", "0.0.0.0:8081", "HTTP/HTTPs Host and Port to listen on")
	nsqlookupdAddr := fs.String("nsqlookupd-http-addr", "127.0.0.1:4161", "Specify the nsqlookupd adress")
	nsqdAddr := fs.String("nsqd-tcp-addr", "127.0.0.1:4150", "Specify the nsqd tcp adress")
	nsqTopics := fs.String("topics", "counter,match,dsaction", "Specify the NSQ topic for consume, flag format should be topic1, topic2...")
	nsqChannel := fs.String("channel", "deepstats_attribution", "Specify the NSQ channel for consumer")
	nsqChannelPush := fs.String("push-channel", "attribution_push", "Specify the NSQ channel for consumer")
	mongoDB := fs.String("mongodb", "deepstats", "Specify the Mongo database")
	mongoColl := fs.String("mongocoll", "appchannel", "Specify the Mongo collection")
	mongoAddr := fs.String("mongo-addr", "127.0.0.1:27017", "Specify the raw data mongo database URL")
	pushWorkerNum := fs.Int("push-worker-num", 2, "How many workers to work on pushing attribution")
	urlRedis := fs.String("redis-url", "", "The redis url to use as DB")
	passwordRedis := fs.String("redis-password", "", "The redis password to use as DB")
	sentinelUrls := fs.String("redis-sentinel-urls", "", "The redis sentinel urls, splited by , if there are multiple")
	redisMasterName := fs.String("redis-master-name", "mymaster", "The redis master name on sentinel deployment")
	poolSizeRedis := fs.Int("redis-pool-size", 100, "Specify the size of connection pool to redis")
	log.InitLog("[DEEPSTATS][ATTRIBUTION]", "", log.LevelDebug)
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	if *nsqlookupdAddr == "" {
		log.Fatal("nsqlookupd addr is needed")
	}
	if *mongoAddr == "" {
		log.Fatal("mongo-addr is not set!")
	}
	topics := strings.Split(*nsqTopics, ",")
	db := storage.NewInMemSimpleKV()
	if *sentinelUrls != "" {
		log.Info("RedisSentinel urls:", *sentinelUrls, "master name:", *redisMasterName, "pool size:", *poolSizeRedis)
		urls := strings.Split(*sentinelUrls, ",")
		db = storage.NewRedisSentinelSimpleKV(urls, *redisMasterName, *passwordRedis, *poolSizeRedis)
	} else if *urlRedis != "" {
		log.Info("Redis url:", *urlRedis, "pool size:", *poolSizeRedis)
		db = storage.NewRedisSimpleKV(*urlRedis, *passwordRedis, *poolSizeRedis)
	}
	mux := http.NewServeMux()
	au := attribpush.AddHandler(mux, "/apps/callbackurl/", db)

	go func() {
		log.Debug("Start attribution retrive consumer...")
		nsqConsumers := make([]*nsq.Consumer, len(topics))
		for i, topic := range topics {
			var err error
			log.Infof("Listening topic=%s, channel=%s\n", topic, *nsqChannel)
			nsqConsumers[i], err = attribpush.NewAttributionConsumer(*nsqlookupdAddr, topic, *nsqChannel, *nsqdAddr, attribpush.AttributionTopic, *mongoAddr, *mongoDB, *mongoColl)
			if err != nil {
				log.Fatalf("Construct the nsq consumer of attribution service failed, error message = %v", err)
				return
			}
		}

		log.Debug("Start Push Service...")
		attribpush.StartPushService(*nsqlookupdAddr, attribpush.AttributionTopic, *nsqChannelPush, *pushWorkerNum, 10, au, attribpush.NewInMemBuffer())

		// Wait nsq consumer stop
		// We need to wait all nsq consumers stop,
		// thus we could wait every consumer stop in particular order,
		// instead of waiting all consumer stop repeatly.
		for _, v := range nsqConsumers {
			<-v.StopChan
		}
	}()

	log.Debug("Start app url handler, http server Listening on:", *hl)
	log.Fatal(http.ListenAndServe(*hl, mux))

}
