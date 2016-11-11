package device_stat

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/MISingularity/deepshare2/api"
	in "github.com/MISingularity/deepshare2/pkg/instrumentation"
	"github.com/MISingularity/deepshare2/pkg/storage"

	"github.com/MISingularity/deepshare2/pkg/httputil"
	"github.com/MISingularity/deepshare2/pkg/log"
)

func AddHandler(mux *http.ServeMux, endpoint string, skv storage.SimpleKV, prefix string) {
	mux.Handle(endpoint, newDeviceStatHandler(NewGeneralDeviceStat(skv, prefix)))
}

type deviceStatHandler struct {
	ds DeviceStatService
}

func newDeviceStatHandler(ds DeviceStatService) http.Handler {
	return &deviceStatHandler{ds}
}

func DeviceStatPath(os string) string {
	return path.Join(api.DeviceStatPrefix, os)
}

// Serves following endpoints:
// - GET /device-stat/:os
func (dh *deviceStatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Infof("Request to DeviceStat Handler, %s", r.URL.Path)
	if !httputil.AllowMethod(w, r.Method, "GET") {
		return
	}
	start := time.Now()
	defer in.PrometheusForAppInfo.HTTPGetDuration(start)

	if len(r.URL.Path) <= len(api.DeviceStatPrefix) {
		log.Errorf("Request Path is invalid, request path=%s, want=%s:os.", r.URL.Path, api.DeviceStatPrefix)
		httputil.WriteHTTPError(w, api.ErrPathNotFound)
		return
	}
	fields := strings.Split(string(r.URL.Path[len(api.AppChannelPrefix):]), "/")
	if len(fields) != 1 {
		log.Errorf("Request Path is invalid, request path=%s, want=%s:os.", r.URL.Path, api.DeviceStatPrefix)
		httputil.WriteHTTPError(w, api.ErrPathNotFound)
		return
	}

	switch r.Method {
	case "GET":
		log.Infof("Request Detail, AppID=%s", fields[0])
		deviceSum, err := dh.ds.Count(fields[0])
		if err != nil {
			log.Fatalf("Query DeviceStat service failed! Err Msg=%v", err)
			httputil.WriteHTTPError(w, api.ErrInternalServer)
			return
		}
		en := json.NewEncoder(w)
		if err := en.Encode(struct{ Value int64 }{deviceSum}); err != nil {
			log.Errorf("Encode JSON failed! Err Msg=%v", err)
			httputil.WriteHTTPError(w, api.ErrInternalServer)
			return
		}
	}
}
