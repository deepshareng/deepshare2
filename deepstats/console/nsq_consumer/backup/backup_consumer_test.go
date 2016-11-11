package backup

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	pb "github.com/MISingularity/deepshare2/deepstats/backup/proto"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/nsqio/go-nsq"
)

type expectMsg string

type backupTestService struct {
	st map[expectMsg]int
}

func EventMarshal(e pb.Event) (string, error) {
	res, err := json.Marshal(e)
	return string(res), err
}

func (b *backupTestService) Insert(event pb.Event) error {
	s, err := EventMarshal(event)
	if err != nil {
		return err
	}
	b.st[expectMsg(s)]++
	return nil
}

func (b *backupTestService) RetriveAllEvents() ([]pb.Event, error) {
	return []pb.Event{}, nil
}
func TestAppEventConsume(t *testing.T) {
	testcases := []struct {
		input []messaging.Event
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
					UAInfo:       uainfo.UAInfo{},
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
					UAInfo:       uainfo.UAInfo{Brand: "3"},
					ReceiverInfo: api.MatchReceiverInfo{},
					KVs:          map[string]interface{}{"cookie_id": "", "inapp_data": "dfsdf", "ua": "UA:127.0.0.1_ios_7.0.3_-"},
					TimeStamp:    time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
			},
		},
	}

	for _, tt := range testcases {
		ats := backupTestService{make(map[expectMsg]int)}
		appconsumer := BackupConsumer{&ats}
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
		for _, em := range tt.input {
			e, err := convertPBEvent(em)
			if err != nil {
				t.Errorf("Marshal failed!Err Msg=%v", err)
				continue
			}
			bm, err := EventMarshal(e)
			if err != nil {
				t.Errorf("Marshal failed!Err Msg=%v", err)
				continue
			}
			m := expectMsg(bm)
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
