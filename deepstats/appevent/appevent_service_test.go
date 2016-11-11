package appevent

import (
	"reflect"
	"testing"

	"github.com/MISingularity/deepshare2/deepstats"
)

// It tests that
// mongo AppEvent Service Getevents() returns
// returns the events of specified app
func TestMongoASGetEvents(t *testing.T) {
	dbName := "test_as_getevents"
	collName := "app_events"

	session := deepstats.MustCreateMongoSession(deepstats.MongoTestingAddr)
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()

	m := NewMongoAppEventService(c)
	prepareEvents(m)
	tests := []AppEvents{
		{"a1",
			[]string{
				"match/e1",
				"dsaction/e2",
				"counter/e1",
			},
		},
		{"a2",
			[]string{
				"counter/e1",
			},
		},
		{"a3",
			[]string{
				"counter/e1",
				"counter/e2",
			},
		},
	}

	for i, tt := range tests {
		res, err := m.GetEvents(tt.AppID)
		if err != nil {
			t.Fatalf("#%d: Get Event Failed: %v", i, err)
		}
		if !reflect.DeepEqual(res, tt) {
			t.Errorf("#%d:\nGetEvents(appid)=%s\nwant=%s\n", i,
				res, tt)
		}
	}
}

func prepareEvents(m AppEventService) {
	m.InsertEvent("a1", "match/e1")
	m.InsertEvent("a1", "match/e1")
	m.InsertEvent("a1", "dsaction/e2")
	m.InsertEvent("a1", "counter/e1")
	m.InsertEvent("a1", "counter/e1")
	m.InsertEvent("a1", "counter/e1")
	m.InsertEvent("a2", "counter/e1")
	m.InsertEvent("a2", "counter/e1")
	m.InsertEvent("a3", "counter/e1")
	m.InsertEvent("a3", "counter/e2")

}
