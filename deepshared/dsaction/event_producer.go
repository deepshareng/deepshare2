package dsaction

import (
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg/messaging"
)

//produceDSActionEvent is triggered from SDK side to attach some pre-defined event to our tracking system
// i.e. close, etc
//req:  DSActionRequest packages the receiver info and the action, which are user want to
//      attach to our tracking system.
//ua:	User Agent info about the client - remote address, os, os version, device brand,
//      browser info, and other information that can be parsed from the http request
func produceDSActionEvent(appID string, req *DSActionRequest, ua *uainfo.UAInfo, p messaging.Producer, endpoint string) error {
	event := messaging.Event{
		AppID:     appID,
		EventType: endpoint + req.Action,
		Channels:  req.Channels,
		SenderID:  req.SenderID,
		UniqueID:  req.ReceiverInfo.UniqueID,
		Count:     1,
		UAInfo:    *ua,
		KVs:       req.KVs,
	}
	p.Produce(messaging.DSActionTopic, &event)
	return nil
}
