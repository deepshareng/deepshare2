package main

import (
	"flag"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepstats"
	"github.com/MISingularity/deepshare2/deepstats/aggregate"
	"github.com/MISingularity/deepshare2/deepstats/appchannel"
	"github.com/MISingularity/deepshare2/deepstats/appevent"
	"github.com/MISingularity/deepshare2/deepstats/device_stat"
	"github.com/MISingularity/deepshare2/pkg/httputil"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/storage"
)

func main() {
	fs := flag.NewFlagSet("deepstatsd", flag.ExitOnError)
	hl := fs.String("http-listen", "0.0.0.0:16759", "HTTP/HTTPs Host and Port to listen on")
	mongoAddr := fs.String("mongo-addr", "", "Specify the raw data mongo database URL")
	mongoDBDayAggreagate := fs.String("mongodb-day-aggregate", "deepstats", "Specify the Mongo database for aggregate")
	mongoCollDayAggregate := fs.String("mongocoll-day-aggregate", "day", "Specify the Mongo collection for aggregate")

	mongoDBTotalAggreagate := fs.String("mongodb-total-aggregate", "deepstats", "Specify the Mongo database for aggregate")
	mongoCollTotalAggregate := fs.String("mongocoll-total-aggregate", "total", "Specify the Mongo collection for aggregate")

	mongoDBHourAggreagate := fs.String("mongodb-hour-aggregate", "deepstats", "Specify the Mongo database for aggregate")
	mongoCollHourAggregate := fs.String("mongocoll-hour-aggregate", "hour", "Specify the Mongo collection for aggregate")

	mongoDBAppchannel := fs.String("mongodb-appchannel", "deepstats", "Specify the Mongo database for app channel")
	mongoCollAppChannel := fs.String("mongocoll-appchannel", "appchannel", "Specify the Mongo collection for app channel")
	mongoDBAppevent := fs.String("mongodb-appevent", "deepstats", "Specify the Mongo database for app event")
	mongoCollAppevent := fs.String("mongocoll-appevent", "appevent", "Specify the Mongo collection for app event")

	urlRedis := fs.String("redis-url", "", "The redis url to use as DB")
	passwordRedis := fs.String("redis-password", "", "The redis password to use as DB")
	sentinelUrls := fs.String("redis-sentinel-urls", "", "The redis sentinel urls, splited by , if there are multiple")
	clusterNodeUrl := fs.String("redis-cluster-node-url", "", "redis 3.0 cluster nodes urls")
	redisMasterName := fs.String("redis-master-name", "mymaster", "The redis master name on sentinel deployment")
	poolSizeRedis := fs.Int("redis-pool-size", 100, "Specify the size of connection pool to redis")

	log.InitLog("[DEEPSTATS][API]", "", log.LevelDebug)
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	if *mongoAddr == "" {
		log.Fatal("mongo-addr is not set!")
	}
	session := deepstats.MustCreateMongoSession(*mongoAddr)
	defer session.Close()
	collNameHour := *mongoCollHourAggregate
	collNameDay := *mongoCollDayAggregate
	collNameTotal := *mongoCollTotalAggregate
	collAppchannelService := session.DB(*mongoDBAppchannel).C(*mongoCollAppChannel)
	collAppEventService := session.DB(*mongoDBAppevent).C(*mongoCollAppevent)

	redisDB := storage.NewInMemSimpleKV()

	if *clusterNodeUrl != "" {
		redisClusterUrls := strings.Split(*clusterNodeUrl, ",")
		log.Info("RedisCluster urls:", redisClusterUrls, "pool size:", *poolSizeRedis)
		redisDB = storage.NewRedisClusterSimpleKV(redisClusterUrls, *passwordRedis, *poolSizeRedis)
	} else if *sentinelUrls != "" {
		log.Info("RedisSentinel urls:", *sentinelUrls, "master name:", *redisMasterName, "pool size:", *poolSizeRedis)
		urls := strings.Split(*sentinelUrls, ",")
		redisDB = storage.NewRedisSentinelSimpleKV(urls, *redisMasterName, *passwordRedis, *poolSizeRedis)
	} else if *urlRedis != "" {
		log.Info("Redis url:", *urlRedis, "pool size:", *poolSizeRedis)
		redisDB = storage.NewRedisSimpleKV(*urlRedis, *passwordRedis, *poolSizeRedis)
	}
	mux := http.NewServeMux()
	aggregate.AddHandler(mux, api.ChannelPrefix, nil, session, session, session, *mongoDBHourAggreagate, *mongoDBDayAggreagate, *mongoDBTotalAggreagate, collNameHour, collNameDay, collNameTotal)
	appchannel.AddHandler(mux, api.AppChannelPrefix, collAppchannelService)
	appevent.AddHandler(mux, api.AppEventPrefix, collAppEventService)
	device_stat.AddHandler(mux, api.DeviceStatPrefix, redisDB, device_stat.RedisPrefixDevice)

	handler := httputil.RequestLogger(mux)
	if handler == nil {
		log.Fatal("handler is nil, impossible")
	}
	hs := http.Server{
		Addr:         *hl,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      handler,
	}
	log.Infof("main: serving HTTP on %s", hs.Addr)
	log.Fatal(hs.ListenAndServe())

}
