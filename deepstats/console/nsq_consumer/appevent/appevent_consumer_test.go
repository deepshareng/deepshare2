package appevent

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/deepstats/appevent"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/nsqio/go-nsq"
)

type expectMsg string

type appeventTestService struct {
	st map[expectMsg]int
}

func (a *appeventTestService) GetEvents(appID string) (appevent.AppEvents, error) {
	return appevent.AppEvents{}, nil
}

func (a *appeventTestService) InsertEvent(appID, event string) (bool, error) {
	a.st[expectMsg(appID+"#"+event)]++
	return true, nil
}

func TestAppEventConsume(t *testing.T) {
	testcases := []struct {
		input  []messaging.Event
		result []expectMsg
	}{
		{
			[]messaging.Event{
				messaging.Event{
					AppID:        "C27A897D12F465AA",
					EventType:    api.CounterPrefix + "install",
					Channels:     []string{"a"},
					SenderID:     "a",
					CookieID:     "",
					UniqueID:     "",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
				messaging.Event{
					AppID:        "C27A897D12F465AA",
					EventType:    api.CounterPrefix + "install",
					Channels:     []string{"a"},
					SenderID:     "a",
					CookieID:     "",
					UniqueID:     "",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
				messaging.Event{
					AppID:        "C27A897D12F465AA",
					EventType:    api.CounterPrefix + "open",
					Channels:     []string{"a"},
					SenderID:     "a",
					CookieID:     "",
					UniqueID:     "",
					Count:        1,
					UAInfo:       uainfo.UAInfo{Brand: "3"},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
				messaging.Event{
					AppID:        "C27A897D12F465AA",
					EventType:    "match/open",
					Channels:     []string{"a"},
					SenderID:     "a",
					CookieID:     "",
					UniqueID:     "",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
			},
			[]expectMsg{
				"C27A897D12F465AA#/v2/counters/install",
				"C27A897D12F465AA#/v2/counters/install",
				"C27A897D12F465AA#/v2/counters/open",
			},
		},
	}

	for _, tt := range testcases {
		ats := appeventTestService{make(map[expectMsg]int)}
		appconsumer := AppeventConsumer{&ats}
		for _, m := range tt.input {
			d, err := json.Marshal(m)
			if err != nil {
				t.Errorf("%v", err)
				continue
			}
			msg := &nsq.Message{
				Body: []byte(d),
			}
			err = appconsumer.Consume(msg)
			if err != nil {
				t.Errorf("Consumer Msg failed! Msg Detail=%v, Err Msg=%v", m, err)
				continue
			}
		}
		for _, m := range tt.result {
			_, ok := ats.st[m]
			if !ok {
				t.Errorf("Expected message don't exist! Message Detail: %v", m)
			}
			ats.st[m]--
			if ats.st[m] == 0 {
				delete(ats.st, m)
			}
		}
		if len(ats.st) != 0 {
			t.Errorf("Message number mismatch!")
		}

	}
}
