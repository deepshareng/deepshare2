package sharelink

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg/messaging"
)

func TestProduceSharelinkEvent(t *testing.T) {
	q := new(bytes.Buffer)
	p := messaging.NewSimpleProducer(q)
	tests := []struct {
		appID      string
		endpoint   string
		channels   string
		senderID   string
		cookieID   string
		ip         string
		ua         string
		shortToken string
		isyyb      bool
		isUL       bool
		dsTag      string

		event *messaging.Event
	}{
		{
			appID:      "AppID",
			endpoint:   "d",
			channels:   "channel1|channel2",
			senderID:   "sender1",
			cookieID:   "cookie1",
			ip:         "ip1",
			ua:         "ua1",
			shortToken: "shorttoken1",
			isyyb:      true,
			isUL:       true,
			dsTag:      "t1",

			event: &messaging.Event{
				AppID:     "AppID",
				EventType: EventTypePrefix + "d",
				Channels:  []string{"channel1", "channel2"},
				SenderID:  "sender1",
				CookieID:  "cookie1",
				Count:     1,
				UAInfo: uainfo.UAInfo{
					Ua:          "ua1",
					Ip:          "ip1",
					Os:          "",
					Brand:       "-",
					Browser:     "Other",
					ChromeMajor: 0,
				},
				KVs: map[string]interface{}{
					"shorturl_token":       "shorttoken1",
					"shorturl_token_valid": true,
					"is_yyb":               true,
					"is_universallink":     true,
					"ds_tag":               "t1",
				},
			},
		},
	}
	for i, tt := range tests {
		q.Reset()

		e := &messaging.Event{}

		e = AddEventShortUrlToken(e, tt.shortToken)
		e = AddEventInfo(e, tt.appID, tt.endpoint, tt.channels, tt.senderID, tt.cookieID, tt.ip, tt.ua, tt.dsTag)
		e = AddEventOfAppInfo(e, tt.isyyb, tt.isUL)
		if err := FireShareLinkEvent(e, p); err != nil {
			t.Fatalf("#%d: Produce Share Link Events failed!\nerr:%v", i, err)
		}

		de := json.NewDecoder(q)
		se := new(messaging.SimpleProducerEvent)
		if err := de.Decode(se); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(se.Topic, messaging.ShareLinkTopic) {
			t.Errorf("#%d, topic=%s, want=%s", i, se.Topic, messaging.CounterTopic)
		}
		ce := new(messaging.Event)
		if err := ce.Unmarshal(se.Msg); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(ce, tt.event) {
			t.Errorf("#%d, event=%v, want=%v", i, ce, tt.event)
		}
	}
}
