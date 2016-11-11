package attribution

import (
	"testing"

	"log"
	"time"

	"reflect"

	"github.com/MISingularity/deepshare2/pkg/storage"
	"github.com/MISingularity/deepshare2/pkg/testutil"
	"gopkg.in/mgo.v2/bson"
)

func TestAttributionRetriever(t *testing.T) {
	testAppID := "testAppID"
	timestampStrs := []string{
		"2016-06-15_12:00:00",
		"2016-06-15_13:00:00",
		"2016-06-15_14:00:00",
		"2016-06-15_15:00:00",
		"2016-06-15_16:00:00",
	}
	unixTimes := []int64{}
	for _, t := range timestampStrs {
		timestamp, err := time.Parse("2006-01-02_15:04:05", t)
		if err != nil {
			panic(err)
		}
		unixTimes = append(unixTimes, timestamp.Unix())
	}
	tests := []struct {
		eventType string
		senderID  string
		channels  []string
		timestamp int64

		uniqueID      string
		wSenderID     string
		wAttr         *Attribution
		wAttrsInMongo []Attribution
	}{
		{"install", "s1", []string{"c1", "c2"}, unixTimes[0], "r", "s1",
			&Attribution{
				AppID:    testAppID,
				UniqueID: "r",
				SenderID: "s1",
				Channels: []string{"c1", "c2"},
				Tracking: trackingInstall,
				CreateAt: unixTimes[0],
			},
			[]Attribution{
				Attribution{
					AppID:    testAppID,
					UniqueID: "r",
					SenderID: "s1",
					Channels: []string{"c1", "c2"},
					Tracking: trackingInstall,
					CreateAt: unixTimes[0],
				},
			},
		},
		{"close", "", nil, unixTimes[1], "r", "",
			&Attribution{
				AppID:    testAppID,
				UniqueID: "r",
				SenderID: "s1",
				Channels: []string{"c1", "c2"},
				Tracking: trackingInstall,
				CreateAt: unixTimes[0],
				CloseAt:  unixTimes[1],
			},
			[]Attribution{
				Attribution{
					AppID:    testAppID,
					UniqueID: "r",
					SenderID: "s1",
					Channels: []string{"c1", "c2"},
					Tracking: trackingInstall,
					CreateAt: unixTimes[0],
					CloseAt:  unixTimes[1],
				},
			},
		},
		{"open", "", nil, unixTimes[2], "r", "",
			&Attribution{
				AppID:    testAppID,
				UniqueID: "r",
				SenderID: "s1",
				Channels: []string{"c1", "c2"},
				Tracking: trackingInstall,
				CreateAt: unixTimes[0],
				CloseAt:  unixTimes[1],
			},
			[]Attribution{
				Attribution{
					AppID:    testAppID,
					UniqueID: "r",
					SenderID: "s1",
					Channels: []string{"c1", "c2"},
					Tracking: trackingInstall,
					CreateAt: unixTimes[0],
					CloseAt:  unixTimes[1],
				},
			},
		},
		{"open", "s2", []string{"x", "y"}, unixTimes[3], "r", "s2",
			&Attribution{
				AppID:    testAppID,
				UniqueID: "r",
				SenderID: "s2",
				Channels: []string{"x", "y"},
				Tracking: trackingOpen,
				CreateAt: unixTimes[3],
			},
			[]Attribution{
				Attribution{
					AppID:    testAppID,
					UniqueID: "r",
					SenderID: "s1",
					Channels: []string{"c1", "c2"},
					Tracking: trackingInstall,
					CreateAt: unixTimes[0],
					CloseAt:  unixTimes[1],
				},
				Attribution{
					AppID:    testAppID,
					UniqueID: "r",
					SenderID: "s2",
					Channels: []string{"x", "y"},
					Tracking: trackingOpen,
					CreateAt: unixTimes[3],
				},
			},
		},
		{"close", "", nil, unixTimes[4], "r", "",
			&Attribution{
				AppID:    testAppID,
				UniqueID: "r",
				SenderID: "s2",
				Channels: []string{"x", "y"},
				Tracking: trackingOpen,
				CreateAt: unixTimes[3],
				CloseAt:  unixTimes[4],
			},
			[]Attribution{
				Attribution{
					AppID:    testAppID,
					UniqueID: "r",
					SenderID: "s1",
					Channels: []string{"c1", "c2"},
					Tracking: trackingInstall,
					CreateAt: unixTimes[0],
					CloseAt:  unixTimes[1],
				},
				Attribution{
					AppID:    testAppID,
					UniqueID: "r",
					SenderID: "s2",
					Channels: []string{"x", "y"},
					Tracking: trackingOpen,
					CreateAt: unixTimes[3],
					CloseAt:  unixTimes[4],
				},
			},
		},
	}

	session := testutil.MustNewLocalSession()
	mgoDB := session.DB("testAttr")
	mgoDB.DropDatabase()
	db := storage.NewRedisSimpleKV("127.0.0.1:6379", "", 200)
	as := NewAttributionRetriever(session, "testAttr", "attribution", db)

	for i, tt := range tests {
		log.Printf("#%d %v", i, tt)
		switch tt.eventType {
		case "install":
			as.OnInstall(testAppID, tt.uniqueID, tt.senderID, tt.channels, tt.timestamp)
		case "open":
			as.OnOpen(testAppID, tt.uniqueID, tt.senderID, tt.channels, tt.timestamp)
		case "close":
			as.OnClose(testAppID, tt.uniqueID, tt.timestamp)
		}

		k := formKey(testAppID, tt.uniqueID)
		v, err := db.HGet(k, "install")
		log.Println("install:", string(v), err)
		v, err = db.HGet(k, "open")
		log.Println("open:", string(v), err)
		attr := as.GetAttribution(testAppID, tt.uniqueID)
		log.Println("attr:", attr)
		log.Println("wattr:", tt.wAttr)
		if !reflect.DeepEqual(attr, tt.wAttr) {
			t.Errorf("#%d Expected Attribution: %v, got: %v", i, tt.wAttr, attr)
		}

		//mongoDB should save the attributions as wMgoData
		coll := mgoDB.C("attribution_" + testAppID)
		results := []Attribution{}
		if err := coll.Find(bson.M{}).Sort("create_at").All(&results); err != nil {
			t.Fatalf("#%d Failed to find from mongo, err: %v\n", i, err)
		}
		for j, _ := range results {
			results[j].ID = bson.ObjectId("")
		}
		if !reflect.DeepEqual(results, tt.wAttrsInMongo) {
			t.Errorf("#%d attribution saved in mongo does not match, attrs = %v, want = %v\n", i, results, tt.wAttrsInMongo)
		}
	}
}
