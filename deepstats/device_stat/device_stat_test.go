package device_stat

import (
	"fmt"
	"testing"
	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/MISingularity/deepshare2/pkg/storage"
)

func TestDeviceStatInsert(t *testing.T) {
	ds := NewGeneralDeviceStat(storage.NewInMemSimpleKV(), "testdeviceinsert")
	testcase := []messaging.Event{
		messaging.Event{
			AppID:        "C27A897D12F465AA",
			EventType:    api.CounterPrefix + "install",
			Channels:     []string{"a"},
			SenderID:     "a",
			CookieID:     "",
			UniqueID:     "u1",
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
			UniqueID:     "u2",
			Count:        1,
			UAInfo:       uainfo.UAInfo{Os: "android"},
			ReceiverInfo: api.MatchReceiverInfo{},
			KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
			TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC).Unix(),
		},
	}
	for _, v := range testcase {
		if err := ds.Insert(&v); err != nil {
			t.Fatal(err)
		}
	}
}

func TestDeviceStatCount(t *testing.T) {
	testcase := []struct {
		insertDevice []messaging.Event
		os           string
		count        int64
	}{
		{
			[]messaging.Event{
				messaging.Event{
					AppID:        "C27A897D12F465AA",
					EventType:    api.CounterPrefix + "install",
					Channels:     []string{"a"},
					SenderID:     "u1",
					CookieID:     "",
					UniqueID:     "u1",
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
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "u2",
					Count:        1,
					UAInfo:       uainfo.UAInfo{Os: "android"},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
			},
			"all",
			2,
		},
		{
			[]messaging.Event{
				messaging.Event{
					AppID:        "C27A897D12F465AA",
					EventType:    api.CounterPrefix + "install",
					Channels:     []string{"a"},
					SenderID:     "uu1",
					CookieID:     "",
					UniqueID:     "",
					Count:        1,
					UAInfo:       uainfo.UAInfo{Os: "ios"},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
				messaging.Event{
					AppID:        "C27A897D12F465AA",
					EventType:    api.CounterPrefix + "install",
					Channels:     []string{"a"},
					SenderID:     "uu2",
					CookieID:     "",
					UniqueID:     "",
					Count:        1,
					UAInfo:       uainfo.UAInfo{Os: "ios"},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
				messaging.Event{
					AppID:        "C27A897D12F465AA",
					EventType:    api.CounterPrefix + "install",
					Channels:     []string{"a"},
					SenderID:     "uu2",
					CookieID:     "",
					UniqueID:     "uu2",
					Count:        1,
					UAInfo:       uainfo.UAInfo{Os: "ios"},
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
					UniqueID:     "uu2",
					Count:        1,
					UAInfo:       uainfo.UAInfo{Os: "android"},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
			},
			"ios",
			2,
		},
	}
	for i, v := range testcase {
		ds := NewGeneralDeviceStat(storage.NewInMemSimpleKV(), fmt.Sprintf("testCount%d", i))
		for _, device := range v.insertDevice {
			if err := ds.Insert(&device); err != nil {
				t.Fatal(err)
			}
		}
		count, err := ds.Count(v.os)
		if err != nil {
			t.Fatal(err)
		}
		if count != v.count {
			t.Errorf("#%d, wanted device number=%d, get=%d\n", i, v.count, count)
		}
	}
}
