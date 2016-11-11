package counter

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg/messaging"
)

func TestProduceCounterEvent(t *testing.T) {
	q := new(bytes.Buffer)
	p := messaging.NewSimpleProducer(q)

	tests := []struct {
		appID   string
		request *CounterRequest
		events  []*messaging.Event
	}{
		{
			appID: "AppID",
			request: &CounterRequest{
				ReceiverInfo: api.ReceiverInfo{UniqueID: "d1"},
				Counters: []*Counter{
					{
						Event: "e1",
						Count: 1,
					},
				},
			},

			events: []*messaging.Event{
				{
					AppID:     "AppID",
					UniqueID:  "d1",
					EventType: "/counters/e1",
					Count:     1,
					UAInfo:    uainfo.UAInfo{},
				},
			},
		},
		{
			appID: "AppID",
			request: &CounterRequest{
				ReceiverInfo: api.ReceiverInfo{UniqueID: "d1"},
				Counters: []*Counter{
					{
						Event: "e1",
						Count: 1,
					},
					{
						Event: "e2",
						Count: 2,
					},
				},
			},
			events: []*messaging.Event{
				{
					AppID:     "AppID",
					UniqueID:  "d1",
					EventType: "/counters/e1",
					Count:     1,
					UAInfo:    uainfo.UAInfo{},
				},
				{
					AppID:     "AppID",
					UniqueID:  "d1",
					EventType: "/counters/e2",
					Count:     2,
					UAInfo:    uainfo.UAInfo{},
				},
			},
		},
	}

	for i, tt := range tests {
		q.Reset()

		if err := produceCounterEvents(tt.appID, tt.request, &uainfo.UAInfo{}, p, "/counters/"); err != nil {
			t.Fatalf("#%d: ProduceCounterEvents failed!\nerr:%v", i, err)
		}
		// Simple producer should produce events to "counter" topic
		de := json.NewDecoder(q)
		for _, event := range tt.events {
			// first extract SimpleProducerEvent
			se := new(messaging.SimpleProducerEvent)
			if err := de.Decode(se); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(se.Topic, messaging.CounterTopic) {
				t.Errorf("#%d, topic=%s, want=%s", i, se.Topic, messaging.CounterTopic)
			}
			// then extract CounterEvent from SimpleProducerEvent#Msg
			ce := new(messaging.Event)
			if err := ce.Unmarshal(se.Msg); err != nil {
				t.Fatal(err)
			}
			// Empty the UAInfo to ensuring we could compare it with our origin data
			ce.UAInfo = uainfo.UAInfo{}
			if !reflect.DeepEqual(ce, event) {
				t.Errorf("#%d, event=%v, want=%v", i, ce, event)
			}
		}
	}
}
