package backup

import (
	"reflect"
	"testing"
	"time"

	"fmt"
	"regexp"

	"strconv"
	"strings"

	"github.com/MISingularity/deepshare2/deepstats"
	pb "github.com/MISingularity/deepshare2/deepstats/backup/proto"
	"gopkg.in/mgo.v2/bson"
)

func TestRawDayBackup(t *testing.T) {
	dbNamePrefix := "test_raw_day_backup"
	session := deepstats.MustCreateMongoSession(deepstats.MongoTestingAddr)
	defer session.Close()
	defer func() {
		dbs, err := session.DatabaseNames()
		if err != nil {
			return
		}
		reg := regexp.MustCompile(fmt.Sprintf("%s_\\d{6}", dbNamePrefix))
		for _, dbName := range dbs {
			if reg.Match([]byte(dbName)) {
				session.DB(dbName).DropDatabase()
			}
		}
	}()
	l := NewMongoRawDayBackupService(session, dbNamePrefix)

	var testcases = testMessage()
	for _, v := range testcases {
		go func(v pb.Event) {
			err := l.Insert(v)
			if err != nil {
				t.Error(v)
			}
		}(v)
	}
	time.Sleep(time.Duration(200) * time.Millisecond)
	// read All events and compare
	events, err := l.RetriveAllEvents()
	if err != nil {
		t.Fatalf("RetriveAllEvents failed! Err Msg=%v", err)
	}

	for _, v := range events {
		_, ok := testcases[v.AppID+"_"+v.EventType]
		if !ok {
			t.Errorf("Exist unexpected item, %v, %s", v, v.AppID+"_"+v.EventType)
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

	// every events should be in the right db & collection
	dbs, err := session.DatabaseNames()
	if err != nil {
		panic(err)
	}
	regDB := regexp.MustCompile(fmt.Sprintf("%s_\\d{6}", dbNamePrefix))
	regColl := regexp.MustCompile(`coll\d{2}`)
	for _, dbName := range dbs {
		if !regDB.Match([]byte(dbName)) {
			continue
		}
		db := session.DB(dbName)
		dbYear, _ := strconv.Atoi(strings.TrimPrefix(dbName, dbNamePrefix+"_")[:4])
		dbMonth, _ := strconv.Atoi(strings.TrimPrefix(dbName, dbNamePrefix+"_")[4:])
		colls, err := db.CollectionNames()
		if err != nil {
			panic(err)
		}
		for _, collName := range colls {
			if !regColl.Match([]byte(collName)) {
				continue
			}
			collDay, _ := strconv.Atoi(strings.TrimPrefix(collName, "coll"))
			c := db.C(collName)
			e := pb.Event{}
			iter := c.Find(bson.M{}).Iter()
			for iter.Next(&e) {
				eventTime := time.Unix(e.Timestamp, 0)
				if eventTime.Year() != dbYear || int(eventTime.Month()) != dbMonth || eventTime.Day() != collDay {
					t.Errorf("event is stored to the wrong db/collection, event:%v, dbName:%s, collName:%s", e, dbName, collName)
				}
			}
		}
	}
}

func testMessage() map[string]pb.Event {
	tt1 := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	tt2 := time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC)
	tt3 := time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC)
	tt4 := time.Date(2015, 1, 3, 0, 0, 0, 0, time.UTC)
	tt5 := time.Date(2015, 3, 1, 0, 0, 0, 0, time.UTC)
	tt6 := time.Date(2016, 1, 6, 0, 0, 0, 0, time.UTC)
	tt7 := time.Date(2015, 1, 3, 0, 0, 0, 0, time.UTC)
	tt8 := time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC)
	tt9 := time.Date(2016, 1, 5, 0, 0, 0, 0, time.UTC)
	UA1 := pb.UAInfo{Ua: "u1", Ip: "i1", Os: "o1", OsVersion: "v1", Brand: "b1", Browser: "b1", IsWechat: true, IsWeibo: false, IsQQ: true, ChromeMajor: 23}
	UA2 := pb.UAInfo{Ua: "u2", Ip: "i2", Os: "o2", OsVersion: "v2", Brand: "b2", Browser: "b2", IsWechat: false, IsWeibo: true, IsQQ: false, ChromeMajor: -1}
	return map[string]pb.Event{
		"a1_e1":  pb.Event{AppID: "a1", Channels: []string{"c1", "c2"}, Timestamp: tt1.Unix(), EventType: "e1", SenderID: "s1", CookieID: "c1", UniqueID: "u1", UAInfo: &UA1, KVs: []byte{1}, Count: 1},
		"a1_e2":  pb.Event{AppID: "a1", Channels: []string{"c2"}, Timestamp: tt2.Unix(), EventType: "e2", SenderID: "s2", CookieID: "c2", UniqueID: "u2", UAInfo: &UA2, KVs: []byte{2}, Count: 2},
		"a1_e3":  pb.Event{AppID: "a1", Channels: []string{"c3", "c45", "c41"}, Timestamp: tt2.Unix(), SenderID: "s3", CookieID: "c3", UniqueID: "u3", UAInfo: &UA1, KVs: []byte{3}, EventType: "e3", Count: 3},
		"a1_e4":  pb.Event{AppID: "a1", Channels: []string{"c4"}, Timestamp: tt3.Unix(), SenderID: "s4", CookieID: "c4", UniqueID: "u4", UAInfo: &UA1, KVs: []byte{4}, EventType: "e4", Count: 4},
		"a1_e5":  pb.Event{AppID: "a1", Channels: []string{"c5"}, Timestamp: tt4.Unix(), SenderID: "s5", CookieID: "c5", UniqueID: "u5", UAInfo: &UA1, KVs: []byte{5}, EventType: "e5", Count: 5},
		"a1_e6":  pb.Event{AppID: "a1", Channels: []string{"c6"}, Timestamp: tt5.Unix(), SenderID: "s6", CookieID: "c6", UniqueID: "u5", UAInfo: &UA1, KVs: []byte{6}, EventType: "e6", Count: 6},
		"a1_e7":  pb.Event{AppID: "a1", Channels: []string{"c7"}, Timestamp: tt6.Unix(), SenderID: "s7", CookieID: "c7", UniqueID: "u5", UAInfo: &UA1, KVs: []byte{7}, EventType: "e7", Count: 7},
		"a1_e8":  pb.Event{AppID: "a1", Channels: []string{"c8"}, Timestamp: tt7.Unix(), SenderID: "s8", CookieID: "c8", UniqueID: "u5", UAInfo: &UA1, KVs: []byte{8}, EventType: "e8", Count: 8},
		"a1_e9":  pb.Event{AppID: "a1", Channels: []string{"c9"}, Timestamp: tt8.Unix(), SenderID: "s9", CookieID: "c9", UniqueID: "u5", UAInfo: &UA1, KVs: []byte{9}, EventType: "e9", Count: 9},
		"a1_e10": pb.Event{AppID: "a1", Channels: []string{"c10"}, Timestamp: tt9.Unix(), SenderID: "s10", CookieID: "c10", UniqueID: "u5", UAInfo: &UA1, KVs: []byte{10}, EventType: "e10", Count: 10},
	}
}
