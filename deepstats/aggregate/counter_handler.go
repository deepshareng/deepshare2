package aggregate

import (
	"encoding/json"
	"math"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/pkg/httputil"
	in "github.com/MISingularity/deepshare2/pkg/instrumentation"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"gopkg.in/mgo.v2"
)

type counterHandler struct {
	hourAgs, dayAgs, totalAgs AggregateService
}

func AddHandler(mux *http.ServeMux, endpoint string, redisDB storage.SimpleKV, sessionHour, sessionDay, sessionTotal *mgo.Session, dbNameHour, dbNameDay, dbNameTotal string, hourCollNamePrefix, dayCollNamePrefix, totalCollNamePrefix string) {
	mux.Handle(
		endpoint,
		newCounterHandler(
			NewHourAggregateService(redisDB, sessionHour, dbNameHour, hourCollNamePrefix, "", ""),
			NewDayAggregateService(redisDB, sessionDay, dbNameDay, dayCollNamePrefix, "", ""), //TODO
			NewTotalAggregateService(redisDB, sessionTotal, dbNameTotal, totalCollNamePrefix, "", ""),
		),
	)
}

func newCounterHandler(hourAgs, dayAgs, totalAgs AggregateService) http.Handler {
	return &counterHandler{hourAgs, dayAgs, totalAgs}
}

func CounterPath(channelID string) string {
	return path.Join(api.ChannelPrefix, channelID, "counters")
}

// Serves following endpoints:
// - GET /channels/:channel_id/counters
// 		- parameters: "?appid=...&gran=[d,w,y]&limit=10&event=install&event=..."
// 		- response counts follow chronological order on the basis of event.
// return inverted order, the smaller index, the recent data

// Support several format rules requesting data, for pattern, we check every rule in turn.
// If there is more than one match, we will return the first one matched.
// 1. start=x&end=y
// 2. start=x&limit=y
// 3. end=x&limit=y
// 4. limit=x

