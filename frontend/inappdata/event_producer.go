package inappdata

import (
	"path"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
)

const (
	SourceTypeFromRawUrl = "fromRawUrl"
	SourceTypeFromMatch  = "fromMatch"
)

func produceInAppDataFromRawUrlEvent(appID, endpoint, shortSeg, tracking, inapp_data, sender_id, ip, ua string, channels []string, receiverInfo api.MatchReceiverInfo, p messaging.Producer) {
	if p == nil {
		log.Error("messaging producer is nil")
		return
	}
	event := path.Join(endpoint, tracking)
	uaInfo := uainfo.ExtractUAInfoFromUAString(ip, ua)
	e := &messaging.Event{
		AppID:        appID,
		EventType:    event,
		UniqueID:     receiverInfo.UniqueID,
		Channels:     channels,
		SenderID:     sender_id,
		CookieID:     "",
		Count:        1,
		UAInfo:       *uaInfo,
		ReceiverInfo: receiverInfo,
		KVs: map[string]interface{}{
			"type":       SourceTypeFromRawUrl,
			"short_seg":  shortSeg,
			"inapp_data": inapp_data,
		},
	}
	p.Produce(messaging.InAppDataTopic, e)
}

func produceInAppDataFromMatchEvent(appID, endpoint, tracking, inapp_data, sender_id, ip, ua, clickId, cookieId string, channels []string, receiverInfo api.MatchReceiverInfo, p messaging.Producer) {
	if p == nil {
		log.Error("messaging producer is nil")
		return
	}
	event := path.Join(endpoint, tracking)
	uaInfo := uainfo.ExtractUAInfoFromUAString(ip, ua)
	e := &messaging.Event{
		AppID:        appID,
		EventType:    event,
		UniqueID:     receiverInfo.UniqueID,
		Channels:     channels,
		SenderID:     sender_id,
		CookieID:     "",
		Count:        1,
		UAInfo:       *uaInfo,
		ReceiverInfo: receiverInfo,
		KVs: map[string]interface{}{
			"type":       SourceTypeFromMatch,
			"inapp_data": inapp_data,
			"click_id":   clickId,
			"cookie_id":  cookieId,
		},
	}
	p.Produce(messaging.InAppDataTopic, e)
}

func produceInAppDataErrorEvent(appID, endpoint, tracking, ip, ua string, receiverInfo api.MatchReceiverInfo, err error, p messaging.Producer) {
	if p == nil {
		log.Error("messaging producer is nil")
		return
	}
	event := path.Join(endpoint, tracking)
	uaInfo := uainfo.ExtractUAInfoFromUAString(ip, ua)
	e := &messaging.Event{
		AppID:        appID,
		EventType:    event,
		UniqueID:     receiverInfo.UniqueID,
		CookieID:     "",
		Count:        1,
		UAInfo:       *uaInfo,
		ReceiverInfo: receiverInfo,
		KVs: map[string]interface{}{
			"error": err.Error(),
		},
	}
	p.Produce(messaging.InAppDataTopic, e)
}
