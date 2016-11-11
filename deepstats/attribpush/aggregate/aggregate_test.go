package aggregate

import (
	"reflect"
	"testing"
	"time"

	"github.com/MISingularity/deepshare2/deepstats/attribpush"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/MISingularity/deepshare2/pkg/testutil"
	"gopkg.in/mgo.v2/bson"
)

func TestAggregate(t *testing.T) {
	testSenderID := "test_sender"
	testAppID := "test_app"
	time1 := time.Date(2015, time.November, 1, 0, 0, 0, 0, time.Local).Unix()
	time2 := time.Date(2015, time.November, 2, 0, 0, 0, 0, time.Local).Unix()
	time3 := time.Date(2015, time.November, 3, 0, 0, 0, 0, time.Local).Unix()
	tests := []struct {
		eventType string
		count     int
		timestamp int64

		wAttribution AttributionInfo
	}{
		{
			eventType: "match/install",
			count:     1,
			timestamp: time1,

			wAttribution: AttributionInfo{
				testAppID,
				testSenderID,
				map[string][]ValueTimePair{
					"ds/install": []ValueTimePair{
						{1, time1},
					},
				},
			},
		},
		{
			eventType: "match/open",
			count:     1,
			timestamp: time2,

			wAttribution: AttributionInfo{
				testAppID,
				testSenderID,
				map[string][]ValueTimePair{
					"ds/install": []ValueTimePair{
						{1, time1},
					},
					"ds/open": []ValueTimePair{
						{1, time2},
					},
				},
			},
		},
		{
			eventType: "match/open",
			count:     1,
			timestamp: time3,

			wAttribution: AttributionInfo{
				testAppID,
				testSenderID,
				map[string][]ValueTimePair{
					"ds/install": []ValueTimePair{
						{1, time1},
					},
					"ds/open": []ValueTimePair{
						{1, time2},
						{1, time3},
					},
				},
			},
		},
	}
	sa := newTestSimpleAggregator()

	for i, tt := range tests {
		event := messaging.Event{
			AppID:     testAppID,
			EventType: tt.eventType,
			SenderID:  testSenderID,
			Count:     tt.count,
			TimeStamp: tt.timestamp,
		}
		if err := sa.Aggregate(event); err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Duration(time.Second))
		var result AttributionInfo
		if err := sa.mgocoll.Find(bson.M{}).One(&result); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(result, tt.wAttribution) {
			t.Errorf("#%d got: %v, want: %v\n", i, result, tt.wAttribution)
		}
	}

}

func newTestSimpleAggregator() *simpleAggregator {
	attrColl := testutil.MustNewMongoColl("testAttribution", "attrcoll")
	rsColl := testutil.MustNewMongoColl("testAttribution", "rscoll")
	return &simpleAggregator{
		ar:      attribpush.NewAttributionParser(rsColl),
		mgocoll: attrColl,
	}
}
