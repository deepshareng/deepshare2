package retention

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/nsqio/go-nsq"
)

type expectMsg string

type retentionTestService struct {
	coll map[expectMsg]int
}

func (r *retentionTestService) FindMatchUserForRetention(dp *messaging.Event) error {
	res, err := json.Marshal(dp)
	if err != nil {
		return err
	}
	r.coll["open#"+expectMsg(res)]++
	return nil
}

func (r *retentionTestService) InsertInstallEventForRetention(dp *messaging.Event) error {
	res, err := json.Marshal(dp)
	if err != nil {
		return err
	}
	r.coll["input#"+expectMsg(res)]++
	return nil
}

func TestThreeDayRetention(t *testing.T) {
	retentionService := &retentionTestService{make(map[expectMsg]int)}
	retentionConsumer := NewRetentionConsumer(retentionService)
	testcases := []struct {
		input         []messaging.Event
		expectMessage []expectMsg
	}{

		{
			[]messaging.Event{
				// Don't meet the requirement
				messaging.Event{
					AppID:        "app3",
					EventType:    api.MatchPrefix + "install",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "4",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
				messaging.Event{
					AppID:        "app3",
					EventType:    "install",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "4",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},

				// events satisfy rules
				// same event
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
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC).Unix(),
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
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC).Unix(),
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
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},

				// distinct event

				messaging.Event{
					AppID:        "app3",
					EventType:    api.GetInAppDataPrefix + "install",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "5",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2014, time.December, 29, 0, 0, 0, 0, time.UTC).Unix(),
				},
			},
			[]expectMsg{
				`input#{"appid":"app3","event":"/v2/inappdata/install","channels":[""],"unique_id":"4","count":1,"ua_info":{"Ua":"","Ip":"","Os":"","OsVersion":"","Brand":"","Browser":"","IsWechat":false,"IsWeibo":false,"IsQQ":false,"IsTwitter":false,"IsFacebook":false,"IsFirefox":false,"IsQQBrowser":false,"ChromeMajor":0,"CannotDeeplink":false,"CannotGoMarket":false,"CannotGetWindowsEvent":false,"ForceUseScheme":false},"receiver_info":{"unique_id":"","app_version_name":"","app_version_code":0,"app_version_build":"","sdk_info":"","is_wifi_connected":false,"carrier_name":"","blueTooth_enable":false,"model":"","brand":"","is_emulator":false,"os":"","os_version":"","has_nfc":false,"has_telephone":false,"bluetooth_version":"","screen_dpi":0,"screen_width":0,"screen_height":0,"uri_scheme":"","hardware_id":""},"kvs":{"cookie_id":"","inapp_data":"dfsdf","ua":"UA:127.0.0.1_ios_7.0.3_-"},"timestamp":1420070400}`,
				`input#{"appid":"app3","event":"/v2/inappdata/install","channels":[""],"unique_id":"4","count":1,"ua_info":{"Ua":"","Ip":"","Os":"","OsVersion":"","Brand":"","Browser":"","IsWechat":false,"IsWeibo":false,"IsQQ":false,"IsTwitter":false,"IsFacebook":false,"IsFirefox":false,"IsQQBrowser":false,"ChromeMajor":0,"CannotDeeplink":false,"CannotGoMarket":false,"CannotGetWindowsEvent":false,"ForceUseScheme":false},"receiver_info":{"unique_id":"","app_version_name":"","app_version_code":0,"app_version_build":"","sdk_info":"","is_wifi_connected":false,"carrier_name":"","blueTooth_enable":false,"model":"","brand":"","is_emulator":false,"os":"","os_version":"","has_nfc":false,"has_telephone":false,"bluetooth_version":"","screen_dpi":0,"screen_width":0,"screen_height":0,"uri_scheme":"","hardware_id":""},"kvs":{"cookie_id":"","inapp_data":"dfsdf","ua":"UA:127.0.0.1_ios_7.0.3_-"},"timestamp":1420070400}`,
				`input#{"appid":"app3","event":"/v2/inappdata/install","channels":[""],"unique_id":"4","count":1,"ua_info":{"Ua":"","Ip":"","Os":"","OsVersion":"","Brand":"","Browser":"","IsWechat":false,"IsWeibo":false,"IsQQ":false,"IsTwitter":false,"IsFacebook":false,"IsFirefox":false,"IsQQBrowser":false,"ChromeMajor":0,"CannotDeeplink":false,"CannotGoMarket":false,"CannotGetWindowsEvent":false,"ForceUseScheme":false},"receiver_info":{"unique_id":"","app_version_name":"","app_version_code":0,"app_version_build":"","sdk_info":"","is_wifi_connected":false,"carrier_name":"","blueTooth_enable":false,"model":"","brand":"","is_emulator":false,"os":"","os_version":"","has_nfc":false,"has_telephone":false,"bluetooth_version":"","screen_dpi":0,"screen_width":0,"screen_height":0,"uri_scheme":"","hardware_id":""},"kvs":{"cookie_id":"","inapp_data":"dfsdf","ua":"UA:127.0.0.1_ios_7.0.3_-"},"timestamp":1420070400}`,
				`input#{"appid":"app3","event":"/v2/inappdata/install","channels":[""],"unique_id":"5","count":1,"ua_info":{"Ua":"","Ip":"","Os":"","OsVersion":"","Brand":"","Browser":"","IsWechat":false,"IsWeibo":false,"IsQQ":false,"IsTwitter":false,"IsFacebook":false,"IsFirefox":false,"IsQQBrowser":false,"ChromeMajor":0,"CannotDeeplink":false,"CannotGoMarket":false,"CannotGetWindowsEvent":false,"ForceUseScheme":false},"receiver_info":{"unique_id":"","app_version_name":"","app_version_code":0,"app_version_build":"","sdk_info":"","is_wifi_connected":false,"carrier_name":"","blueTooth_enable":false,"model":"","brand":"","is_emulator":false,"os":"","os_version":"","has_nfc":false,"has_telephone":false,"bluetooth_version":"","screen_dpi":0,"screen_width":0,"screen_height":0,"uri_scheme":"","hardware_id":""},"kvs":{"cookie_id":"","inapp_data":"dfsdf","ua":"UA:127.0.0.1_ios_7.0.3_-"},"timestamp":1419811200}`,
			},
		},
		{
			[]messaging.Event{
				// install
				messaging.Event{
					AppID:        "app3",
					EventType:    api.MatchPrefix + "open",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "4",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
				messaging.Event{
					AppID:        "app3",
					EventType:    "open",
					Channels:     []string{""},
					SenderID:     "",
					CookieID:     "",
					UniqueID:     "4",
					Count:        1,
					UAInfo:       uainfo.UAInfo{},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},

				// open invent
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
					TimeStamp:    time.Date(2015, time.January, 4, 0, 0, 0, 0, time.UTC).Unix(),
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
					TimeStamp:    time.Date(2015, time.January, 4, 0, 0, 0, 0, time.UTC).Unix(),
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
					TimeStamp:    time.Date(2015, time.January, 4, 0, 0, 0, 0, time.UTC).Unix(),
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
					TimeStamp:    time.Date(2015, time.January, 3, 0, 0, 0, 0, time.UTC).Unix(),
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
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
			},
			[]expectMsg{
				`open#{"appid":"app3","event":"/v2/inappdata/open","channels":[""],"unique_id":"87","count":1,"ua_info":{"Ua":"","Ip":"","Os":"","OsVersion":"","Brand":"","Browser":"","IsWechat":false,"IsWeibo":false,"IsQQ":false,"IsTwitter":false,"IsFacebook":false,"IsFirefox":false,"IsQQBrowser":false,"ChromeMajor":0,"CannotDeeplink":false,"CannotGoMarket":false,"CannotGetWindowsEvent":false,"ForceUseScheme":false},"receiver_info":{"unique_id":"","app_version_name":"","app_version_code":0,"app_version_build":"","sdk_info":"","is_wifi_connected":false,"carrier_name":"","blueTooth_enable":false,"model":"","brand":"","is_emulator":false,"os":"","os_version":"","has_nfc":false,"has_telephone":false,"bluetooth_version":"","screen_dpi":0,"screen_width":0,"screen_height":0,"uri_scheme":"","hardware_id":""},"kvs":{"cookie_id":"","inapp_data":"dfsdf","ua":"UA:127.0.0.1_ios_7.0.3_-"},"timestamp":1420329600}`,
				`open#{"appid":"app3","event":"/v2/inappdata/open","channels":[""],"unique_id":"87","count":1,"ua_info":{"Ua":"","Ip":"","Os":"","OsVersion":"","Brand":"","Browser":"","IsWechat":false,"IsWeibo":false,"IsQQ":false,"IsTwitter":false,"IsFacebook":false,"IsFirefox":false,"IsQQBrowser":false,"ChromeMajor":0,"CannotDeeplink":false,"CannotGoMarket":false,"CannotGetWindowsEvent":false,"ForceUseScheme":false},"receiver_info":{"unique_id":"","app_version_name":"","app_version_code":0,"app_version_build":"","sdk_info":"","is_wifi_connected":false,"carrier_name":"","blueTooth_enable":false,"model":"","brand":"","is_emulator":false,"os":"","os_version":"","has_nfc":false,"has_telephone":false,"bluetooth_version":"","screen_dpi":0,"screen_width":0,"screen_height":0,"uri_scheme":"","hardware_id":""},"kvs":{"cookie_id":"","inapp_data":"dfsdf","ua":"UA:127.0.0.1_ios_7.0.3_-"},"timestamp":1420329600}`,
				`open#{"appid":"app3","event":"/v2/inappdata/open","channels":[""],"unique_id":"4","count":1,"ua_info":{"Ua":"","Ip":"","Os":"","OsVersion":"","Brand":"","Browser":"","IsWechat":false,"IsWeibo":false,"IsQQ":false,"IsTwitter":false,"IsFacebook":false,"IsFirefox":false,"IsQQBrowser":false,"ChromeMajor":0,"CannotDeeplink":false,"CannotGoMarket":false,"CannotGetWindowsEvent":false,"ForceUseScheme":false},"receiver_info":{"unique_id":"","app_version_name":"","app_version_code":0,"app_version_build":"","sdk_info":"","is_wifi_connected":false,"carrier_name":"","blueTooth_enable":false,"model":"","brand":"","is_emulator":false,"os":"","os_version":"","has_nfc":false,"has_telephone":false,"bluetooth_version":"","screen_dpi":0,"screen_width":0,"screen_height":0,"uri_scheme":"","hardware_id":""},"kvs":{"cookie_id":"","inapp_data":"dfsdf","ua":"UA:127.0.0.1_ios_7.0.3_-"},"timestamp":1420329600}`,
				`open#{"appid":"app3","event":"/v2/inappdata/open","channels":[""],"unique_id":"4","count":1,"ua_info":{"Ua":"","Ip":"","Os":"","OsVersion":"","Brand":"","Browser":"","IsWechat":false,"IsWeibo":false,"IsQQ":false,"IsTwitter":false,"IsFacebook":false,"IsFirefox":false,"IsQQBrowser":false,"ChromeMajor":0,"CannotDeeplink":false,"CannotGoMarket":false,"CannotGetWindowsEvent":false,"ForceUseScheme":false},"receiver_info":{"unique_id":"","app_version_name":"","app_version_code":0,"app_version_build":"","sdk_info":"","is_wifi_connected":false,"carrier_name":"","blueTooth_enable":false,"model":"","brand":"","is_emulator":false,"os":"","os_version":"","has_nfc":false,"has_telephone":false,"bluetooth_version":"","screen_dpi":0,"screen_width":0,"screen_height":0,"uri_scheme":"","hardware_id":""},"kvs":{"cookie_id":"","inapp_data":"dfsdf","ua":"UA:127.0.0.1_ios_7.0.3_-"},"timestamp":1420243200}`,
				`open#{"appid":"app3","event":"/v2/inappdata/open","channels":[""],"unique_id":"5","count":1,"ua_info":{"Ua":"","Ip":"","Os":"","OsVersion":"","Brand":"","Browser":"","IsWechat":false,"IsWeibo":false,"IsQQ":false,"IsTwitter":false,"IsFacebook":false,"IsFirefox":false,"IsQQBrowser":false,"ChromeMajor":0,"CannotDeeplink":false,"CannotGoMarket":false,"CannotGetWindowsEvent":false,"ForceUseScheme":false},"receiver_info":{"unique_id":"","app_version_name":"","app_version_code":0,"app_version_build":"","sdk_info":"","is_wifi_connected":false,"carrier_name":"","blueTooth_enable":false,"model":"","brand":"","is_emulator":false,"os":"","os_version":"","has_nfc":false,"has_telephone":false,"bluetooth_version":"","screen_dpi":0,"screen_width":0,"screen_height":0,"uri_scheme":"","hardware_id":""},"kvs":{"cookie_id":"","inapp_data":"dfsdf","ua":"UA:127.0.0.1_ios_7.0.3_-"},"timestamp":1420070400}`,
			},
		},
	}
	for _, tt := range testcases {
		retentionService.coll = make(map[expectMsg]int)
		for _, m := range tt.input {
			d, err := json.Marshal(m)
			if err != nil {
				t.Errorf("%v", err)
				continue
			}
			msg := &nsq.Message{
				Body: []byte(d),
			}
			err = retentionConsumer.Consume(msg)

			if err != nil {
				t.Errorf("Consumer Msg failed! Msg Detail=%v, Err Msg=%v", m, err)
				continue
			}
		}
		for _, v := range tt.expectMessage {
			_, ok := retentionService.coll[v]
			if !ok {
				t.Errorf("Expected retention message is not exist! Expected Msg = %s", v)
				continue
			}
			retentionService.coll[v]--
			if retentionService.coll[v] == 0 {
				delete(retentionService.coll, v)
			}
		}
		if len(retentionService.coll) != 0 {
			t.Errorf("Remain message number mismatsh! Expected = %d, get = %d.", 0, len(retentionService.coll))
		}
		for k, v := range retentionService.coll {
			t.Log(k, v)
			fmt.Println(k, v)
		}
	}
}
