package main

import (
	"flag"
	"net/http"
	"os"
	"strings"

	"github.com/MISingularity/deepshare2/deepstats"
	nsqConsumer "github.com/MISingularity/deepshare2/deepstats/console/nsq_consumer/appchannel"
	"github.com/MISingularity/deepshare2/pkg/log"
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
	nsqChannel := fs.String("channel", "deepstats_appchannel", "Specify the NSQ channel for consumer")
	mongoDB := fs.String("mongodb", "deepstats", "Specify the Mongo database")
	mongoColl := fs.String("mongocoll", "appchannel", "Specify the Mongo collection")
	mongoAddr := fs.String("mongo-addr", "", "Specify the raw data mongo database URL")
	log.InitLog("[DEEPSTATS][APPCHANNEL]", "", log.LevelDebug)
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	if *nsqsel == "nsqd" && *nsqdAddr == "" || *nsqsel == "nsqlookupd" && *nsqlookupdAddr == "" {
		log.Fatal("nsq message genrator addr is not set! Please see nsqset/nsqd-tcp-addr/nsqlookupd-http-addr usage.")
	}
	if *mongoAddr == "" {
		log.Fatal("mongo-addr is not set!")
	}

	topics := strings.Split(*nsqTopics, ",")
	session := deepstats.MustCreateMongoSession(*mongoAddr)
	defer session.Close()
	c := session.DB(*mongoDB).C(*mongoColl)

	// This part waits any nsq topic we required
	nsqConsumers := make([]*nsq.Consumer, len(topics))
	for i, topic := range topics {
		var err error
		log.Info("Listening topic=%s, channel=%s\n", topic, *nsqChannel)
		nsqConsumers[i], err = nsqConsumer.ApplyNSQConsumer(c, *nsqsel, *nsqlookupdAddr, *nsqdAddr, topic, *nsqChannel, 4)
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
