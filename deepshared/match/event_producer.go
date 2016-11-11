package match

import (
	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg/messaging"
)

// Always use service as prefix
const (
	EventInstall      = api.MatchPrefix + "install"
	EventOpen         = api.MatchPrefix + "open"
	EventMatchUnknown = api.MatchPrefix + "unknown"
	EventBind         = api.MatchPrefix + "bind"
)

//produceBindEvent when binding some payload to user session
//mp:	the payload to bind, including senderID and channels info
//ua:	User Agent info about the client - remote address, os, os version, device brand,
//      browser info, and other information that can be parsed from the http request
func produceBindEvent(appID string, mp *MatchPayload, cookieID string, ua *uainfo.UAInfo,
	p messaging.Producer) error {
	l := &messaging.Event{
		AppID:     appID,
		EventType: EventBind,
		Channels:  mp.SenderInfo.Channels,
		SenderID:  mp.SenderInfo.SenderID,
		CookieID:  cookieID,
		Count:     1,
		UAInfo:    *ua,
		KVs: map[string]interface{}{
			"ua":         uainfo.NewUAFingerPrinter(ua).Transform(),
			"cookie_id":  cookieID,
			"inapp_data": mp.InappData,
		},
	}
	p.Produce(messaging.MatchTopic, l)
	return nil
}

//ProduceMatchEvent when a user try to retrieve some payload that saved by previous binding
//mp:	the payload retrieved
//receiverInfo: device properties about the client, for tracking
//ua:	User Agent info about the client - remote address, os, os version, device brand, browser info, and other information that can be parsed from the http request
func produceMatchEvent(appID string, mp *MatchPayload, cookieID string, tracking string,
	receiverInfo api.MatchReceiverInfo, ua *uainfo.UAInfo, p messaging.Producer) error {
	event := EventMatchUnknown
	switch tracking {
	case "install":
		event = EventInstall
	case "open":
		event = EventOpen
	}
	e := &messaging.Event{
		AppID:     appID,
		EventType: event,
		Channels:  mp.SenderInfo.Channels,
		SenderID:  mp.SenderInfo.SenderID,
		UniqueID:  receiverInfo.UniqueID,
		CookieID:  "",
		Count:     1,
		UAInfo:    *ua,
		KVs: map[string]interface{}{
			"ua":         uainfo.NewUAFingerPrinter(ua).Transform(),
			"cookie_id":  cookieID,
			"inapp_data": mp.InappData,
		},
	}
	p.Produce(messaging.MatchTopic, e)
	return nil
}
