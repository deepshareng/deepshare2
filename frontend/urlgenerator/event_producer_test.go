package urlgenerator

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"net/url"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg/messaging"
)

func TestProduceGenURLEvent(t *testing.T) {
	q := new(bytes.Buffer)
	p := messaging.NewSimpleProducer(q)

	tests := []struct {
		appID      string
		rawUrl     string
		shortUrl   string
		genUrlPost *GenURLPostBody
		event      *messaging.Event
	}{
		{
			appID:    "AppID",
			rawUrl:   "rawUrl",
			shortUrl: "shortUrl",
			genUrlPost: &GenURLPostBody{
				InAppData:          "inappData",
				DownloadTitle:      "dTitle",
				DownloadBtnText:    "dBtnText",
				DownloadMsg:        "dMsg",
				DownloadUrlIos:     "dlUrlIos",
				DownloadUrlAndroid: "dlUrlAndroid",
				UninstallUrl:       "uninstallUrl",
				RedirectUrl:        "rUrl",
				IsShort:            true,
				SDKInfo:            "android2.0",
				SenderID:           "uniqueId1",
				Channels: []string{
					"c1",
					"c2",
				},
			},

			event: &messaging.Event{
				AppID:     "AppID",
				EventType: "/v2/url/",
				Channels: []string{
					"c1",
					"c2",
				},
				SenderID: "uniqueId1",
				Count:    1,
				UAInfo:   uainfo.UAInfo{},
				KVs: map[string]interface{}{
					"is_short":             true,
					"sdk_info":             "android2.0",
					"inapp_data":           "inappData",
					"download_title":       "dTitle",
					"download_btn_text":    "dBtnText",
					"download_msg":         "dMsg",
					"download_url_ios":     "dlUrlIos",
					"download_url_android": "dlUrlAndroid",
					"uninstall_url":        "uninstallUrl",
					"redirect_url":         "rUrl",
					"rawUrl":               "rawUrl",
					"shortUrl":             "shortUrl",
					"role":                 "sender",
				},
			},
		},
		{
			appID:    "AppID",
			rawUrl:   "rawUrl",
			shortUrl: "shortUrl",
			genUrlPost: &GenURLPostBody{
				InAppData:          "inappData",
				DownloadTitle:      "dTitle",
				DownloadBtnText:    "dBtnText",
				DownloadMsg:        "dMsg",
				DownloadUrlIos:     "dlUrlIos",
				DownloadUrlAndroid: "dlUrlAndroid",
				UninstallUrl:       "uninstallUrl",
				RedirectUrl:        "rUrl",
				IsShort:            true,
				SDKInfo:            "android2.0",
				ForwardedSenderID:  "uniqueId1",
				Channels: []string{
					"c1",
					"c2",
				},
			},

			event: &messaging.Event{
				AppID:     "AppID",
				EventType: "/v2/url/",
				Channels: []string{
					"c1",
					"c2",
				},
				SenderID: "uniqueId1",
				Count:    1,
				UAInfo:   uainfo.UAInfo{},
				KVs: map[string]interface{}{
					"is_short":             true,
					"sdk_info":             "android2.0",
					"inapp_data":           "inappData",
					"download_title":       "dTitle",
					"download_btn_text":    "dBtnText",
					"download_msg":         "dMsg",
					"download_url_ios":     "dlUrlIos",
					"download_url_android": "dlUrlAndroid",
					"uninstall_url":        "uninstallUrl",
					"redirect_url":         "rUrl",
					"rawUrl":               "rawUrl",
					"shortUrl":             "shortUrl",
					"role":                 "receiver",
				},
			},
		},
	}

	for i, tt := range tests {
		q.Reset()

		if err := produceGenerateUrlEvent(tt.appID, tt.rawUrl, tt.shortUrl, api.GenerateUrlPrefix, tt.genUrlPost, &uainfo.UAInfo{}, p); err != nil {
			t.Fatalf("#%d: produceDSActionEvent failed!\nerr:%v", i, err)
		}
		// Simple producer should produce events to "dsaction" topic
		de := json.NewDecoder(q)
		// first extract SimpleProducerEvent
		se := new(messaging.SimpleProducerEvent)
		if err := de.Decode(se); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(se.Topic, messaging.GenUrlTopic) {
			t.Errorf("#%d, topic=%s, want=%s", i, se.Topic, messaging.GenUrlTopic)
		}
		// then extract Event from SimpleProducerEvent#Msg
		ce := new(messaging.Event)
		if err := ce.Unmarshal(se.Msg); err != nil {
			t.Fatal(err)
		}
		// Empty the UAInfo to ensuring we could compare it with our origin data
		ce.UAInfo = uainfo.UAInfo{}
		if !reflect.DeepEqual(ce, tt.event) {
			t.Errorf("#%d, event=%v, want=%v", i, ce, tt.event)
		}
	}
}

