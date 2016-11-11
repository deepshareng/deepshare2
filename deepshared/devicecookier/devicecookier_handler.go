package devicecookier

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

type deviceCookierHandler struct {
	dc           DeviceCookier
	ep           EventProducer
	cookiePrefix string
}

func AddHandler(mux *http.ServeMux, endpoint string, db storage.SimpleKV, p messaging.Producer) {

	mux.Handle(endpoint, newDeviceCookieHandler(NewDeviceCookier(db), p, endpoint))
}

func newDeviceCookieHandler(s DeviceCookier, p messaging.Producer, pfx string) *deviceCookierHandler {
	return &deviceCookierHandler{
		dc:           s,
		ep:           NewEventProducer(p),
		cookiePrefix: pfx,
	}
}

func NewDeviceCookieTestHandler(pfx string) *deviceCookierHandler {
	s := NewDeviceCookier(storage.NewInMemSimpleKV())
	return &deviceCookierHandler{
		dc:           s,
		ep:           messaging.NewSimpleProducer(bytes.NewBuffer(nil)),
		cookiePrefix: pfx,
	}
}

func (dch *deviceCookierHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !httputil.AllowMethod(w, r.Method, "GET", "PUT") {
		return
	}

	start := time.Now()
	fields := strings.Split(r.URL.Path[len(dch.cookiePrefix):], "/")
	if len(fields) != 2 {
		log.Error("Path not supported:", r.URL.Path, "method:", r.Method)
		httputil.WriteHTTPError(w, api.ErrPathNotFound)
		return
	}

	switch r.Method {
	case "GET":
		defer in.PrometheusForDeviceCookier.HTTPGetDuration(start)
		switch fields[0] {
		case ApiDevicesPrefix:
			cookieID := fields[1]
			// This include both sender info and inapp_data
			did, err := dch.dc.GetDeviceID(cookieID)
			if err != nil {
				log.Error(api.ErrCookieNotFound)
				httputil.WriteHTTPError(w, api.ErrCookieNotFound)
				return
			}
			en := json.NewEncoder(w)
			resp := DeviceInfo{
				DeviceID: did,
			}
			if err := en.Encode(resp); err != nil {
				log.Errorf("api: failed to encode devicecookie response to %s", r.RemoteAddr)
				httputil.WriteHTTPError(w, api.ErrInternalServer)
				return
			}
			log.Debugf("GetDevice based on cookieID succeed, cookieID: %s, deviceID: %s\n", did, cookieID)

		case ApiCookiesPrefix:
			deviceID := fields[1]
			// This include both sender info and inapp_data
			cid, err := dch.dc.GetCookieID(deviceID)
			if err != nil {
				log.Error(api.ErrCookieNotFound)
				httputil.WriteHTTPError(w, api.ErrCookieNotFound)
				return
			}
			en := json.NewEncoder(w)
			resp := CookieInfo{
				CookieID: cid,
			}
			log.Debugf("GetCookieID based on deviceID succeed, deviceID: %s, cookieID: %s\n", deviceID, cid)
			if err := en.Encode(resp); err != nil {
				log.Errorf("api: failed to encode devicecookie response to %s", r.RemoteAddr)
				httputil.WriteHTTPError(w, api.ErrInternalServer)
				return
			}
		default:
			log.Error("Unknown path:", r.URL.Path, "method:", r.Method, fields[0])
			httputil.WriteHTTPError(w, api.ErrPathNotFound)
			return
		}

		// TODO: should we produce event here? Is this user triggered?

	case "PUT":

		defer in.PrometheusForDeviceCookier.HTTPPutDuration(start)
		switch fields[0] {
		case ApiCookiesPrefix:
			deviceID := fields[1]
			decoder := json.NewDecoder(r.Body)
			cookieInfo := new(CookieInfo)
			if err := decoder.Decode(cookieInfo); err != nil {
				log.Error(api.ErrBadJSONBody)
				httputil.WriteHTTPError(w, api.ErrBadJSONBody)
				return
			}
			if err := dch.dc.SaveCookieID(deviceID, cookieInfo.CookieID); err != nil {
				log.Error("BindDeviceCookie by uniqueID failed, err:", err)
				httputil.WriteHTTPError(w, api.ErrInternalServer)
				return
			}

		case ApiDevicesPrefix:
			cookieID := fields[1]
			decoder := json.NewDecoder(r.Body)
			deviceInfo := new(DeviceInfo)
			if err := decoder.Decode(deviceInfo); err != nil {
				log.Error(api.ErrBadJSONBody)
				httputil.WriteHTTPError(w, api.ErrBadJSONBody)
				return
			}
			if err := dch.dc.SaveDeviceID(cookieID, deviceInfo.DeviceID); err != nil {
				log.Error("BindCookieDevice by uniqueID failed, err:", err)
				httputil.WriteHTTPError(w, api.ErrInternalServer)
				return
			}

		default:
			log.Error("Unknown path:", r.URL.Path, "method:", r.Method, fields[0])
			httputil.WriteHTTPError(w, api.ErrPathNotFound)
			return
		}
	}
}
