package appchannel

import (
	"net/http"
	"path"
	"testing"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepstats"
	"github.com/MISingularity/deepshare2/deepstats/aggregate"
	"github.com/MISingularity/deepshare2/pkg/testutil"
)

func TestChannelGet(t *testing.T) {
	dbName := "test_appchannel"
	collName := "counter"

	session := deepstats.MustCreateMongoSession(deepstats.MongoTestingAddr)
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()
	ap := NewMongoAppChannelService(c)

	prepareChannels(ap)

	type counters struct {
		Counters []*aggregate.AggregateResult
	}
	tests := []struct {
		requestPath string
		appID       string

		wcode int
		wbody string
	}{
		{ //good request, no filters
			"",
			"a1",

			http.StatusOK,
			`{"AppID":"a1","Channels":["c1"]}` +
				"\n",
		},
		{ //good request, filter in "install"
			"",
			"a2",

			http.StatusOK,
			`{"AppID":"a2","Channels":["c2","c3"]}` +
				"\n",
		},
		{ //wrong path request
			"apps/",
			"non-exist-app-id",

			http.StatusNotFound,
			api.ErrPathNotFound.Error() + "\n",
		},
		{ //wrong path request
			"apps/channels",
			"non-exist-app-id",

			http.StatusNotFound,
			api.ErrPathNotFound.Error() + "\n",
		},
		{ //wrong path request
			"apps/non/channel",
			"non-exist-app-id",

			http.StatusNotFound,
			api.ErrPathNotFound.Error() + "\n",
		},
	}

	handler := newAppChannelHandler(ap)
	// Do some GET requests with different channel IDs and event filters.
	for i, tt := range tests {
		var url string
		if tt.requestPath == "" {
			url = "http://" + path.Join("example.com", AppChannelPath(tt.appID))
		} else {
			url = "http://" + path.Join("exmaple.com", tt.requestPath)
		}
		w := testutil.HandleWithBody(handler, "GET", url, "")

		if w.Code != tt.wcode {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.wcode)
		}

		if string(w.Body.Bytes()) != tt.wbody {
			t.Errorf("#%d: HTTP response body = %q, want = %q", i, string(w.Body.Bytes()), tt.wbody)
		}
	}

}

func TestChannelDelete(t *testing.T) {
	dbName := "test_appchannel_delete"
	collName := "appchannel"

	session := deepstats.MustCreateMongoSession(deepstats.MongoTestingAddr)
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()
	ap := NewMongoAppChannelService(c)

	prepareChannels(ap)

	type counters struct {
		Counters []*aggregate.AggregateResult
	}
	tests := []struct {
		requestPath string
		appID       string
		params      string

		wcode int
		wbody string
	}{
		{ //good request, no filters
			"",
			"a1",
			"?channel=c1",

			http.StatusOK,
			`{"AppID":"a1","Channels":[]}` +
				"\n",
		},
		{ //good request, filter in "install"
			"",
			"a2",
			"?channel=c2",

			http.StatusOK,
			`{"AppID":"a2","Channels":["c3"]}` +
				"\n",
		},
	}

	handler := newAppChannelHandler(ap)
	// Do some GET requests with different channel IDs and event filters.
	for i, tt := range tests {
		var deleteurl, geturl string
		if tt.requestPath == "" {
			deleteurl = "http://" + path.Join("example.com", AppChannelPath(tt.appID)) + tt.params
			geturl = "http://" + path.Join("example.com", AppChannelPath(tt.appID))
		} else {
			deleteurl = "http://" + path.Join("exmaple.com", tt.requestPath) + tt.params
			geturl = "http://" + path.Join("example.com", AppChannelPath(tt.appID))
		}
		deletew := testutil.HandleWithBody(handler, "DELETE", deleteurl, "")
		if deletew.Code != http.StatusOK {
			t.Errorf("#%d: DELETE HTTP status code = %d, want = %d", i, deletew.Code, http.StatusOK)
		}
		getw := testutil.HandleWithBody(handler, "GET", geturl, "")
		if getw.Code != tt.wcode {
			t.Errorf("#%d: GET HTTP status code = %d, want = %d", i, getw.Code, tt.wcode)
		}

		if string(getw.Body.Bytes()) != tt.wbody {
			t.Errorf("#%d: GET HTTP response body = %q, want = %q", i, string(getw.Body.Bytes()), tt.wbody)
		}
	}

}
