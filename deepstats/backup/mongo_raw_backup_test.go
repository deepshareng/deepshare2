package backup

import (
	"reflect"
	"testing"

	"github.com/MISingularity/deepshare2/deepstats"
)

func TestMongoBackupRawServiceInsertAndDecode(t *testing.T) {
	dbName := "test_mongo_raw_backup_insert_decode"
	collName := "backup"
	session := deepstats.MustCreateMongoSession(deepstats.MongoTestingAddr)
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()

	l := NewMongoRawBackupService(c)

	var testcases = testMessage()
	for _, v := range testcases {
		err := l.Insert(v)
		if err != nil {
			t.Error(err)
		}
	}

	events, err := l.RetriveAllEvents()
	if err != nil {
		t.Errorf("RetriveAllEvents %v failed! Err Msg=%v", c, err)
	}
	for _, v := range events {
		_, ok := testcases[v.AppID+"_"+v.EventType]
		if !ok {
			t.Errorf("Exist unexpected item, %v", v)
			continue
		}
		if !reflect.DeepEqual(v, testcases[v.AppID+"_"+v.EventType]) {
			t.Errorf("Mismatch, want=%v, get=%v", testcases[v.AppID+"_"+v.EventType], v)
		}
		delete(testcases, v.AppID+"_"+v.EventType)
	}
	if len(testcases) != 0 {
		t.Errorf("Testcase match failed, %v", testcases)
	}
}
