package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/appcookiedevice"
	"github.com/MISingularity/deepshare2/deepshared/appinfo"
	"github.com/MISingularity/deepshare2/deepshared/counter"
	"github.com/MISingularity/deepshare2/deepshared/devicecookier"
	"github.com/MISingularity/deepshare2/deepshared/dsaction"
	"github.com/MISingularity/deepshare2/deepshared/match"
	"github.com/MISingularity/deepshare2/deepshared/token"
	"github.com/MISingularity/deepshare2/frontend/binddevicetocookie"
	"github.com/MISingularity/deepshare2/frontend/inappdata"
	"github.com/MISingularity/deepshare2/frontend/jsapi"
	"github.com/MISingularity/deepshare2/frontend/sharelink"
	"github.com/MISingularity/deepshare2/frontend/urlgenerator"
	"github.com/MISingularity/deepshare2/frontendplus/dsusage"
	"github.com/MISingularity/deepshare2/pkg/httputil"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"github.com/nsqio/go-nsq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/cors"
)

func main() {
	fs := flag.NewFlagSet("deepshared", flag.ExitOnError)
	sts := fs.String("service-types", "appcookiedevice,devicecookie,match,counter,dsaction,appinfo,token,sharelinkfront,urlgenerator,jsapi,inappdata,binddevicetocookie,dsusage", "service type to start")
	hl := fs.String("http-listen", "0.0.0.0:8080", "HTTP/HTTPs Host and Port to listen on")
	hk := fs.String("http-key", "", "The key file HTTPs server to use")
	hc := fs.String("http-cert", "", "The certificate file HTTPs server to use")
	workerID := fs.Int64("worker-id", 0, "The worker ID of current docker")
	dataCenterID := fs.Int64("data-center-id", 0, "The data center ID of current docker")
	urlRedis := fs.String("redis-url", "", "The redis url to use as DB")
	passwordRedis := fs.String("redis-password", "", "The redis password to use as DB")
	sentinelUrls := fs.String("redis-sentinel-urls", "", "The redis sentinel urls, splited by , if there are multiple")
	clusterNodeUrl := fs.String("redis-cluster-node-url", "", "The redis cluster urls, splited by , if there are multiple")
	redisMasterName := fs.String("redis-master-name", "mymaster", "The redis master name on sentinel deployment")
	poolSizeRedis := fs.Int("redis-pool-size", 100, "Specify the size of connection pool to redis")
	urlNSQ := fs.String("nsq-url", "", "The redis url to use as DB")
	urlMatch := fs.String("match-url", "127.0.0.1:8080", "matching service url")
	urlCookie := fs.String("cookie-url", "127.0.0.1:8080", "cookie service url")
	urlAppCookie := fs.String("appcookie-url", "127.0.0.1:8080", "app cookie service url")
	urlAppInfo := fs.String("appinfo-url", "127.0.0.1:8080", "appinfo service url")
	urlUrlGenerator := fs.String("urlgenerator-url", "127.0.0.1:8080", "appinfo service url")
	urlToken := fs.String("token-url", "127.0.0.1:8080", "token service url")
	urlBaseUrlGenerator := fs.String("genurl-base", "http://127.0.0.1:8080", "generated url will start with this base")
	logLevel := fs.String("log-level", "error", "Specify the log level: debug, info, error")
	uaMatchValidMinutes := fs.Int("ua-match-valid-minutes", 15, "Specify ua match valid time in minutes")

	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Error(err)
		os.Exit(1)
	}

	switch *logLevel {
	case "debug":
		log.InitLog("[DEEPSHARED]", "", log.LevelDebug)
	case "info":
		log.InitLog("[DEEPSHARED]", "", log.LevelInfo)
	case "error":
		log.InitLog("[DEEPSHARED]", "", log.LevelError)
	}

	log.Info("Args from cmd line:")
	fs.VisitAll(func(f *flag.Flag) {
		fmt.Println("	", f.Name, ":", f.Value)
	})

	mux := http.NewServeMux()
	producer := messaging.NewSimpleProducer(nil)
	db := storage.NewInMemSimpleKV()

	if *urlNSQ != "" {
		log.Info("NSQ url:", *urlNSQ)
		p, err := messaging.NewNSQProducer(*urlNSQ, log.GetInfoLogger(), nsq.LogLevelDebug)
		if err != nil {
			log.Fatal("Failed to new NSQ client:", err)
		}
		producer = p
	}

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

	matchUrl := "http://" + path.Join(*urlMatch, api.MatchPrefix)
	cookieUrl := "http://" + path.Join(*urlCookie, api.DeviceCookiePrefix)
	appcookieUrl := "http://" + path.Join(*urlAppCookie, api.AppCookieDevicePrefix)
	appinfoUrl := "http://" + path.Join(*urlAppInfo, api.AppInfoPrefix)
	urlGeneratorUrl := "http://" + path.Join(*urlUrlGenerator, api.GenerateUrlPrefix)
	tokenUrl := "http://" + path.Join(*urlToken, api.TokenPrefix)
	log.Info("match url:        ", matchUrl)
	log.Info("cookie url:       ", cookieUrl)
	log.Info("appcookie url:    ", appcookieUrl)
	log.Info("appinfo url:      ", appinfoUrl)
	log.Info("urlgenerator url: ", urlGeneratorUrl)
	log.Info("token url:        ", tokenUrl)

	var urlHandler http.Handler = nil
	var jsApiHandler http.Handler = nil
	var dsActionHandler http.Handler = nil
	var needCors bool = false

	parts := strings.Split(*sts, ",")
	for _, st := range parts {
		switch st {
		case "appcookiedevice":
			log.Infof("Handle endpoint %s", api.AppCookieDevicePrefix)
			appcookiedevice.AddHandler(mux, api.AppCookieDevicePrefix, db, producer)
		case "devicecookie":
			log.Infof("Handle endpoint %s", api.DeviceCookiePrefix)
			devicecookier.AddHandler(mux, api.DeviceCookiePrefix, db, producer)
		case "match":
			log.Infof("Handle endpoint %s", api.MatchPrefix)
			match.AddHandler(mux, api.MatchPrefix, db, producer, int64(*uaMatchValidMinutes)*60)
		case "counter":
			log.Infof("Handle endpoint %s", api.CounterPrefix)
			counter.AddHandler(mux, api.CounterPrefix, producer)
		case "dsaction":
			log.Infof("Handle endpoint %s", api.DSActionsPrefix)
			dsActionHandler = dsaction.CreateHandler(api.DSActionsPrefix, producer)
			needCors = true
		case "appinfo":
			log.Infof("Handle endpoint %s", api.AppInfoPrefix)
			appinfo.AddHandler(mux, api.AppInfoPrefix, db, producer)
		case "token":
			log.Infof("Handle endpoint %s", api.TokenPrefix)
			token.AddHandler(mux, *workerID, *dataCenterID, api.TokenPrefix)
		case "sharelinkfront":
			log.Infof("Handle endpoint %s", api.ShareLinkPrefix)
			sharelink.AddHandler(mux, api.ShareLinkPrefix, *urlBaseUrlGenerator, db, matchUrl, cookieUrl, appcookieUrl, appinfoUrl, tokenUrl, producer)
		case "urlgenerator":
			urlHandler = urlgenerator.CreateHandler(api.GenerateUrlPrefix, *urlBaseUrlGenerator, db, tokenUrl, producer)
			needCors = true
		case "inappdata":
			log.Infof("Handle endpoint %s", api.GetInAppDataPrefix)
			inappdata.AddHandler(mux, api.GetInAppDataPrefix, matchUrl, cookieUrl, appcookieUrl, urlGeneratorUrl, appinfoUrl, producer)
		case "binddevicetocookie":
			log.Infof("Handle endpoint %s", api.BindDeviceToCookiePrefix)
			binddevicetocookie.AddHandler(mux, api.BindDeviceToCookiePrefix, cookieUrl, tokenUrl)
		case "dsusage":
			log.Infof("Handle endpoint %s", api.DSUsagesPrefix)
			dsusage.AddHandler(mux, api.DSUsagesPrefix, db)
		case "jsapi":
			log.Infof("Handle endpoint %s", api.JSApiPrefix)
			jsApiHandler = jsapi.CreateHandler(api.JSApiPrefix, *urlBaseUrlGenerator, db, matchUrl, cookieUrl, appcookieUrl, appinfoUrl, tokenUrl, producer)
			needCors = true
		}
	}
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./files"))))
	log.Infof("Handle endpoint %s", "/metrics")
	mux.Handle("/metrics", prometheus.Handler())

	//CORS APIs
	if needCors {
		c := cors.New(cors.Options{
			AllowOriginFunc: func(o string) bool {
				log.Info("CORS from:", o)
				return true
			},
			AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
			AllowCredentials: true,
			//allow all headers for CORS request (fix for Talicai)
			AllowedHeaders: []string{"*"},
			Debug:          true,
		})

		if urlHandler != nil {
			urlHandler = c.Handler(urlHandler)
			log.Infof("Handle endpoint %s (CORS)", api.GenerateUrlPrefix)
			mux.Handle(api.GenerateUrlPrefix, urlHandler)
		}

		if jsApiHandler != nil {
			jsApiHandler = c.Handler(jsApiHandler)
			log.Infof("Handle endpoint %s (CORS)", api.JSApiPrefix)
			mux.Handle(api.JSApiPrefix, jsApiHandler)
		}

		if dsActionHandler != nil {
			dsActionHandler = c.Handler(dsActionHandler)
			log.Infof("Handle endpoint %s (CORS)", api.DSActionsPrefix)
			mux.Handle(api.DSActionsPrefix, dsActionHandler)
		}
	}

	handler := httputil.RequestLogger(mux)
	if handler == nil {
		log.Fatal("[Error], handler is nil, impossible")
	}

	hs := http.Server{
		Addr:         *hl,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      handler,
	}

	if len(*hk) == 0 || len(*hc) == 0 {
		log.Infof("main: serving HTTP on %s", hs.Addr)
		log.Fatal(hs.ListenAndServe())
	}

	log.Infof("main: serving HTTPs on %s", hs.Addr)
	log.Fatal(hs.ListenAndServeTLS(*hk, *hc))
}
