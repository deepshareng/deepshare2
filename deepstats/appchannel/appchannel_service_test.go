package appchannel

import (
	"reflect"
	"testing"

	"github.com/MISingularity/deepshare2/deepstats"
)

// It tests that
// mongo AppChannel Service Getchannels() returns
// returns the channels of specified app
func TestMongoASGetChannels(t *testing.T) {
	dbName := "test_as_getchannels"
	collName := "app_channels"

	session := deepstats.MustCreateMongoSession(deepstats.MongoTestingAddr)
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()

	m := NewMongoAppChannelService(c)
	prepareChannels(m)
	tests := []AppChannels{
		{"a1",
			[]string{
				"c1",
			},
		},
		{"a2",
			[]string{
				"c2", "c3",
			},
		},
		{"a3",
			[]string{
				"c4",
			},
		},
	}

	for i, tt := range tests {
		res, err := m.GetChannels(tt.AppID)
		if err != nil {
			t.Fatalf("#%d: Get Channels Failed: %v", i, err)
		}
		if !reflect.DeepEqual(res, tt) {
			t.Errorf("#%d:\nGetChannels(appid)=%s\nwant=%s\n", i,
				res, tt)
		}
	}
}

func prepareChannels(m AppChannelService) {
	m.InsertChannel("a2", "c2")
	m.InsertChannel("a1", "c1")
	m.InsertChannel("a1", "c1")
	m.InsertChannel("a2", "c3")
	m.InsertChannel("a2", "c2")
	m.InsertChannel("a3", "c4")
}
