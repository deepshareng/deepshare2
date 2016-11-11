package urlgenerator

import (
	"path"

	"net/url"

	"strings"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg/messaging"
)

const (
	EventInstall      = "/install"
	EventOpen         = "/open"
	EventMatchUnknown = "/unknown"
)

func produceGenerateUrlEvent(appID, rawUrl, shortUrl, endpoint string, genUrlPost *GenURLPostBody, ua *uainfo.UAInfo, p messaging.Producer) error {
	senderId := ""
	role := ""
	if genUrlPost.SenderID != "" {
		senderId = genUrlPost.SenderID
		role = "sender"
	} else {
		senderId = genUrlPost.ForwardedSenderID
		role = "receiver"
	}
	l := &messaging.Event{
		AppID:     appID,
		EventType: endpoint,
		Channels:  genUrlPost.Channels,
		SenderID:  senderId,
		CookieID:  "",
		Count:     1,
		UAInfo:    *ua,
		KVs: map[string]interface{}{
			"is_short":             genUrlPost.IsShort,
			"sdk_info":             genUrlPost.SDKInfo,
			"inapp_data":           genUrlPost.InAppData,
			"download_title":       genUrlPost.DownloadTitle,
			"download_btn_text":    genUrlPost.DownloadBtnText,
			"download_msg":         genUrlPost.DownloadMsg,
			"download_url_ios":     genUrlPost.DownloadUrlIos,
			"download_url_android": genUrlPost.DownloadUrlAndroid,
			"uninstall_url":        genUrlPost.UninstallUrl,
			"redirect_url":         genUrlPost.RedirectUrl,
			"rawUrl":               rawUrl,
			"shortUrl":             shortUrl,
			"role":                 role,
		},
	}
	p.Produce(messaging.GenUrlTopic, l)
	return nil
}

func produceGetInappDataEvent(appID, endpoint, shortSeg, tracking string, rawUrl *url.URL, receiverInfo api.MatchReceiverInfo, ua *uainfo.UAInfo, p messaging.Producer) error {
	event := path.Join(endpoint, EventMatchUnknown)
	switch tracking {
	case "install":
		event = path.Join(endpoint, EventInstall)
	case "open":
		event = path.Join(endpoint, EventOpen)
	}
	e := &messaging.Event{
		AppID:        appID,
		EventType:    event,
		UniqueID:     receiverInfo.UniqueID,
		CookieID:     "",
		Count:        1,
		UAInfo:       *ua,
		ReceiverInfo: receiverInfo,
		KVs: map[string]interface{}{
			"short_seg": shortSeg,
		},
	}
	values := rawUrl.Query()
	if contexts, ok := values["inapp_data"]; ok {
		e.KVs["inapp_data"] = contexts[0]
	}
	if contexts, ok := values["channels"]; ok {
		s := contexts[0]
		e.Channels = strings.Split(s, "|")
	}
	if contexts, ok := values["sender_id"]; ok {
		e.SenderID = contexts[0]
	}
	p.Produce(messaging.GenUrlTopic, e)
	return nil
}
