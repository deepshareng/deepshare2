package dsaction

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

func TestDSActionAppPost(t *testing.T) {
	q := new(bytes.Buffer)
	p := messaging.NewSimpleProducer(q)
	handler := NewDSActionHandler(p, api.DSActionsPrefix)
	url := "http://" + path.Join("example.com", api.DSActionsPrefix, "AppID")
	tests := []struct {
		code   int
		header map[string]string
		body   string
		event  *messaging.Event
	}{
		{
			code:   http.StatusOK,
			header: map[string]string{},
			body:   `{"receiver_info":{"unique_id":"rec_dev"},"action":"app/close"}` + "\n",
			event: &messaging.Event{
				AppID:     "AppID",
				UniqueID:  "rec_dev",
				EventType: api.DSActionsPrefix + "app/close",
				Count:     1,
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
		// first extract SimpleProducerEvent
		se := new(messaging.SimpleProducerEvent)
		if err := de.Decode(se); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(se.Topic, messaging.DSActionTopic) {
			t.Errorf("#%d, topic=%s, want=%s", i, se.Topic, messaging.DSActionTopic)
		}
		// then extract DSAction event from SimpleProducerEvent#Msg
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

func TestDSActionJSPost(t *testing.T) {
	q := new(bytes.Buffer)
	p := messaging.NewSimpleProducer(q)
	handler := NewDSActionHandler(p, api.DSActionsPrefix)
	url := "http://" + path.Join("example.com", api.DSActionsPrefix, "AppID")
	tests := []struct {
		code   int
		header map[string]string
		body   string
		event  *messaging.Event
	}{
		{
			code:   http.StatusOK,
			header: map[string]string{},
			body:   `{"action":"js/dst","kvs":{"click_id":"click1","destination":"dst-weixin-tip-ios"}}` + "\n",

			event: &messaging.Event{
				AppID:     "AppID",
				EventType: api.DSActionsPrefix + "js/dst",
				Count:     1,
				KVs: map[string]interface{}{
					"click_id":    "click1",
					"destination": "dst-weixin-tip-ios",
				},
			},
		},
		{
			code:   http.StatusOK,
			header: map[string]string{},
			body:   `{"action":"js/deeplink","kvs":{"click_id":"click2","destination":"deepshare://open"}}` + "\n",
			event: &messaging.Event{
				AppID:     "AppID",
				EventType: api.DSActionsPrefix + "js/deeplink",
				Count:     1,
				KVs: map[string]interface{}{
					"click_id":    "click2",
					"destination": "deepshare://open",
				},
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
		// first extract SimpleProducerEvent
		se := new(messaging.SimpleProducerEvent)
		if err := de.Decode(se); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(se.Topic, messaging.DSActionTopic) {
			t.Errorf("#%d, topic=%s, want=%s", i, se.Topic, messaging.DSActionTopic)
		}
		// then extract DSAction event from SimpleProducerEvent#Msg
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
