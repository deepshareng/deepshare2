package aggregate

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/deepstats/aggregate"
	"github.com/MISingularity/deepshare2/pkg/condition"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/nsqio/go-nsq"
)

type expectMsg aggregate.CounterEvent

type aggregateTestService struct {
	st map[expectMsg]int
}

func (a *aggregateTestService) Insert(appID string, aggregate aggregate.CounterEvent) error {
	a.st[expectMsg(aggregate)]++
	return nil
}

func (a *aggregateTestService) QueryDuration(appid string, channel string, eventFilters []string, start time.Time, granularity time.Duration, limit int, os string) ([]*aggregate.AggregateResult, error) {
	return []*aggregate.AggregateResult{}, nil
}

func (a *aggregateTestService) Aggregate(appID string) error {
	return nil
}

func (a *aggregateTestService) StartRefreshLoop() {

}

func (a *aggregateTestService) ConvertTimeToGranularity(time.Time) time.Time {
	return time.Now()
}

func TestAppEventConsume(t *testing.T) {

	testcases := []struct {
		input  []messaging.Event
		result []expectMsg
	}{
		//sharelink test
		{
			[]messaging.Event{
				messaging.Event{
					AppID:        "app1",
					EventType:    "sharelink:/d/",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"shorturl_token": "t1", "ds_tag": "abc", "shorturl_token_valid": true},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:        "app1",
					EventType:    "sharelink:/d/",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"shorturl_token": "t2", "shorturl_token_valid": true},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:        "app1",
					EventType:    "sharelink:/d/",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"shorturl_token": "t3", "shorturl_token_valid": true},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},

				// without shorturl_token, event should not be converted
				messaging.Event{
					AppID:        "app1",
					EventType:    "sharelink:/d/",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
			},

			[]expectMsg{
				expectMsg(aggregate.CounterEvent{
					AppID:     "app1",
					Channel:   "all",
					Event:     "sharelink:/d/_ds_tag",
					Count:     1,
					Timestamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
					Os:        "other",
				}),
				expectMsg(aggregate.CounterEvent{
					AppID:     "app1",
					Channel:   "all",
					Event:     "sharelink:/d/_no_ds_tag",
					Count:     1,
					Timestamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
					Os:        "other",
				}),
				expectMsg(aggregate.CounterEvent{
					AppID:     "app1",
					Channel:   "all",
					Event:     "sharelink:/d/_no_ds_tag",
					Count:     1,
					Timestamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
					Os:        "other",
				}),
				expectMsg(aggregate.CounterEvent{
					AppID:     "app1",
					Channel:   "all",
					Event:     "sharelink:/d/",
					Count:     1,
					Timestamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
					Os:        "other",
				}),
			},
		},
		// bydeepshare/inappdate test
		{
			[]messaging.Event{
				messaging.Event{
					AppID:        "app1",
					EventType:    api.CounterPrefix + "install",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "",
					Count:        1,
					UAInfo:       uainfo.UAInfo{Os: "ios"},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
			},
			[]expectMsg{
				expectMsg(aggregate.CounterEvent{
					AppID:     "app1",
					Channel:   "all",
					Event:     "/v2/counters/install",
					Count:     1,
					Timestamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
					Os:        "ios",
				}),
				expectMsg(aggregate.CounterEvent{
					AppID:     "app1",
					Channel:   "all",
					Event:     "/v2/counters/install_with_inappdata",
					Count:     1,
					Timestamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
					Os:        "ios",
				}),
				expectMsg(aggregate.CounterEvent{
					AppID:     "app1",
					Channel:   "all",
					Event:     "/v2/counters/install_with_params",
					Count:     1,
					Timestamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
					Os:        "ios",
				}),
			},
		},

		{
			[]messaging.Event{
				messaging.Event{
					AppID:        "app2",
					EventType:    "install",
					Channels:     []string{"a", "b"},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "",
					Count:        1,
					UAInfo:       uainfo.UAInfo{Os: "android"},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:        "app3",
					EventType:    api.GetInAppDataPrefix + "open",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "",
					Count:        1,
					UAInfo:       uainfo.UAInfo{Brand: "3", Os: "android"},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:        "app3",
					EventType:    api.GetInAppDataPrefix + "open",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "",
					Count:        1,
					UAInfo:       uainfo.UAInfo{Brand: "3", Os: "android"},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
			},
			[]expectMsg{
				// 1st msg
				expectMsg(aggregate.CounterEvent{
					AppID:     "app2",
					Channel:   "a",
					Event:     "install",
					Count:     1,
					Timestamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
					Os:        "android",
				}),

				expectMsg(aggregate.CounterEvent{
					AppID:     "app2",
					Channel:   "a",
					Event:     "install_with_params",
					Count:     1,
					Timestamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
					Os:        "android",
				}),
				expectMsg(aggregate.CounterEvent{
					AppID:     "app2",
					Channel:   "b",
					Event:     "install",
					Count:     1,
					Timestamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
					Os:        "android",
				}),

				expectMsg(aggregate.CounterEvent{
					AppID:     "app2",
					Channel:   "b",
					Event:     "install_with_params",
					Count:     1,
					Timestamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
					Os:        "android",
				}),
				expectMsg(aggregate.CounterEvent{
					AppID:     "app2",
					Channel:   "all",
					Event:     "install",
					Count:     1,
					Timestamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
					Os:        "android",
				}),
				expectMsg(aggregate.CounterEvent{
					AppID:     "app2",
					Channel:   "all",
					Event:     "install_with_params",
					Count:     1,
					Timestamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
					Os:        "android",
				}),

				// 2nd msg
				expectMsg(aggregate.CounterEvent{
					AppID:     "app3",
					Channel:   "all",
					Event:     api.GetInAppDataPrefix + "open",
					Count:     1,
					Timestamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
					Os:        "android",
				}),
				expectMsg(aggregate.CounterEvent{
					AppID:     "app3",
					Channel:   "all",
					Event:     condition.DisplayPrefix + "open",
					Count:     1,
					Timestamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
					Os:        "android",
				}),

				// 3rd msg
				expectMsg(aggregate.CounterEvent{
					AppID:     "app3",
					Channel:   "all",
					Event:     api.GetInAppDataPrefix + "open",
					Count:     1,
					Timestamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
					Os:        "android",
				}),
				expectMsg(aggregate.CounterEvent{
					AppID:     "app3",
					Channel:   "all",
					Event:     condition.DisplayPrefix + "open",
					Count:     1,
					Timestamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
					Os:        "android",
				}),
			},
		},
		{
			[]messaging.Event{
				messaging.Event{
					AppID:        "app3",
					EventType:    api.GetInAppDataPrefix + "open",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "",
					Count:        1,
					UAInfo:       uainfo.UAInfo{Brand: "3", Os: "android"},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "", "ua": "UA:127.0.0.1_ios_7.0.3_-", "short_seg": ""},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:        "app1",
					EventType:    api.GetInAppDataPrefix + "open",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "",
					Count:        1,
					UAInfo:       uainfo.UAInfo{Brand: "3", Os: "android"},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "", "ua": "UA:127.0.0.1_ios_7.0.3_-", "short_seg": "123"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
			},
			[]expectMsg{
				expectMsg(aggregate.CounterEvent{
					AppID:     "app3",
					Channel:   "all",
					Event:     api.GetInAppDataPrefix + "open",
					Count:     1,
					Timestamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
					Os:        "android",
				}),
				expectMsg(aggregate.CounterEvent{
					AppID:     "app3",
					Channel:   "all",
					Event:     condition.DisplayPrefix + "open",
					Count:     1,
					Timestamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
					Os:        "android",
				}),
				expectMsg(aggregate.CounterEvent{
					AppID:     "app1",
					Channel:   "all",
					Event:     api.GetInAppDataPrefix + "open",
					Count:     1,
					Timestamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
					Os:        "android",
				}),
				expectMsg(aggregate.CounterEvent{
					AppID:     "app1",
					Channel:   "all",
					Event:     condition.DisplayPrefix + "open",
					Count:     1,
					Timestamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
					Os:        "android",
				}),
				expectMsg(aggregate.CounterEvent{
					AppID:     "app1",
					Channel:   "all",
					Event:     api.GenerateUrlPrefix,
					Count:     1,
					Timestamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
					Os:        "android",
				}),
			},
		},
	}

	for _, tt := range testcases {
		ats := aggregateTestService{make(map[expectMsg]int)}
		appconsumer := AggregateConsumer{&ats}
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
			t.Errorf("Message number mismatch! Details=%v", ats.st)
		}

	}
}
