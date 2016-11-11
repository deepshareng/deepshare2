package appevent

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/MISingularity/deepshare2/api"
	in "github.com/MISingularity/deepshare2/pkg/instrumentation"

	"github.com/MISingularity/deepshare2/pkg/httputil"
	"github.com/MISingularity/deepshare2/pkg/log"
	"gopkg.in/mgo.v2"
)

func AddHandler(mux *http.ServeMux, endpoint string, coll *mgo.Collection) {
	mux.Handle(endpoint, newAppEventHandler(NewMongoAppEventService(coll)))
}

type appEventHandler struct {
	aes AppEventService
}

func newAppEventHandler(aes AppEventService) http.Handler {
	return &appEventHandler{aes}
}

func AppEventPath(appID string) string {
	return path.Join(api.AppEventPrefix, appID, "events")
}

// Serves following endpoints:
// - GET /apps/:appid/events
func (ch *appEventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Infof("Request to AppEvent Handler, %s", r.URL.Path)
	if !httputil.AllowMethod(w, r.Method, "GET") {
		return
	}
	start := time.Now()
	defer in.PromCounter.HTTPGETDuration(start)

	if len(r.URL.Path) <= len(api.AppEventPrefix) {
		log.Errorf("Request Path is invalid, request path=%s, want=/appevents/:appid/events.", r.URL.Path)
		httputil.WriteHTTPError(w, api.ErrPathNotFound)
		return
	}
	fields := strings.Split(string(r.URL.Path[len(api.AppEventPrefix):]), "/")
	if len(fields) != 2 {
		log.Errorf("Request Path is invalid, request path=%s, want=/appevents/:appid/events.", r.URL.Path)
		httputil.WriteHTTPError(w, api.ErrPathNotFound)
		return
	}
	if fields[1] != "events" {
		log.Errorf("Request Path is invalid, request path=%s, want=/appevents/:appid/events.", r.URL.Path)
		httputil.WriteHTTPError(w, api.ErrPathNotFound)
		return
	}

	log.Infof("Request Detail, AppID=%s", fields[0])
	events, err := ch.aes.GetEvents(fields[0])
	if err != nil {
		log.Fatalf("Query AppEvent service for event list of app failed! Err Msg=%v", err)
		httputil.WriteHTTPError(w, api.ErrInternalServer)
		return
	}
	en := json.NewEncoder(w)
	if err := en.Encode(events); err != nil {
		log.Errorf("Encode JSON failed! Err Msg=%v", err)
		httputil.WriteHTTPError(w, api.ErrInternalServer)
		return
	}
}
