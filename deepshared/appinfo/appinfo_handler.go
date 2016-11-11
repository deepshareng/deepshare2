package appinfo

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"bytes"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/pkg/httputil"
	in "github.com/MISingularity/deepshare2/pkg/instrumentation"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/MISingularity/deepshare2/pkg/storage"
)

type appInfoHandler struct {
	ais *appInfoService
	p   messaging.Producer
	pfx string
}

// Used for unit testing handler core logic.
func AddHandler(mux *http.ServeMux, endpoint string, db storage.SimpleKV, mp messaging.Producer) {
	mux.Handle(endpoint, NewAppInfoHandler(db, mp, endpoint))
}

func NewAppInfoHandler(db storage.SimpleKV, pp messaging.Producer, ppfx string) http.Handler {
	return &appInfoHandler{
		ais: NewAppInfoService(db),
		p:   pp,
		pfx: ppfx,
	}
}

func NewAppInfoTestHandler(pfx string) http.Handler {
	return &appInfoHandler{
		ais: NewAppInfoService(storage.NewInMemSimpleKV()),
		p:   messaging.NewSimpleProducer(bytes.NewBuffer(nil)),
		pfx: pfx,
	}
}

func (aih *appInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if !httputil.AllowMethod(w, r.Method, "GET", "PUT") {
		return
	}

	start := time.Now()

	fields := strings.Split(r.URL.Path[len(aih.pfx):], "/")

	if len(fields) != 1 {
		log.Error("appInfoHandler; path not found:", r.URL.Path, aih.pfx)
		httputil.WriteHTTPError(w, api.ErrPathNotFound)
		return
	}

	appID := fields[0]

	switch r.Method {
	case "GET":
		defer in.PrometheusForAppInfo.HTTPGetDuration(start)

		// This include both sender info and inapp_data
		appInfo, err := aih.ais.GetAppInfo(appID)
		if err != nil {
			log.Error("GetAppInfo err:", err)
			httputil.WriteHTTPError(w, api.ErrAppIDNotFound)
			return
		}
		en := json.NewEncoder(w)
		if err := en.Encode(appInfo); err != nil {
			// TODO: use a logger pkg and change this to debug level
			log.Errorf("api: failed to encode match response to %s", r.RemoteAddr)
		}
		// TODO: should we produce event here? Is this user triggered?

	case "PUT":
		defer in.PrometheusForAppInfo.HTTPPutDuration(start)
		decoder := json.NewDecoder(r.Body)
		req := new(AppInfo)
		if err := decoder.Decode(req); err != nil {
			httputil.WriteHTTPError(w, api.ErrBadJSONBody)
			return
		}

		err := aih.ais.SetAppInfo(appID, req)
		if err != nil {
			log.Error("Failed to set appinfo, err:", err)
			httputil.WriteHTTPError(w, api.ErrInternalServer)
			return
		}

		if req.ShortID != "" {
			if err := SetAppIDShortIDPair(aih.ais.db, req.ShortID, appID); err != nil {
				log.Error("Faild to SetAppIDShortIDPair, err:", err)
				httputil.WriteHTTPError(w, api.ErrInternalServer)
				return
			}
		}
	}
}
