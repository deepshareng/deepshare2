package attribution

import (
	"bytes"
	"testing"

	"time"

	"log"

	"encoding/json"

	"reflect"

	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"github.com/MISingularity/deepshare2/pkg/testutil"
)

func TestAttributionService_Consume(t *testing.T) {
	timestampStrs := []string{
		"2016-06-17_12:00:00",
		"2016-06-17_13:00:00",
		"2016-06-17_14:00:00",
		"2016-06-17_15:00:00",
		"2016-06-17_16:00:00",
	}
	unixTimes := []int64{}
	for _, t := range timestampStrs {
		timestamp, err := time.Parse("2006-01-02_15:04:05", t)
		if err != nil {
			panic(err)
		}
		unixTimes = append(unixTimes, timestamp.Unix())
	}
	testAppID := "testAppID"
	tests := []struct {
		topic         string
		event         *messaging.Event
		wCounterEvent *messaging.Event
	}{
		{
			"inappdata",
			&messaging.Event{
				AppID:     testAppID,
				EventType: "/v2/inappdata/install",
				Channels:  []string{"1", "2"},
				SenderID:  "s1",
				UniqueID:  "u",
				Count:     1,
				TimeStamp: unixTimes[0],
			},
			nil,
		},
		{
			"counter_raw",
			&messaging.Event{
				AppID:     testAppID,
				EventType: "/v2/counters/c1",
				UniqueID:  "u",
				Count:     10,
				TimeStamp: unixTimes[1],
			},
			&messaging.Event{
				AppID:     testAppID,
				EventType: "/v2/counters/c1",
				Channels:  []string{"1", "2"},
				SenderID:  "s1",
				UniqueID:  "u",
				Count:     10,
				TimeStamp: unixTimes[1],
			},
		},
		{
			"inappdata",
			&messaging.Event{
				AppID:     testAppID,
				EventType: "/v2/inappdata/open",
				Channels:  []string{"a", "b"},
				SenderID:  "s2",
				UniqueID:  "u",
				Count:     1,
				TimeStamp: unixTimes[0],
			},
			nil,
		},
		{
			"counter_raw",
			&messaging.Event{
				AppID:     testAppID,
				EventType: "/v2/counters/c1",
				UniqueID:  "u",
				Count:     10,
				TimeStamp: unixTimes[1],
			},
			&messaging.Event{
				AppID:     testAppID,
				EventType: "/v2/counters/c1",
				Channels:  []string{"a", "b"},
				SenderID:  "s2",
				UniqueID:  "u",
				Count:     10,
				TimeStamp: unixTimes[1],
			},
		},
	}

	q := new(bytes.Buffer)
	producer := messaging.NewSimpleProducer(q)
	session := testutil.MustNewLocalSession()
	mgoDB := session.DB("testAS")
	mgoDB.DropDatabase()
	as := NewAttributionService(NewAttributionRetriever(session, "testAS", "attribution", storage.NewInMemSimpleKV()), producer)
	for i, tt := range tests {
		q.Reset()
		log.Printf("#%d topic: %s, event: %v\n", i, tt.topic, tt.event)
		if err := as.OnEvent(tt.event); err != nil {
			t.Fatalf("#%d OnEvent failed, err: %v\n", i, err)
		}
		if q.Len() == 0 {
			if tt.wCounterEvent != nil {
				t.Errorf("#%d OnEvent failed, counter event = %v, want = %v\n", i, q.String(), tt.wCounterEvent)
			}
		} else {
			se := new(messaging.SimpleProducerEvent)
			if err := json.Unmarshal(q.Bytes(), se); err != nil {
				t.Fatalf("#%d Failed to unmarshal json from messaging queue, err: %v, q.Bytes(): %v\n", i, err, q.Bytes())
			}
			e := &messaging.Event{}
			if err := json.Unmarshal(se.Msg, e); err != nil {
				t.Fatalf("#%d Failed to unmarshal json from se.Msg, err: %v, se.Msg: %v\n", i, err, string(se.Msg))
			}
			if !reflect.DeepEqual(e, tt.wCounterEvent) {
				t.Errorf("#%d Counter event: %v, want = %v\n", i, e, tt.wCounterEvent)
			}
		}
	}
}
