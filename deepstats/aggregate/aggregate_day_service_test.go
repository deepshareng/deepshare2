package aggregate

import (
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/MISingularity/deepshare2/deepstats"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"github.com/MISingularity/deepshare2/pkg/testutil"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// It tests that
// mongo AggregateService Insert(appID string) can add counter event and doesn't return error.
func TestMongoASInsert(t *testing.T) {
	dbName := "test_as_insert"
	collName := "counter"
	session := deepstats.MustCreateMongoSession(deepstats.MongoTestingAddr)
	defer session.Close()
	defer session.DB(dbName).DropDatabase()

	appID := "appid1"
	tests := []struct {
		event CounterEvent
	}{
		{CounterEvent{
			AppID:     appID,
			Channel:   "c1",
			Event:     "open",
			Timestamp: time.Now(),
			Count:     3,
			Os:        "android",
		}},
		{CounterEvent{
			AppID:     appID,
			Channel:   "c2",
			Event:     "open",
			Timestamp: time.Now(),
			Count:     2,
			Os:        "ios",
		}},
	}

	m := NewDayAggregateService(storage.NewInMemSimpleKV(), session, dbName, collName, "", "")

	intervalRefreshToMongoSec = 0

	for i, tt := range tests {
		err := m.Insert(tt.event.AppID, tt.event)
		if err != nil {
			t.Errorf("#%d: Insert failed: %v", i, err)
		}
	}
	m.Aggregate(appID)
	m.refreshToMongo()
	// check: all records were inserted
	c := GetColl(session, dbName, make(map[string]*mgo.Collection), collName, appID)
	num, err := c.Find(bson.M{}).Count()
	if err != nil {
		panic(err)
	}
	if num != len(tests) {
		t.Fatalf("number of inserted records=%d, want=%d", num, len(tests))
	}
}

// It tests that
// mongo AggregateService Aggregate() returns
// counts aggregated by "day" granularity and by event of a channel.
func TestMongoASAggregateByDay(t *testing.T) {
	dbName := "test_as_day"
	collName := "counter"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	defer session.DB(dbName).DropDatabase()

	m := NewDayAggregateService(storage.NewInMemSimpleKV(), session, dbName, collName, "", "")

	intervalRefreshToMongoSec = 0

	prepareEvents(m)

	tests := []struct {
		appid        string
		channel      string
		gran         time.Duration
		start        time.Time
		eventFilters []string
		os           string
		aggrs        aggrResults
	}{
		{ // #0
			"w33",
			"c11",
			time.Hour * 24,
			time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
			[]string{`version/.*`},
			"",
			aggrResults{
				{
					Event: "version/install",
					Counts: []*AggregateCount{
						{1}, {0}, {0}, {0}, {0}, {0}, {0},
					},
				},
				{
					Event: "version/2555",
					Counts: []*AggregateCount{
						{1}, {0}, {0}, {0}, {0}, {0}, {0},
					},
				},
				{
					Event: "version/234",
					Counts: []*AggregateCount{
						{1}, {0}, {0}, {0}, {0}, {0}, {0},
					},
				},
			},
		},
		{ // #1
			"w1",
			"c1",
			time.Hour * 24,
			time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
			[]string{"open"},
			"",
			aggrResults{
				{
					Event: "open",
					Counts: []*AggregateCount{
						{3}, {1}, {0}, {0}, {0}, {0}, {1},
					},
				},
			},
		},
		{ // #2
			"w1",
			"c1",
			time.Hour * 24,
			time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
			[]string{"install"},
			"",
			aggrResults{
				{
					Event: "install",
					Counts: []*AggregateCount{
						{1}, {0}, {0}, {0}, {0}, {0}, {0},
					},
				},
			},
		},
		{ // #3
			"w2",
			"c1",
			time.Hour * 24,
			time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
			[]string{"open"},
			"",
			aggrResults{
				{
					Event: "open",
					Counts: []*AggregateCount{
						{1}, {0}, {0}, {0}, {0}, {0}, {0},
					},
				},
			},
		},
		{ // #4
			"w3",
			"c1",
			time.Hour * 24 * 7,
			time.Date(2014, time.December, 17, 0, 0, 0, 0, time.Local),
			[]string{"open"},
			"",
			aggrResults{
				{
					Event: "open",
					Counts: []*AggregateCount{
						{1}, {1}, {3}, {0}, {0}, {0}, {0},
					},
				},
			},
		},
		{ // #5
			"w3",
			"c1",
			time.Hour * 24 * 30,
			time.Date(2014, time.December, 2, 0, 0, 0, 0, time.Local),
			[]string{},
			"",
			aggrResults{
				{
					Event: "open",
					Counts: []*AggregateCount{
						{3}, {2}, {0}, {0}, {0}, {0}, {0},
					},
				},
			},
		},
		{ // #6
			"w1",
			"non-exist",
			time.Hour * 24,
			time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local),
			nil,
			"",
			aggrResults{},
		},
		{ // #7 test os
			"w33",
			"c11",
			time.Hour * 24,
			time.Date(2020, time.January, 1, 0, 0, 0, 0, time.Local),
			nil,
			"",
			aggrResults{
				{
					Event: "version/install",
					Counts: []*AggregateCount{
						{13}, {0}, {0}, {0}, {0}, {0}, {0},
					},
				},
			},
		},
		{ // #8 test os
			"w33",
			"c11",
			time.Hour * 24,
			time.Date(2020, time.January, 1, 0, 0, 0, 0, time.Local),
			nil,
			"ios",
			aggrResults{
				{
					Event: "version/install",
					Counts: []*AggregateCount{
						{2}, {0}, {0}, {0}, {0}, {0}, {0},
					},
				},
			},
		},
	}

	m.refreshToMongo()

	for i, tt := range tests {
		aggrraw, err := m.QueryDuration(tt.appid, tt.channel, tt.eventFilters, tt.start, tt.gran, 7, tt.os)
		aggrs := aggrResults(aggrraw)
		sort.Sort(aggrs)
		if err != nil {
			t.Fatalf("#%d: Aggregate Failed: %v", i, err)
		}
		if !reflect.DeepEqual(aggrs, tt.aggrs) {
			t.Errorf("#%d:\nAggregateResult(s)=%s\nwant=%s\n", i,
				AggregateResultsToString(aggrs), AggregateResultsToString(tt.aggrs))
		}
	}
}

