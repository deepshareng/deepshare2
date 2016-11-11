package appchannel

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/deepstats/appchannel"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/nsqio/go-nsq"
)

type expectMsg string

type appchannelTestService struct {
	st map[expectMsg]int
}

func (a *appchannelTestService) GetChannels(appID string) (appchannel.AppChannels, error) {
	return appchannel.AppChannels{}, nil
}

func (a *appchannelTestService) InsertChannel(appID, channel string) (bool, error) {
	a.st[expectMsg(appID+"#"+channel)]++
	return true, nil
}

func (a *appchannelTestService) DeleteChannel(appID, channel string) error {
	return nil
}

func TestAppEventConsume(t *testing.T) {
	testcases := []struct {
		input  []messaging.Event
		result []expectMsg
	}{
		{
			[]messaging.Event{
				messaging.Event{
					AppID:        "app1",
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
					AppID:        "app2",
					EventType:    api.CounterPrefix + "install",
					Channels:     []string{"a", "b"},
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
					AppID:        "app1",
					EventType:    api.CounterPrefix + "open",
					Channels:     []string{"a", "c"},
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
					AppID:        "app3",
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
				// 1st msg
				"app1#a",
				// 2nd msg
				"app2#a",
				"app2#b",
				// 3rd msg
				"app1#a",
				"app1#c",
				// 4th msg
				"app3#a",
			},
		},
	}

	for _, tt := range testcases {
		ats := appchannelTestService{make(map[expectMsg]int)}
		appconsumer := AppchannelConsumer{&ats}
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
