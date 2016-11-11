package dsaction

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg/messaging"
)

func TestProduceDSActionEvent(t *testing.T) {
	q := new(bytes.Buffer)
	p := messaging.NewSimpleProducer(q)

	tests := []struct {
		appID   string
		request *DSActionRequest
		event   *messaging.Event
	}{
		{
			appID: "AppID",
			request: &DSActionRequest{
				Action:       "app/close",
				ReceiverInfo: api.ReceiverInfo{UniqueID: "d1"},
			},

			event: &messaging.Event{
				AppID:     "AppID",
				UniqueID:  "d1",
				EventType: "/dsactions/app/close",
				Count:     1,
				UAInfo:    uainfo.UAInfo{},
			},
		},
		{
			appID: "AppID",
			request: &DSActionRequest{
				Action: "js/dst",
				KVs: map[string]interface{}{
					"click_id":    "clickId1",
					"destination": "dst1",
				},
			},

			event: &messaging.Event{
				AppID:     "AppID",
				EventType: "/dsactions/js/dst",
				Count:     1,
				UAInfo:    uainfo.UAInfo{},
				KVs: map[string]interface{}{
					"click_id":    "clickId1",
					"destination": "dst1",
				},
			},
		},
	}

	for i, tt := range tests {
		q.Reset()

		if err := produceDSActionEvent(tt.appID, tt.request, &uainfo.UAInfo{}, p, "/dsactions/"); err != nil {
			t.Fatalf("#%d: produceDSActionEvent failed!\nerr:%v", i, err)
		}
		// Simple producer should produce events to "dsaction" topic
		de := json.NewDecoder(q)
		// first extract SimpleProducerEvent
		se := new(messaging.SimpleProducerEvent)
		if err := de.Decode(se); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(se.Topic, messaging.DSActionTopic) {
			t.Errorf("#%d, topic=%s, want=%s", i, se.Topic, messaging.DSActionTopic)
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
