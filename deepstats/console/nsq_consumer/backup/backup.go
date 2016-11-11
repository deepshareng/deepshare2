package backup

import (
	"encoding/json"

	"github.com/MISingularity/deepshare2/deepstats"
	"github.com/MISingularity/deepshare2/deepstats/backup"
	pb "github.com/MISingularity/deepshare2/deepstats/backup/proto"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/nsqio/go-nsq"
)

type BackupConsumer struct {
	backupService backup.BackupService
}

func ApplyNSQConsumer(bs backup.BackupService, nsqsel, nsqlookupdAddr, nsqdAddr, nsqtopic, nsqchannel string, maxinflight int) (*nsq.Consumer, error) {
	var err error
	backupConsumer := &BackupConsumer{bs}
	backupNSQConsumer := deepstats.MustCreateNSQConsumerObj(nsqtopic, nsqchannel)
	backupNSQConsumer.ChangeMaxInFlight(4)
	if nsqsel == "nsqlookupd" {
		err = messaging.NsqlookupdConsumeMessage(backupNSQConsumer, nsqlookupdAddr, backupConsumer)
	} else {
		err = messaging.NsqdConsumeMessage(backupNSQConsumer, nsqdAddr, backupConsumer)
	}
	if err != nil {
		return nil, err
	}
	return backupNSQConsumer, nil
}

func convertPBEvent(dp messaging.Event) (pb.Event, error) {
	kvs, err := json.Marshal(dp.KVs)
	if err != nil {
		return pb.Event{}, err
	}
	return pb.Event{
		AppID:     dp.AppID,
		EventType: dp.EventType,
		Channels:  dp.Channels,
		SenderID:  dp.SenderID,
		CookieID:  dp.CookieID,
		UniqueID:  dp.UniqueID,
		KVs:       kvs,
		Count:     int64(dp.Count),
		UAInfo: &pb.UAInfo{
			Ua:          dp.UAInfo.Ua,
			Ip:          dp.UAInfo.Ip,
			Os:          dp.UAInfo.Os,
			OsVersion:   dp.UAInfo.OsVersion,
			Brand:       dp.UAInfo.Brand,
			Browser:     dp.UAInfo.Browser,
			IsWechat:    dp.UAInfo.IsWechat,
			IsWeibo:     dp.UAInfo.IsWeibo,
			IsQQ:        dp.UAInfo.IsQQ,
			ChromeMajor: int64(dp.UAInfo.ChromeMajor),
		},
		Timestamp: dp.TimeStamp,
	}, nil

}

func (c *BackupConsumer) Consume(msg *nsq.Message) error {
	var dp messaging.Event
	log.Info("Receive Msg : ", string(msg.Body))
	if string(msg.Body) == "Aggregate" {
		return nil
	}
	err := json.Unmarshal(msg.Body, &dp)
	if err != nil {
		log.Errorf("Marshal message body failed! Err Msg=%v", err)
		return err
	}
	event, err := convertPBEvent(dp)
	if err != nil {
		log.Errorf("Convert PB Event failed! Err Msg=%v", err)
		return err
	}
	err = c.backupService.Insert(event)
	if err != nil {
		log.Fatalf("Query appchannel service failed! Err Msg=%v", err)
		return err
	}

	return nil
}
