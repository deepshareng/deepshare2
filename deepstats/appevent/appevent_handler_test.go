package appevent

import (
	"net/http"
	"path"
	"testing"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepstats"
	"github.com/MISingularity/deepshare2/pkg/testutil"
)

func TestEventGet(t *testing.T) {
	dbName := "test_appeventhandler"
	collName := "appeventhandler"

	session := deepstats.MustCreateMongoSession(deepstats.MongoTestingAddr)
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()
	ae := NewMongoAppEventService(c)
	prepareEvents(ae)

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
			`{"AppID":"a1","Events":["match/e1","dsaction/e2","counter/e1"]}` +
				"\n",
		},
		{ //good request, filter in "install"
			"",
			"a2",

			http.StatusOK,
			`{"AppID":"a2","Events":["counter/e1"]}` +
				"\n",
		},
		{ //wrong path request
			"apps/",
			"non-exist-app-id",

			http.StatusNotFound,
			api.ErrPathNotFound.Error() + "\n",
		},
		{ //wrong path request
			"apps/events",
			"non-exist-app-id",

			http.StatusNotFound,
			api.ErrPathNotFound.Error() + "\n",
		},
		{ //wrong path request
			"apps/non/event",
			"non-exist-app-id",

			http.StatusNotFound,
			api.ErrPathNotFound.Error() + "\n",
		},
	}

	handler := newAppEventHandler(ae)
	// Do some GET requests with different channel IDs and event filters.
	for i, tt := range tests {
		var url string
		if tt.requestPath == "" {
			url = "http://" + path.Join("example.com", AppEventPath(tt.appID))
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
