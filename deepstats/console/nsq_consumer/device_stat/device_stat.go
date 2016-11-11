package device_number

import (
	"encoding/json"

	"github.com/MISingularity/deepshare2/deepstats"
	"github.com/MISingularity/deepshare2/deepstats/device_stat"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"github.com/nsqio/go-nsq"
)

type DeviceStatConsumer struct {
	deviceStatService device_stat.DeviceStatService
}

func ApplyNSQConsumer(skv storage.SimpleKV, devicePrefix string, nsqsel, nsqlookupdAddr, nsqdAddr, nsqtopic, nsqchannel string, maxinflight int) (*nsq.Consumer, error) {
	var err error
	deviceStatConsumer := &DeviceStatConsumer{deviceStatService: device_stat.NewGeneralDeviceStat(skv, devicePrefix)}
	deviceStatNSQConsumer := deepstats.MustCreateNSQConsumerObj(nsqtopic, nsqchannel)
	deviceStatNSQConsumer.ChangeMaxInFlight(4)
	if nsqsel == "nsqlookupd" {
		err = messaging.NsqlookupdConsumeMessage(deviceStatNSQConsumer, nsqlookupdAddr, deviceStatConsumer)
	} else {
		err = messaging.NsqdConsumeMessage(deviceStatNSQConsumer, nsqdAddr, deviceStatConsumer)
	}
	if err != nil {
		return nil, err
	}
	return deviceStatNSQConsumer, nil
}

func (c *DeviceStatConsumer) Consume(msg *nsq.Message) error {
	var dp messaging.Event
	log.Info("Receive Msg : ", string(msg.Body))

	err := json.Unmarshal(msg.Body, &dp)
	if err != nil {
		log.Errorf("Marshal message body failed! Err Msg=%v\n", err)
		return err
	}
	err = c.deviceStatService.Insert(&dp)
	if err != nil {
		log.Errorf("Insert device into device stat failed! Err Msg=%v\n", err)
		return err
	}
	return nil
}
