package inappdata

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"errors"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg/messaging"
)

func TestProduceInAppDataFromRawUrlEvent(t *testing.T) {
	tests := []struct {
		appID        string
		endpoint     string
		shortSeg     string
		tracking     string
		inAppData    string
		channels     []string
		senderId     string
		ip           string
		ua           string
		receiverInfo api.MatchReceiverInfo
		wEvent       *messaging.Event
	}{
		{
			appID:     "TestAppID",
			endpoint:  api.GetInAppDataPrefix,
			shortSeg:  "s1",
			tracking:  "open",
			inAppData: `{"k1":"v1","k2":"v2"}`,
			channels:  []string{"channel1", "channel2"},
			senderId:  "sender1",
			ip:        "ip1",
			ua:        "ua1",
			receiverInfo: api.MatchReceiverInfo{
				UniqueID:  "u1",
				OS:        "iOS",
				OSVersion: "9.1",
			},
			wEvent: &messaging.Event{
				AppID:     "TestAppID",
				EventType: api.GetInAppDataPrefix + "open",
				UniqueID:  "u1",
				Count:     1,
				Channels:  []string{"channel1", "channel2"},
				SenderID:  "sender1",
				UAInfo: uainfo.UAInfo{
					Ua:          "ua1",
					Ip:          "ip1",
					Os:          "",
					Brand:       "-",
					Browser:     "Other",
					ChromeMajor: 0,
				},
				ReceiverInfo: api.MatchReceiverInfo{
					UniqueID:  "u1",
					OS:        "iOS",
					OSVersion: "9.1",
				},
				KVs: map[string]interface{}{
					"type":       "fromRawUrl",
					"short_seg":  "s1",
					"inapp_data": `{"k1":"v1","k2":"v2"}`,
				},
			},
		},
		{
			appID:     "TestAppID",
			endpoint:  api.GetInAppDataPrefix,
			shortSeg:  "s1",
			tracking:  "install",
			inAppData: `{"k1":"v1","k2":"v2"}`,
			channels:  []string{"channel1", "channel2"},
			senderId:  "sender1",
			ip:        "ip1",
			ua:        "ua1",
			receiverInfo: api.MatchReceiverInfo{
				UniqueID:  "u1",
				OS:        "iOS",
				OSVersion: "9.1",
			},
			wEvent: &messaging.Event{
				AppID:     "TestAppID",
				EventType: api.GetInAppDataPrefix + "install",
				UniqueID:  "u1",
				Count:     1,
				Channels:  []string{"channel1", "channel2"},
				SenderID:  "sender1",
				UAInfo: uainfo.UAInfo{
					Ua:          "ua1",
					Ip:          "ip1",
					Os:          "",
					Brand:       "-",
					Browser:     "Other",
					ChromeMajor: 0,
				},
				ReceiverInfo: api.MatchReceiverInfo{
					UniqueID:  "u1",
					OS:        "iOS",
					OSVersion: "9.1",
				},
				KVs: map[string]interface{}{
					"inapp_data": `{"k1":"v1","k2":"v2"}`,
					"short_seg":  "s1",
					"type":       "fromRawUrl",
				},
			},
		},
	}
	q := new(bytes.Buffer)
	p := messaging.NewSimpleProducer(q)
	for i, tt := range tests {
		q.Reset()

		produceInAppDataFromRawUrlEvent(tt.appID, tt.endpoint, tt.shortSeg, tt.tracking, tt.inAppData, tt.senderId, tt.ip, tt.ua, tt.channels, tt.receiverInfo, p)
		// Simple producer should produce events to "dsaction" topic
		de := json.NewDecoder(q)
		// first extract SimpleProducerEvent
		se := new(messaging.SimpleProducerEvent)
		if err := de.Decode(se); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(se.Topic, messaging.InAppDataTopic) {
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

func TestProduceInAppDataFromMatchEvent(t *testing.T) {
	tests := []struct {
		appID        string
		endpoint     string
		tracking     string
		inAppData    string
		channels     []string
		senderId     string
		ip           string
		ua           string
		clickId      string
		cookieId     string
		receiverInfo api.MatchReceiverInfo
		wEvent       *messaging.Event
	}{
		{
			appID:     "TestAppID",
			endpoint:  api.GetInAppDataPrefix,
			tracking:  "open",
			inAppData: `{"k1":"v1","k2":"v2"}`,
			channels:  []string{"channel1", "channel2"},
			senderId:  "sender1",
			ip:        "ip1",
			ua:        "ua1",
			clickId:   "click1",
			cookieId:  "cookie1",
			receiverInfo: api.MatchReceiverInfo{
				UniqueID:  "u1",
				OS:        "ios",
				OSVersion: "9.1",
			},
			wEvent: &messaging.Event{
				AppID:     "TestAppID",
				EventType: api.GetInAppDataPrefix + "open",
				UniqueID:  "u1",
				Count:     1,
				Channels:  []string{"channel1", "channel2"},
				SenderID:  "sender1",
				UAInfo: uainfo.UAInfo{
					Ua:          "ua1",
					Ip:          "ip1",
					Os:          "",
					Brand:       "-",
					Browser:     "Other",
					ChromeMajor: 0,
				},
				ReceiverInfo: api.MatchReceiverInfo{
					UniqueID:  "u1",
					OS:        "ios",
					OSVersion: "9.1",
				},
				KVs: map[string]interface{}{
					"type":       "fromMatch",
					"inapp_data": `{"k1":"v1","k2":"v2"}`,
					"click_id":   "click1",
					"cookie_id":  "cookie1",
				},
			},
		},
		{
			appID:     "TestAppID",
			endpoint:  api.GetInAppDataPrefix,
			tracking:  "install",
			inAppData: `{"k1":"v1","k2":"v2"}`,
			channels:  []string{"channel1", "channel2"},
			senderId:  "sender1",
			ip:        "ip1",
			ua:        "ua1",
			clickId:   "click1",
			cookieId:  "cookie1",
			receiverInfo: api.MatchReceiverInfo{
				UniqueID:  "u1",
				OS:        "ios",
				OSVersion: "9.1",
			},
			wEvent: &messaging.Event{
				AppID:     "TestAppID",
				EventType: api.GetInAppDataPrefix + "install",
				UniqueID:  "u1",
				Count:     1,
				Channels:  []string{"channel1", "channel2"},
				SenderID:  "sender1",
				UAInfo: uainfo.UAInfo{
					Ua:          "ua1",
					Ip:          "ip1",
					Os:          "",
					Brand:       "-",
					Browser:     "Other",
					ChromeMajor: 0,
				},
				ReceiverInfo: api.MatchReceiverInfo{
					UniqueID:  "u1",
					OS:        "ios",
					OSVersion: "9.1",
				},
				KVs: map[string]interface{}{
					"type":       "fromMatch",
					"inapp_data": `{"k1":"v1","k2":"v2"}`,
					"click_id":   "click1",
					"cookie_id":  "cookie1",
				},
			},
		},
	}
	q := new(bytes.Buffer)
	p := messaging.NewSimpleProducer(q)
	for i, tt := range tests {
		q.Reset()
		produceInAppDataFromMatchEvent(tt.appID, tt.endpoint, tt.tracking, tt.inAppData, tt.senderId, tt.ip, tt.ua, tt.clickId, tt.cookieId, tt.channels, tt.receiverInfo, p)
		// Simple producer should produce events to "dsaction" topic
		de := json.NewDecoder(q)
		// first extract SimpleProducerEvent
		se := new(messaging.SimpleProducerEvent)
		if err := de.Decode(se); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(se.Topic, messaging.InAppDataTopic) {
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

func TestProduceInAppDataErrorEvent(t *testing.T) {
	tests := []struct {
		appID        string
		endpoint     string
		tracking     string
		ip           string
		ua           string
		err          error
		receiverInfo api.MatchReceiverInfo
		wEvent       *messaging.Event
	}{
		{
			appID:    "TestAppID",
			endpoint: api.GetInAppDataPrefix,
			tracking: "open",
			ip:       "ip1",
			ua:       "ua1",
			err:      errors.New("error1"),
			receiverInfo: api.MatchReceiverInfo{
				UniqueID:  "u1",
				OS:        "ios",
				OSVersion: "9.1",
			},
			wEvent: &messaging.Event{
				AppID:     "TestAppID",
				EventType: api.GetInAppDataPrefix + "open",
				UniqueID:  "u1",
				Count:     1,
				UAInfo: uainfo.UAInfo{
					Ua:          "ua1",
					Ip:          "ip1",
					Os:          "",
					Brand:       "-",
					Browser:     "Other",
					ChromeMajor: 0,
				},
				ReceiverInfo: api.MatchReceiverInfo{
					UniqueID:  "u1",
					OS:        "ios",
					OSVersion: "9.1",
				},
				KVs: map[string]interface{}{
					"error": "error1",
				},
			},
		},
		{
			appID:    "TestAppID",
			endpoint: api.GetInAppDataPrefix,
			tracking: "install",
			ip:       "ip1",
			ua:       "ua1",
			err:      errors.New("error1"),
			receiverInfo: api.MatchReceiverInfo{
				UniqueID:  "u1",
				OS:        "ios",
				OSVersion: "9.1",
			},
			wEvent: &messaging.Event{
				AppID:     "TestAppID",
				EventType: api.GetInAppDataPrefix + "install",
				UniqueID:  "u1",
				Count:     1,
				UAInfo: uainfo.UAInfo{
					Ua:          "ua1",
					Ip:          "ip1",
					Os:          "",
					Brand:       "-",
					Browser:     "Other",
					ChromeMajor: 0,
				},
				ReceiverInfo: api.MatchReceiverInfo{
					UniqueID:  "u1",
					OS:        "ios",
					OSVersion: "9.1",
				},
				KVs: map[string]interface{}{
					"error": "error1",
				},
			},
		},
	}
	q := new(bytes.Buffer)
	p := messaging.NewSimpleProducer(q)
	for i, tt := range tests {
		q.Reset()
		produceInAppDataErrorEvent(tt.appID, tt.endpoint, tt.tracking, tt.ip, tt.ua, tt.receiverInfo, tt.err, p)
		// Simple producer should produce events to "dsaction" topic
		de := json.NewDecoder(q)
		// first extract SimpleProducerEvent
		se := new(messaging.SimpleProducerEvent)
		if err := de.Decode(se); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(se.Topic, messaging.InAppDataTopic) {
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
