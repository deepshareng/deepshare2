package retention

import (
	"encoding/json"

	"github.com/MISingularity/deepshare2/deepstats"
	"github.com/MISingularity/deepshare2/deepstats/retention"
	"github.com/MISingularity/deepshare2/pkg/condition"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/nsqio/go-nsq"
)

// RetentionConsumer used as retention service to consumer message
type RetentionConsumer struct {
	retentionService retention.RetentionService
}

func NewRetentionConsumer(rs retention.RetentionService) RetentionConsumer {
	return RetentionConsumer{rs}
}

func (r *RetentionConsumer) Consume(msg *nsq.Message) error {
	log.Info("Receive Msg : ", string(msg.Body))
	var dp messaging.Event
	err := json.Unmarshal(msg.Body, &dp)
	if err != nil {
		log.Errorf("Marshal message body failed! Err Msg=%v", err)
		return err
	}
	if condition.IsOpenEvent(&dp) {
		return r.retentionService.FindMatchUserForRetention(&dp)
	}
	if condition.IsInstallEvent(&dp) && condition.ByDeepshare(&dp) {
		return r.retentionService.InsertInstallEventForRetention(&dp)
	}
	return nil
}

// TO-DO:
// time zone difference is not fix

func ApplyNSQConsumer(retentionConsumer *RetentionConsumer, nsqsel, nsqlookupdAddr, nsqdAddr, nsqtopic, nsqchannel string, maxinflight int) (*nsq.Consumer, error) {
	var err error
	retentionNSQConsumer := deepstats.MustCreateNSQConsumerObj(nsqtopic, nsqchannel)
	retentionNSQConsumer.ChangeMaxInFlight(maxinflight)
	if nsqsel == "nsqlookupd" {
		err = messaging.NsqlookupdConsumeMessage(retentionNSQConsumer, nsqlookupdAddr, retentionConsumer)
	} else {
		err = messaging.NsqdConsumeMessage(retentionNSQConsumer, nsqdAddr, retentionConsumer)
	}
	if err != nil {
		return nil, err
	}
	return retentionNSQConsumer, nil
}