func (ch *counterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Infof("Request to Aggregate CounterHandler, %s", r.URL.Path)
	if !httputil.AllowMethod(w, r.Method, "GET") {
		log.Errorf("Method is not permitted, request method=%s.", r.Method)
		return
	}
	startTime := time.Now()
	defer in.PromCounter.HTTPGETDuration(startTime)

	if len(r.URL.Path) <= len(api.AppChannelPrefix) {
		log.Errorf("Request Path is invalid, request path=%s, want=/channels/:channel_id/counters.", r.URL.Path)
		httputil.WriteHTTPError(w, api.ErrPathNotFound)
		return
	}
	fields := strings.Split(string(r.URL.Path[len(api.ChannelPrefix):]), "/")
	if len(fields) != 2 {
		log.Errorf("Request Path is invalid, request path=%s, want=/channels/:channel_id/counters.", r.URL.Path)
		httputil.WriteHTTPError(w, api.ErrPathNotFound)
		return
	}
	if fields[1] != "counters" {
		log.Errorf("Request Path is invalid, request path=%s, want=/channels/:channel_id/counters.", r.URL.Path)
		httputil.WriteHTTPError(w, api.ErrPathNotFound)
		return
	}
	var start time.Time
	var err, err2 error
	gran := time.Hour * 24
	baseGran := time.Hour * 24
	var ags AggregateService
	ags = ch.dayAgs
	limit := -1
	log.Debugf("Request Params, appid=%v, start=%v, end=%v, gran=%v, limit=%v", r.URL.Query()["appid"], r.URL.Query()["start"], r.URL.Query()["end"], r.URL.Query()["gran"], r.URL.Query()["limit"])

	if len(r.URL.Query()["appid"]) == 0 {
		log.Error("Lack appid attribute.")
		httputil.WriteHTTPError(w, api.ErrBadRequestBody)
		return
	}
	appid, err := url.QueryUnescape(r.URL.Query()["appid"][0])
	if err != nil {
		log.Error("Unescape appid failed.")
		httputil.WriteHTTPError(w, api.ErrBadRequestBody)
		return
	}

	if len(r.URL.Query()["gran"]) != 0 {
		switch r.URL.Query()["gran"][0] {
		case "h":
			gran = time.Hour
			baseGran = time.Hour
			ags = ch.hourAgs
		case "d":
			gran = time.Hour * 24
			baseGran = time.Hour * 24
			ags = ch.dayAgs
		case "w":
			gran = time.Hour * 24 * 7
			baseGran = time.Hour * 24
			ags = ch.dayAgs
		case "m":
			gran = time.Hour * 24 * 30
			baseGran = time.Hour * 24
			ags = ch.dayAgs
		case "t":
			gran = time.Hour
			baseGran = time.Hour
			ags = ch.totalAgs
			start = time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)
			limit = 1
		}
	}

	if len(r.URL.Query()["appid"]) == 0 {
		log.Error("Lack appid attribute.")
		httputil.WriteHTTPError(w, api.ErrBadRequestBody)
		return
	}
	appid = r.URL.Query()["appid"][0]
	if len(r.URL.Query()["start"]) != 0 && len(r.URL.Query()["end"]) != 0 {
		// start and end
		startNum, err1 := strconv.ParseInt(r.URL.Query()["start"][0], 10, 64)
		endNum, err2 := strconv.ParseInt(r.URL.Query()["end"][0], 10, 64)
		start = time.Unix(startNum, 0)
		end := time.Unix(endNum, 0)
		if err1 != nil || err2 != nil {
			log.Error("Convert start/end attribute failed!")
			httputil.WriteHTTPError(w, api.ErrBadRequestBody)
			return
		}
		limit = int(math.Ceil(float64(end.Sub(start)) / float64(gran)))
	} else if len(r.URL.Query()["start"]) != 0 && len(r.URL.Query()["limit"]) != 0 {
		// start and limit
		startNum, err1 := strconv.ParseInt(r.URL.Query()["start"][0], 10, 64)
		limit, err2 = strconv.Atoi(r.URL.Query()["limit"][0])
		start = time.Unix(startNum, 0)
		if err1 != nil || err2 != nil {
			log.Error("Convert start/limit attribute failed!")
			httputil.WriteHTTPError(w, api.ErrBadRequestBody)
			return
		}
	} else if len(r.URL.Query()["end"]) != 0 && len(r.URL.Query()["limit"]) != 0 {
		// end and limit
		endNum, err1 := strconv.ParseInt(r.URL.Query()["end"][0], 10, 64)
		limit, err2 = strconv.Atoi(r.URL.Query()["limit"][0])
		end := time.Unix(endNum, 0)
		if err1 != nil || err2 != nil {
			log.Error("Convert end/limit1 attribute failed!")
			httputil.WriteHTTPError(w, api.ErrBadRequestBody)
			return
		}
		end = ags.ConvertTimeToGranularity(end)
		start = end.Add(time.Duration(-limit)*gran + baseGran)
	} else if len(r.URL.Query()["limit"]) != 0 {
		// only limit
		end := ags.ConvertTimeToGranularity(time.Now()).Add(baseGran)
		limit, err = strconv.Atoi(r.URL.Query()["limit"][0])
		if err != nil {
			log.Error("Convert end/limit1 attribute failed!")
			httputil.WriteHTTPError(w, api.ErrBadRequestBody)
			return
		}
		start = end.Add(time.Duration(-limit) * gran)
	} else {
		log.Error("Lack duration(start/end/limit) attribute!")
		httputil.WriteHTTPError(w, api.ErrBadRequestBody)
		return
	}
	eventFilters := r.URL.Query()["event"]

	// os extraction
	os := ""
	if len(r.URL.Query()["os"]) != 0 {
		os = r.URL.Query()["os"][0]
	}

	log.Infof("Request Detail, AppID=%s, Channel=%s, Event=%v, Start=%s, Granularity=%s, Limit=%d, Os=%s", appid, fields[0], eventFilters, start.String(), gran.String(), limit, os)
	aggrs, err := ags.QueryDuration(appid, fields[0], eventFilters, start, gran, limit, os)
	if err != nil {
		log.Fatalf("Query Aggregate service for result of duration failed! Err Msg=%v", err)
		httputil.WriteHTTPError(w, api.ErrInternalServer)
		return
	}
	res := &struct {
		Counters []*AggregateResult
	}{
		aggrs,
	}

	en := json.NewEncoder(w)
	if err := en.Encode(res); err != nil {
		log.Errorf("Encode JSON failed! Err Msg=%v", err)
		httputil.WriteHTTPError(w, api.ErrInternalServer)
		return
	}
}
