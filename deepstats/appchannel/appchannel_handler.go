package appchannel

import (
	"encoding/json"
	"net/http"
	"net/url"
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
	mux.Handle(endpoint, newAppChannelHandler(NewMongoAppChannelService(coll)))
}

type appChannelHandler struct {
	ags AppChannelService
}

func newAppChannelHandler(acs AppChannelService) http.Handler {
	return &appChannelHandler{acs}
}

func AppChannelPath(appID string) string {
	return path.Join(api.AppChannelPrefix, appID, "channels")
}

// Serves following endpoints:
// - GET /appchannels/:appid/channel
// - DELETE /appchannels/:appid/channel
func (ch *appChannelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Infof("Request to AppChannel Handler, %s", r.URL.Path)
	if !httputil.AllowMethod(w, r.Method, "GET", "DELETE") {
		return
	}
	start := time.Now()
	defer in.PromCounter.HTTPGETDuration(start)

	if len(r.URL.Path) <= len(api.AppChannelPrefix) {
		log.Errorf("Request Path is invalid, request path=%s, want=/appchannels/:appid/channels.", r.URL.Path)
		httputil.WriteHTTPError(w, api.ErrPathNotFound)
		return
	}
	fields := strings.Split(string(r.URL.Path[len(api.AppChannelPrefix):]), "/")
	if len(fields) != 2 {
		log.Errorf("Request Path is invalid, request path=%s, want=/appchannels/:appid/channels.", r.URL.Path)
		httputil.WriteHTTPError(w, api.ErrPathNotFound)
		return
	}
	if fields[1] != "channels" {
		log.Errorf("Request Path is invalid, request path=%s, want=/appchannels/:appid/channels.", r.URL.Path)
		httputil.WriteHTTPError(w, api.ErrPathNotFound)
		return
	}

	switch r.Method {
	case "GET":
		log.Infof("Request Detail, AppID=%s", fields[0])
		channels, err := ch.ags.GetChannels(fields[0])
		if err != nil {
			log.Fatalf("Query AppChannel service for channel list of app failed! Err Msg=%v", err)
			httputil.WriteHTTPError(w, api.ErrInternalServer)
			return
		}
		en := json.NewEncoder(w)
		if err := en.Encode(channels); err != nil {
			log.Errorf("Encode JSON failed! Err Msg=%v", err)
			httputil.WriteHTTPError(w, api.ErrInternalServer)
			return
		}
	case "DELETE":
		if len(r.URL.Query()["channel"]) == 0 {
			log.Errorf("Request param is invalid, require specific channel.")
			httputil.WriteHTTPError(w, api.ErrBadRequestBody)
			return
		}

		channel, err := url.QueryUnescape(r.URL.Query()["channel"][0])
		if err != nil || channel == "" {
			log.Errorf("Unescape channel failed! Err Msg=%v", err)
			httputil.WriteHTTPError(w, api.ErrBadRequestBody)
			return
		}
		if err := ch.ags.DeleteChannel(fields[0], channel); err != nil && !strings.Contains(err.Error(), "not found") {
			log.Errorf("Delete channel failed! Err Msg=%v", err)
			httputil.WriteHTTPError(w, api.ErrInternalServer)
			return
		}
	}
}
