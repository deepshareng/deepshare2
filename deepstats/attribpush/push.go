package attribpush

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/MISingularity/deepshare2/deepstats"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/nsqio/go-nsq"
)

type PushService struct {
	au AppToUrl
	ab AttributionBuffer
}

func StartPushService(nsqlookupAddr, topic, channel string, numOfWorkers int, pushIntervel int, au AppToUrl, ab AttributionBuffer) {
	ncs := []*nsq.Consumer{}
	for i := 0; i < numOfWorkers; i++ {
		nc := deepstats.MustCreateNSQConsumerObj(topic, channel)
		pc := &PushService{
			au: au,
			ab: ab,
		}
		go pc.startTick(pushIntervel)
		err := messaging.NsqlookupdConsumeMessage(nc, nsqlookupAddr, pc)
		if err != nil {
			log.Fatal(err)
		}
		ncs = append(ncs, nc)
	}
	for _, nc := range ncs {
		<-nc.StopChan
	}
}

func (pc *PushService) startTick(sec int) {
	for {
		timer := time.NewTimer(time.Second * time.Duration(sec))
		<-timer.C
		log.Info("Timer expired... need to push")

		for _, appID := range pc.ab.ListAppIDs() {
			attrs := pc.ab.PopAttributions(appID)
			url, _ := pc.au.GetUrl(appID)
			if url == "" || attrs == nil {
				continue
			}
			PostAttributions(url, attrs)
		}
	}
}

func (pc *PushService) Consume(msg *nsq.Message) error {
	log.Info("[PushService] Receive Msg : ", string(msg.Body))
	var e messaging.Event
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		log.Errorf("Unmarshal message body failed! Err Msg=%v", err)
		return err
	}
	attr := &AttributionPushInfo{
		SenderID:  e.SenderID,
		Tag:       e.EventType,
		Value:     e.Count,
		Timestamp: e.TimeStamp,
	}
	pc.ab.PutAttribution(e.AppID, attr)

	return nil
}

func PostAttributions(url string, attributions []*AttributionPushInfo) error {
	log.Debugf("[PostAttributions] url=%s, attributions=%#v\n", url, attributions)
	b, err := json.Marshal(attributions)
	log.Debug(string(b))
	if err != nil {
		log.Error("marshal data failed:", err)
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		log.Error("Failed to post data to server url:", url, "err:", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("Failed to post to server url:", url, "code:", resp.StatusCode)
		return errors.New("http request failed")
	}

	return nil
}