func TestProduceGetInAppDataEvent(t *testing.T) {
	tests := []struct {
		appID        string
		endpoint     string
		shortSeg     string
		tracking     string
		rawUrlStr    string
		receiverInfo api.MatchReceiverInfo
		ua           *uainfo.UAInfo
		wEvent       *messaging.Event
	}{
		{
			appID:     "TestAppID",
			endpoint:  api.GenerateUrlPrefix,
			shortSeg:  "s1",
			tracking:  "open",
			rawUrlStr: `http://example.com/d/TestAppID?channels=channel1|channel2&download_msg=dMsg&download_title=dTitle&inapp_data={"k1":"v1","k2":"v2"}&redirect_url=http://www.baidu.com/aaa?abc=eee&sdk_info=android1.3.1&sender_id=aabbcc`,
			receiverInfo: api.MatchReceiverInfo{
				UniqueID:  "u1",
				OS:        "ios",
				OSVersion: "9.1",
			},
			ua: &uainfo.UAInfo{
				Ip:        "testip1",
				Os:        "ios",
				OsVersion: "9.1",
			},
			wEvent: &messaging.Event{
				AppID:     "TestAppID",
				EventType: api.GenerateUrlPrefix + "open",
				UniqueID:  "u1",
				Count:     1,
				Channels:  []string{"channel1", "channel2"},
				SenderID:  "aabbcc",
				UAInfo: uainfo.UAInfo{
					Ip:        "testip1",
					Os:        "ios",
					OsVersion: "9.1",
				},
				ReceiverInfo: api.MatchReceiverInfo{
					UniqueID:  "u1",
					OS:        "ios",
					OSVersion: "9.1",
				},
				KVs: map[string]interface{}{
					"short_seg":  "s1",
					"inapp_data": `{"k1":"v1","k2":"v2"}`,
				},
			},
		},
	}
	q := new(bytes.Buffer)
	p := messaging.NewSimpleProducer(q)
	for i, tt := range tests {
		q.Reset()
		rawUrl, err := url.Parse(tt.rawUrlStr)
		if err != nil {
			t.Fatal(err)
		}
		if err := produceGetInappDataEvent(tt.appID, tt.endpoint, tt.shortSeg, tt.tracking, rawUrl, tt.receiverInfo, tt.ua, p); err != nil {
			t.Fatal(err)
		}
		// Simple producer should produce events to "dsaction" topic
		de := json.NewDecoder(q)
		// first extract SimpleProducerEvent
		se := new(messaging.SimpleProducerEvent)
		if err := de.Decode(se); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(se.Topic, messaging.GenUrlTopic) {
			t.Errorf("#%d, topic=%s, want=%s\n", i, se.Topic, messaging.GenUrlTopic)
		}
		// then extract Event from SimpleProducerEvent#Msg
		ce := new(messaging.Event)
		if err := ce.Unmarshal(se.Msg); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(ce, tt.wEvent) {
			t.Errorf("#%d, event=%#v, want=%#v\n", i, ce, tt.wEvent)
		}
	}
}
