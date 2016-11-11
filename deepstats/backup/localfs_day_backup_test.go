package backup

import (
	"os"
	"reflect"
	"testing"
)

func TestLocalfsBackupServiceInsertAndDecode(t *testing.T) {
	l, err := NewLocalFSBackupService("./testdata")
	defer os.RemoveAll("./testdata")
	if err != nil {
		t.Error(err)
		return
	}
	var testcases = testMessage()

	for _, v := range testcases {
		err := l.Insert(v)
		if err != nil {
			t.Error(err)
		}
	}

	events, err := l.RetriveAllEvents()
	if err != nil {
		t.Errorf("RetriveAllEvents failed! Err Msg=%v", err)
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
