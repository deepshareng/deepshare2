package retention

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/MISingularity/deepshare2/pkg/testutil"
)

func TestRetentionInstall(t *testing.T) {
	cli := testutil.MustNewRedisClient("")
	p := testutil.NewTestUtilProducer()
	retentionService := NewRedisRetentionService(3, "3-day", cli, []byte("retention"), p)
	testcases := []struct {
		input                          []messaging.Event
		expectedCollectionRecordNumber int64
		expectedMatchRetentEvent       []messaging.Event
	}{
		{
			[]messaging.Event{
				messaging.Event{
					AppID:        "app3",
					EventType:    api.GetInAppDataPrefix + "install",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "4",
					Count:        1,
					UAInfo:       uainfo.UAInfo{Os: "android"},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:        "app3",
					EventType:    api.GetInAppDataPrefix + "install",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "4",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:        "app4",
					EventType:    api.GetInAppDataPrefix + "install",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "4",
					Count:        1,
					UAInfo:       uainfo.UAInfo{Os: "ios"},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:        "app3",
					EventType:    api.GetInAppDataPrefix + "install",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "5",
					Count:        1,
					UAInfo:       uainfo.UAInfo{Os: "other"},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2014, time.December, 29, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:        "",
					EventType:    api.GetInAppDataPrefix + "install",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "4",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:        "123",
					EventType:    api.GetInAppDataPrefix + "install",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
			},
			3,
			[]messaging.Event{
				messaging.Event{
					AppID:     "app3",
					EventType: api.RetentionPrefix + retentionService.name + "_install",
					UniqueID:  "4",
					Channels:  []string{},
					SenderID:  "",
					CookieID:  "",
					Count:     1,
					KVs:       map[string]interface{}{},
					UAInfo:    uainfo.UAInfo{Os: "android"},
					TimeStamp: time.Date(2015, time.January, 4, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:     "app4",
					EventType: api.RetentionPrefix + retentionService.name + "_install",
					UniqueID:  "4",
					Channels:  []string{},
					SenderID:  "",
					CookieID:  "",
					Count:     1,
					KVs:       map[string]interface{}{},
					UAInfo:    uainfo.UAInfo{Os: "ios"},
					TimeStamp: time.Date(2015, time.January, 4, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:     "app3",
					EventType: api.RetentionPrefix + retentionService.name + "_install",
					UniqueID:  "5",
					Channels:  []string{},
					SenderID:  "",
					CookieID:  "",
					Count:     1,
					UAInfo:    uainfo.UAInfo{Os: "other"},
					KVs:       map[string]interface{}{},
					TimeStamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
			},
		},
	}
	for i, tt := range testcases {
		retentionService.cli.FlushDb()
		p.Clear()
		for _, m := range tt.input {
			retentionService.InsertInstallEventForRetention(&m)
		}
		n := retentionService.cli.DbSize().Val()
		if n != tt.expectedCollectionRecordNumber {
			t.Errorf("#%d: Get unexpected record number, want = %d, get = %d.", i, tt.expectedCollectionRecordNumber, n)
		}
		for _, e := range tt.expectedMatchRetentEvent {
			res, err := json.Marshal(e)
			if err != nil {
				t.Error(err)
				continue
			}
			v := string(retentionService.produceTopic) + "#" + string(res)
			_, ok := p.C[v]
			if !ok {
				t.Errorf("Expected retention message is not exist! Expected Msg = %s", v)
				continue
			}
			p.C[v]--
			if p.C[v] == 0 {
				delete(p.C, v)
			}
		}
		if len(p.C) != 0 {
			t.Errorf("Remain message number mismatsh! Expected = %d, get = %d.", 0, len(p.C))
		}
	}
}

func TestRetentionFind(t *testing.T) {
	cli := testutil.MustNewRedisClient("")
	p := testutil.NewTestUtilProducer()
	retentionService := NewRedisRetentionService(3, "3-day", cli, []byte("retention"), p)
	testcases := []struct {
		input                          []messaging.Event
		expectedCollectionRecordNumber int64
		expectedMatchRetentEvent       []messaging.Event
		expectedInstallMessage         int
	}{
		{
			[]messaging.Event{
				messaging.Event{
					AppID:        "app1",
					EventType:    api.GetInAppDataPrefix + "install",
					Channels:     []string{"abddc", "ddd"},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "1",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:        "app1",
					EventType:    api.GetInAppDataPrefix + "open",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "1",
					Count:        1,
					UAInfo:       uainfo.UAInfo{Os: "android"},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 4, 0, 0, 0, 0, time.Local).Unix(),
				},
			},
			0,
			[]messaging.Event{
				messaging.Event{
					AppID:     "app1",
					EventType: api.RetentionPrefix + retentionService.name,
					UniqueID:  "1",
					Channels:  []string{"abddc", "ddd"},
					SenderID:  "",
					CookieID:  "",
					Count:     1,
					UAInfo:    uainfo.UAInfo{Os: "android"},
					KVs:       map[string]interface{}{},
					TimeStamp: time.Date(2015, time.January, 4, 0, 0, 0, 0, time.Local).Unix(),
				},
			},
			1,
		},
		{
			[]messaging.Event{
				messaging.Event{
					AppID:        "app2",
					EventType:    api.GetInAppDataPrefix + "install",
					Channels:     []string{"ddd"},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "2",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:        "app2",
					EventType:    api.GetInAppDataPrefix + "install",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "3",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:        "app2",
					EventType:    api.GetInAppDataPrefix + "open",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "3",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 4, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:        "app2",
					EventType:    api.GetInAppDataPrefix + "open",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "3",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 4, 0, 0, 0, 0, time.Local).Unix(),
				},
			},
			1,
			[]messaging.Event{
				messaging.Event{
					AppID:     "app2",
					EventType: api.RetentionPrefix + retentionService.name,
					UniqueID:  "3",
					Channels:  []string{""},
					SenderID:  "",
					CookieID:  "",
					Count:     1,
					KVs:       map[string]interface{}{},
					TimeStamp: time.Date(2015, time.January, 4, 0, 0, 0, 0, time.Local).Unix(),
				},
			},
			2,
		},
		{
			[]messaging.Event{
				messaging.Event{
					AppID:        "app3",
					EventType:    api.GetInAppDataPrefix + "install",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "4",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:        "app3",
					EventType:    api.GetInAppDataPrefix + "install",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "4",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:        "app3",
					EventType:    api.GetInAppDataPrefix + "install",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "4",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:        "app3",
					EventType:    api.GetInAppDataPrefix + "install",
					Channels:     []string{"w1", "w2", "w3"},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "5",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2014, time.December, 29, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:        "app3",
					EventType:    api.GetInAppDataPrefix + "open",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "87",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 4, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:        "app3",
					EventType:    api.GetInAppDataPrefix + "open",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "4",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 4, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:        "app3",
					EventType:    api.GetInAppDataPrefix + "open",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "4",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 3, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:        "app3",
					EventType:    api.GetInAppDataPrefix + "open",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "5",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
			},
			0,
			[]messaging.Event{
				messaging.Event{
					AppID:     "app3",
					EventType: api.RetentionPrefix + retentionService.name,
					UniqueID:  "4",
					Channels:  []string{""},
					SenderID:  "",
					CookieID:  "",
					Count:     1,
					KVs:       map[string]interface{}{},
					TimeStamp: time.Date(2015, time.January, 4, 0, 0, 0, 0, time.Local).Unix(),
				},
				messaging.Event{
					AppID:     "app3",
					EventType: api.RetentionPrefix + retentionService.name,
					UniqueID:  "5",
					Channels:  []string{"w1", "w2", "w3"},
					SenderID:  "",
					CookieID:  "",
					Count:     1,
					KVs:       map[string]interface{}{},
					TimeStamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				},
			},
			2,
		},
	}
	for i, tt := range testcases {
		retentionService.cli.FlushDb()
		p.Clear()
		for _, m := range tt.input {
			if strings.HasSuffix(m.EventType, "install") {
				retentionService.InsertInstallEventForRetention(&m)
			}
			if strings.HasSuffix(m.EventType, "open") {
				retentionService.FindMatchUserForRetention(&m)
			}
		}
		n := retentionService.cli.DbSize().Val()
		if n != tt.expectedCollectionRecordNumber {
			t.Errorf("#%d: Get unexpected record number, want = %d, get = %d.", i, tt.expectedCollectionRecordNumber, n)
		}
		for _, e := range tt.expectedMatchRetentEvent {
			res, err := json.Marshal(e)
			if err != nil {
				t.Error(err)
				continue
			}
			v := string(retentionService.produceTopic) + "#" + string(res)
			_, ok := p.C[v]
			if !ok {
				t.Errorf("Expected retention message is not exist! Expected Msg = %s", v)
				continue
			}
			p.C[v]--
			if p.C[v] == 0 {
				delete(p.C, v)
			}
		}
		for k, v := range p.C {
			t.Log(k, v)
		}
		if len(p.C) != tt.expectedInstallMessage {
			t.Errorf("Remain message number mismatsh! Expected = %d, get = %d.", tt.expectedInstallMessage, len(p.C))
		}
	}
}
