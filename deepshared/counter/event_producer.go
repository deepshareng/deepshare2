package counter

import (
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg/messaging"
)

//ProduceCounterEvent when a user want to attach some relevant event to our tracking system
//cp:   counterRequest packages the receiver info and the event info, which are user want to
//      attach to our tracking system.
//ua:	User Agent info about the client - remote address, os, os version, device brand,
//      browser info, and other information that can be parsed from the http request
func produceCounterEvents(appID string, cp *CounterRequest,
	ua *uainfo.UAInfo, p messaging.Producer, endpoint string) error {
	for _, c := range cp.Counters {
		event := messaging.Event{
			AppID:     appID,
			EventType: endpoint + c.Event,
			UniqueID:  cp.ReceiverInfo.UniqueID,
			Count:     c.Count,
			UAInfo:    *ua,
		}

		p.Produce(messaging.CounterTopic, &event)
	}
	return nil
}
