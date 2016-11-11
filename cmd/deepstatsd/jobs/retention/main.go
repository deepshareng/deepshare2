package main

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/MISingularity/deepshare2/deepstats"
	"github.com/MISingularity/deepshare2/deepstats/console/job/retention"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/nsqio/go-nsq"
)

func main() {
	fs := flag.NewFlagSet("deepstatsd-crontask-retention", flag.ExitOnError)

	serviceName := fs.String("service-name", "3-day-retention", "Specify the service name")
	deepstatsURL := fs.String("deepstats-url", "", "The deepstats url to request data")
	urlNSQ := fs.String("nsq-url", "", "The nsqd url to produce message")
	retentionTopic := fs.String("retention-count-topic", "retention", "Specify the NSQ Topic for retention service produce its message")
	mongoAddr := fs.String("mongo-addr", "", "Specify the mongo database URL")
	mongoAppchannelDB := fs.String("appchannel-mongodb", "deepstats", "Specify the Mongo appchannel database")
	mongoAppchannelColl := fs.String("appchannel-mongocoll", "appchannel", "Specify the Mongo appchannel collection")
	mongoLastTimeDB := fs.String("job-last-time-mongodb", "lasttime", "Specify the Mongo database saving the first unprocessed time of job/channel, signify the job/channel before that time have been processed")
	mongoLastTimeColl := fs.String("job-last-time-mongocoll", "joblasttime", "Specify the Mongo collection saving the first unprocessed time of job")
	mongoLastTimeChannelMarkColl := fs.String("channel-mark-last-time-mongocoll", "channelmarklasttime", "Specify the Mongo collection saving the first unprocessed time of channel")

	log.InitLog("[DEEPSTATS][JOB][RENTENTION]", "", log.LevelDebug)
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	if *deepstatsURL == "" {
		log.Fatal("Deepshare URL is not specified!")
	}
	if !strings.HasPrefix(*deepstatsURL, "http://") {
		*deepstatsURL = "http://" + *deepstatsURL
	}

	session := deepstats.MustCreateMongoSession(*mongoAddr)
	defer session.Close()

	appchannelColl := session.DB(*mongoAppchannelDB).C(*mongoAppchannelColl)
	lasttimeColl := session.DB(*mongoLastTimeDB).C(*mongoLastTimeColl)
	lastTimeChannelMarkColl := session.DB(*mongoLastTimeDB).C(*mongoLastTimeChannelMarkColl)

	p, err := messaging.NewNSQMultiProducer(*urlNSQ, log.GetInfoLogger(), nsq.LogLevelDebug)
	if err != nil {
		log.Fatal("Failed to new NSQ client:", err)
	}
	job := retention.NewRetentionJob(
		*serviceName,
		appchannelColl,
		lasttimeColl,
		lastTimeChannelMarkColl,
		*deepstatsURL,
		[]byte(*retentionTopic),
		p,
		time.Hour*24,
		func(t time.Time) time.Time {
			return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
		},
	)
	ticker := time.NewTicker(1 * time.Hour)
	job.Run()
	for {
		select {
		case <-ticker.C:
			job.Run()
		}
	}
}
