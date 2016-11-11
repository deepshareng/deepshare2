package sharelink

import (
	"errors"

	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg"
	"github.com/MISingularity/deepshare2/pkg/messaging"
)

const (
	EventKeyShareLinkInfo = "shareLinkEventInfo"
	EventTypePrefix       = "sharelink:"
)

func AddEventInfo(e *messaging.Event, appId, endpoint, channels, senderId, cookieId, ip, ua, dsTag string) *messaging.Event {
	e.Count = 1
	e.AppID = appId
	e.EventType = EventTypePrefix + endpoint
	e.Channels = pkg.DecodeStringSlice(channels)
	e.SenderID = senderId
	e.CookieID = cookieId
	uaInfo := uainfo.ExtractUAInfoFromUAString(ip, ua)
	e.UAInfo = *uaInfo
	if e.KVs == nil {
		e.KVs = make(map[string]interface{})
	}
	e.KVs["ds_tag"] = dsTag

	return e
}

func AddEventShortUrlToken(e *messaging.Event, shortToken string) *messaging.Event {
	if e.KVs == nil {
		e.KVs = make(map[string]interface{})
	}
	e.KVs["shorturl_token"] = shortToken
	e.KVs["shorturl_token_valid"] = true
	return e
}

func AddEventShortSegValid(e *messaging.Event, isShortSegValid bool) *messaging.Event {
	if e.KVs == nil {
		e.KVs = make(map[string]interface{})
	}
	e.KVs["shorturl_token_valid"] = isShortSegValid
	return e
}

func AddEventOfAppInfo(e *messaging.Event, isyyb, isUL bool) *messaging.Event {
	if e.KVs == nil {
		e.KVs = make(map[string]interface{})
	}
	e.KVs["is_yyb"] = isyyb
	e.KVs["is_universallink"] = isUL
	return e
}

func FireShareLinkEvent(e *messaging.Event, p messaging.Producer) error {
	if p == nil {
		return errors.New("Share link; message producer is nil")
	}
	p.Produce(messaging.ShareLinkTopic, e)
	return nil
}
