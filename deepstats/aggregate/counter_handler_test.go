package aggregate

import (
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepstats"
	"github.com/MISingularity/deepshare2/pkg/testutil"

	"testing"

	"github.com/MISingularity/deepshare2/pkg/storage"
)

func TestCounterGet(t *testing.T) {
	dbName := "test_aggregate_counter"

	session := deepstats.MustCreateMongoSession(deepstats.MongoTestingAddr)
	defer session.Close()
	/*defer */ session.DB(dbName).DropDatabase()
	redisDB := storage.NewInMemSimpleKV()
	hourAgs := NewHourAggregateService(redisDB, session, dbName, "hour", "1", "")
	dayAgs := NewDayAggregateService(redisDB, session, dbName, "day", "2", "")
	totalAgs := NewTotalAggregateService(redisDB, session, dbName, "total", "3", "")
	intervalRefreshToMongoSec = 0
	prepareEvents(hourAgs)
	prepareEvents(dayAgs)
	prepareEvents(totalAgs)
	hourAgs.refreshToMongo()
	dayAgs.refreshToMongo()
	totalAgs.refreshToMongo()

	type counters struct {
		Counters []*AggregateResult
	}
	tests := []struct {
		requestPath string
		channelID   string
		params      string

		wcode int
		wbody string
	}{
		// hour test
		{ //#0 good request, no filters
			"",
			"c1",
			fmt.Sprintf(
				"?appid=w1&gran=h&start=%d&end=%d&event=open",
				time.Date(2015, time.January, 1, 18, 0, 0, 0, time.Local).Unix(),
				time.Date(2015, time.January, 2, 0, 0, 0, 0, time.Local).Unix(),
			),

			http.StatusOK,
			`{"Counters":[` +
				`{"Event":"open","Counts":[{"Count":1},{"Count":0},{"Count":2},{"Count":0},{"Count":0},{"Count":0}]}]}` +
				"\n",
		},

		// day test
		{ //#1 good request, no filters
			"",
			"c1",
			fmt.Sprintf(
				"?appid=w1&start=%d&end=%d&event=open",
				time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				time.Date(2015, time.January, 8, 0, 0, 0, 0, time.Local).Unix(),
			),

			http.StatusOK,
			`{"Counters":[` +
				`{"Event":"open","Counts":[{"Count":3},{"Count":1},{"Count":0},{"Count":0},{"Count":0},{"Count":0},{"Count":1}]}]}` +
				"\n",
		},
		{ //#2 good request, filter in "install"
			"",
			"c1",
			fmt.Sprintf(
				"?appid=w1&event=install&start=%d&end=%d",
				time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				time.Date(2015, time.January, 8, 0, 0, 0, 0, time.Local).Unix(),
			),

			http.StatusOK,
			`{"Counters":[` +
				`{"Event":"install","Counts":[{"Count":1},{"Count":0},{"Count":0},{"Count":0},{"Count":0},{"Count":0},{"Count":0}]}]}` +
				"\n",
		},

		// week test
		{ //#3 good request
			"",
			"c1",
			fmt.Sprintf(
				"?appid=w1&gran=w&start=%d&end=%d&event=open",
				time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				time.Date(2015, time.January, 9, 0, 0, 0, 0, time.Local).Unix(),
			),

			http.StatusOK,
			`{"Counters":[` +
				`{"Event":"open","Counts":[{"Count":5},{"Count":0}]}]}` +
				"\n",
		},
		{ //#4 good request, filter in "open"
			"",
			"c1",
			fmt.Sprintf(
				"?appid=w3&gran=w&start=%d&limit=3&event=open",
				time.Date(2014, time.December, 17, 0, 0, 0, 0, time.Local).Unix(),
			),

			http.StatusOK,
			`{"Counters":[` +
				`{"Event":"open","Counts":[{"Count":1},{"Count":1},{"Count":3}]}]}` +
				"\n",
		},
		{ //#5 good request, filter in "install"
			"",
			"c1",
			fmt.Sprintf(
				"?appid=w1&gran=w&event=install&start=%d&limit=3&event=install",
				time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
			),

			http.StatusOK,
			`{"Counters":[` +
				`{"Event":"install","Counts":[{"Count":1},{"Count":0},{"Count":0}]}]}` +
				"\n",
		},

		// Month test
		{ //#6 good request, filter in "open"
			"",
			"c1",
			fmt.Sprintf(
				"?appid=w3&gran=m&start=%d&limit=3&event=open",
				time.Date(2014, time.December, 2, 0, 0, 0, 0, time.Local).Unix(),
			),

			http.StatusOK,
			`{"Counters":[` +
				`{"Event":"open","Counts":[{"Count":3},{"Count":2},{"Count":0}]}]}` +
				"\n",
		},

		// total test
		{ //#7 good request, no filters
			"",
			"c1",
			"?appid=w1&gran=t&limit=1&event=open",

			http.StatusOK,
			`{"Counters":[` +
				`{"Event":"open","Counts":[{"Count":5}]}]}` +
				"\n",
		},
		{ //#8 good request, filter in "install"
			"",
			"c1",
			"?appid=w1&gran=t&limit=1&event=install",

			http.StatusOK,
			`{"Counters":[` +
				`{"Event":"install","Counts":[{"Count":1}]}]}` +
				"\n",
		},
		{ //#9 empty request
			"",
			"non-exist-channel-id",
			fmt.Sprintf(
				"?appid=w1&start=%d&end=%d",
				time.Date(2014, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
			),

			http.StatusOK,
			`{"Counters":[]}` + "\n",
		},
		{ // test os
			"",
			"c11",
			fmt.Sprintf(
				"?appid=w33&start=%d&end=%d&os=%s",
				time.Date(2020, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				time.Date(2020, time.January, 2, 0, 0, 0, 0, time.Local).Unix(),
				"ios",
			),
			http.StatusOK,
			`{"Counters":[{"Event":"version/install","Counts":[{"Count":2}]}]}` + "\n",
		},
		{ // test os
			"",
			"c11",
			fmt.Sprintf(
				"?appid=w33&start=%d&end=%d&os=%s",
				time.Date(2020, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				time.Date(2020, time.January, 2, 0, 0, 0, 0, time.Local).Unix(),
				"android",
			),
			http.StatusOK,
			`{"Counters":[{"Event":"version/install","Counts":[{"Count":10}]}]}` + "\n",
		},
		{ // test os
			"",
			"c11",
			fmt.Sprintf(
				"?appid=w33&start=%d&end=%d&os=%s",
				time.Date(2020, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				time.Date(2020, time.January, 2, 0, 0, 0, 0, time.Local).Unix(),
				"other",
			),
			http.StatusOK,
			`{"Counters":[{"Event":"version/install","Counts":[{"Count":1}]}]}` + "\n",
		},
		{ // test os
			"",
			"c11",
			fmt.Sprintf(
				"?appid=w33&start=%d&end=%d&os=%s",
				time.Date(2020, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
				time.Date(2020, time.January, 2, 0, 0, 0, 0, time.Local).Unix(),
				"",
			),
			http.StatusOK,
			`{"Counters":[{"Event":"version/install","Counts":[{"Count":13}]}]}` + "\n",
		},

		{ //wrong path request
			"v2/channels/",
			"non-exist-channel-id",
			"",

			http.StatusNotFound,
			api.ErrPathNotFound.Error() + "\n",
		},
		{ //wrong path request
			"v2/channels/counters",
			"non-exist-channel-id",
			"",

			http.StatusNotFound,
			api.ErrPathNotFound.Error() + "\n",
		},
		{ //wrong path request
			"v2/channels/non/counter",
			"non-exist-channel-id",
			"",

			http.StatusNotFound,
			api.ErrPathNotFound.Error() + "\n",
		},
	}

	handler := newCounterHandler(hourAgs, dayAgs, totalAgs)
	// Do some GET requests with different channel IDs and event filters.
	for i, tt := range tests {
		var url string
		if tt.requestPath == "" {
			url = "http://" + path.Join("example.com", CounterPath(tt.channelID)) + tt.params
		} else {
			url = "http://" + path.Join("exmaple.com", tt.requestPath) + tt.params
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
