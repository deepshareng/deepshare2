package counter

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path"
	"reflect"
	"testing"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/MISingularity/deepshare2/pkg/testutil"
)

func TestCounterPost(t *testing.T) {
	q := new(bytes.Buffer)
	p := messaging.NewSimpleProducer(q)
	handler := NewCounterHandler(p, api.CounterPrefix)
	url := "http://" + path.Join("example.com", api.CounterPrefix, "AppID")
	tests := []struct {
		code   int
		header map[string]string
		body   string
		events []*messaging.Event
	}{
		{
			http.StatusOK,
			map[string]string{},
			`{"receiver_info":{"unique_id":"rec_dev"},` +
				`"counters":[{"event":"e1","count":1},{"event":"e2","count":10}]}` + "\n",
			[]*messaging.Event{
				{AppID: "AppID", UniqueID: "rec_dev", EventType: api.CounterPrefix + "e1", Count: 1},
				{AppID: "AppID", UniqueID: "rec_dev", EventType: api.CounterPrefix + "e2", Count: 10},
			},
		},
	}

	for i, tt := range tests {
		q.Reset()
		w := testutil.HandleWithRequestInfo(handler, "POST", url, tt.body, tt.header, "ip1:port1")
		if w.Code != tt.code {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.code)
		}

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