// Prepare phase:
// - set up events for only one channel "c1"
// - set up a couple of "open" events:
//   - two can be aggregated to the same day
//   - one has a more recent day
//   - one has a more recent month
//   - one has a more recent year
// - set up one other "install" event
func prepareEvents(m AggregateService) {
	// test os at 2020
	m.Insert("w33", CounterEvent{
		AppID:     "w33",
		Channel:   "c11",
		Event:     "version/install",
		Timestamp: time.Date(2020, time.January, 1, 10, 0, 10, 0, time.Local),
		Count:     1,
		Os:        "ios",
	})
	m.Insert("w33", CounterEvent{
		AppID:     "w33",
		Channel:   "c11",
		Event:     "version/install",
		Timestamp: time.Date(2020, time.January, 1, 10, 0, 10, 0, time.Local),
		Count:     1,
		Os:        "ios",
	})
	m.Insert("w33", CounterEvent{
		AppID:     "w33",
		Channel:   "c11",
		Event:     "version/install",
		Timestamp: time.Date(2020, time.January, 1, 10, 0, 10, 0, time.Local),
		Count:     1,
		Os:        "other",
	})
	m.Insert("w33", CounterEvent{
		AppID:     "w33",
		Channel:   "c11",
		Event:     "version/install",
		Timestamp: time.Date(2020, time.January, 1, 10, 0, 10, 0, time.Local),
		Count:     10,
		Os:        "android",
	})
	// w33
	m.Insert("w33", CounterEvent{
		AppID:     "w33",
		Channel:   "c11",
		Event:     "version/install",
		Timestamp: time.Date(2015, time.January, 1, 10, 0, 10, 0, time.Local),
		Count:     1,
	})
	m.Insert("w33", CounterEvent{
		AppID:     "w33",
		Channel:   "c11",
		Event:     "version/234",
		Timestamp: time.Date(2015, time.January, 1, 10, 0, 10, 0, time.Local),
		Count:     1,
	})
	m.Insert("w33", CounterEvent{
		AppID:     "w33",
		Channel:   "c11",
		Event:     "version/2555",
		Timestamp: time.Date(2015, time.January, 1, 10, 0, 10, 0, time.Local),
		Count:     1,
	})
	m.Insert("w33", CounterEvent{
		AppID:     "w33",
		Channel:   "c11",
		Event:     "open",
		Timestamp: time.Date(2015, time.January, 1, 10, 0, 10, 0, time.Local),
		Count:     1,
	})

	//w3
	m.Insert("w3", CounterEvent{
		AppID:     "w3",
		Channel:   "c1",
		Event:     "open",
		Timestamp: time.Date(2014, time.December, 17, 10, 0, 10, 0, time.Local),
		Count:     1,
	})
	m.Insert("w3", CounterEvent{
		AppID:     "w3",
		Channel:   "c1",
		Event:     "open",
		Timestamp: time.Date(2014, time.December, 30, 10, 0, 10, 0, time.Local),
		Count:     1,
	})
	m.Insert("w3", CounterEvent{
		AppID:     "w3",
		Channel:   "c1",
		Event:     "open",
		Timestamp: time.Date(2014, time.December, 31, 10, 0, 10, 0, time.Local),
		Count:     1,
	})
	m.Insert("w3", CounterEvent{
		AppID:     "w3",
		Channel:   "c1",
		Event:     "open",
		Timestamp: time.Date(2015, time.January, 1, 10, 0, 10, 0, time.Local),
		Count:     1,
	})
	m.Insert("w3", CounterEvent{
		AppID:     "w3",
		Channel:   "c1",
		Event:     "open",
		Timestamp: time.Date(2015, time.January, 2, 10, 0, 10, 0, time.Local),
		Count:     1,
	})

	//w2
	m.Insert("w2", CounterEvent{
		AppID:     "w2",
		Channel:   "c1",
		Event:     "open",
		Timestamp: time.Date(2015, time.January, 1, 10, 0, 10, 0, time.Local),
		Count:     1,
	})

	//w1
	m.Insert("w1", CounterEvent{
		AppID:     "w1",
		Channel:   "c1",
		Event:     "open",
		Timestamp: time.Date(2015, time.January, 1, 18, 0, 10, 0, time.Local),
		Count:     1,
	})
	m.Insert("w1", CounterEvent{
		AppID:     "w1",
		Channel:   "c1",
		Event:     "open",
		Timestamp: time.Date(2015, time.January, 1, 20, 0, 10, 0, time.Local),
		Count:     1,
	})
	m.Insert("w1", CounterEvent{
		AppID:     "w1",
		Channel:   "c1",
		Event:     "open",
		Timestamp: time.Date(2015, time.January, 2, 20, 0, 10, 0, time.Local),
		Count:     1,
	})
	m.Insert("w1", CounterEvent{
		AppID:     "w1",
		Channel:   "c1",
		Event:     "open",
		Timestamp: time.Date(2015, time.January, 7, 20, 0, 10, 0, time.Local),
		Count:     1,
	})
	m.Insert("w1", CounterEvent{
		AppID:     "w1",
		Channel:   "c1",
		Event:     "open",
		Timestamp: time.Date(2015, time.January, 1, 20, 0, 10, 0, time.Local),
		Count:     1,
	})
	m.Insert("w1", CounterEvent{
		AppID:     "w1",
		Channel:   "c1",
		Event:     "install",
		Timestamp: time.Date(2015, time.January, 1, 20, 0, 10, 0, time.Local),
		Count:     1,
	})

	m.Aggregate("w1")  //aggregate results insert to <collnamePrefix>w1
	m.Aggregate("w2")  //aggregate results insert to <collnamePrefix>w2
	m.Aggregate("w3")  //aggregate results insert to <collnamePrefix>w3
	m.Aggregate("w33") //aggregate results insert to <collnamePrefix>w33
}
