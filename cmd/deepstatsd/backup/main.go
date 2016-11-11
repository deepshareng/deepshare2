package main

import (
	"flag"
	"os"
	"strings"

	"github.com/MISingularity/deepshare2/deepstats"
	"github.com/MISingularity/deepshare2/deepstats/backup"
	nsqConsumer "github.com/MISingularity/deepshare2/deepstats/console/nsq_consumer/backup"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/nsqio/go-nsq"
)

func main() {
	fs := flag.NewFlagSet("deepstatsd-backup", flag.ExitOnError)
	nsqsel := fs.String("nsqsel", "nsqd", "Specify the way to get nsq message, nsqlookupd/nsqd")
	nsqdAddr := fs.String("nsqd-tcp-addr", "", "Specify the nsqd adress")
	nsqlookupdAddr := fs.String("nsqlookupd-http-addr", "", "Specify the nsqlookupd adress")

	nsqTopics := fs.String("topics", "counter,match,sharelink,dsaction,genurl,inappdata,retention", "Specify the NSQ topic for consume, flag format should be topic1, topic2...")
	nsqChannel := fs.String("channel", "deepstats_backup", "Specify the NSQ channel for consumer")

	backupsel := fs.String("backupsel", "mongo-raw", "Specify the ways to backup our message, currenly only support localfs/mongo-compress-day/mongo-raw.")
	storagePath := fs.String("path", "", "Local filesystem backup path(support localfs backup service)")
	mongoAddr := fs.String("mongo-addr", "", "Specify the mongo database URL(support mongo backup service)")
	mongoDB := fs.String("mongodb", "deepstats", "Specify the Mongo database")
	mongoColl := fs.String("mongocoll", "backup", "Specify the Mongo collection")

	log.InitLog("[BACKUP]", "", log.LevelDebug)
	var bs backup.BackupService
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	if *nsqsel == "nsqd" && *nsqdAddr == "" || *nsqsel == "nsqlookupd" && *nsqlookupdAddr == "" {
		log.Fatal("nsq message genrator addr is not set! Please see nsqset/nsqd-tcp-addr/nsqlookupd-http-addr usage.")
	}
	switch *backupsel {
	case "localfs":
		if *storagePath == "" {
			log.Fatal("Local filesystem backup path is not specified!")
		}
		var err error
		bs, err = backup.NewLocalFSBackupService(*storagePath)
		if err != nil {
			log.Fatalf("Create Mongo backup service failed! Err Msg=%v", err)
		}
	case "mongo-compress-day", "mongo-raw":
		if *mongoAddr == "" {
			log.Fatal("Mongo address is not specified!")
		}
		session := deepstats.MustCreateMongoSession(*mongoAddr)
		defer session.Close()
		db := session.DB(*mongoDB)
		c := db.C(*mongoColl)
		switch *backupsel {
		case "mongo-compress-day":
			bs = backup.NewMongoCompressDayBackupService(c)
		case "mongo-raw":
			bs = backup.NewMongoRawDayBackupService(session, *mongoDB)
		}
	default:
		log.Fatal("Specify the ways to backup our message, currenly only support localfs/mongo-compress-day/mongo-raw.")
	}

	topics := strings.Split(*nsqTopics, ",")

	nsqConsumers := make([]*nsq.Consumer, len(topics))
	for i, topic := range topics {
		var err error
		log.Info("Listening topic=%s, channel=%s\n", topic, *nsqChannel)
		nsqConsumers[i], err = nsqConsumer.ApplyNSQConsumer(bs, *nsqsel, *nsqlookupdAddr, *nsqdAddr, topic, *nsqChannel, 4)
		if err != nil {
			log.Fatalf("Construct the nsq consumer of aggregate service failed, error message = %v", err)
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
