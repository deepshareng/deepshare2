package appcookiedevice

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

type appCookieDeviceHandler struct {
	acd             AppCookieDevice
	ep              EventProducer
	appcookiePrefix string
}

func AddHandler(mux *http.ServeMux, endpoint string, db storage.SimpleKV, p messaging.Producer) {
	mux.Handle(endpoint, newAppCookieDeviceHandler(NewAppCookieDevice(db), p, endpoint))
}

func newAppCookieDeviceHandler(s AppCookieDevice, p messaging.Producer, pfx string) *appCookieDeviceHandler {
	return &appCookieDeviceHandler{
		acd:             s,
		ep:              NewEventProducer(p),
		appcookiePrefix: pfx,
	}
}

func NewAppCookieDeviceTestHandler(pfx string) *appCookieDeviceHandler {
	s := NewAppCookieDevice(storage.NewInMemSimpleKV())
	return &appCookieDeviceHandler{
		acd:             s,
		ep:              messaging.NewSimpleProducer(bytes.NewBuffer(nil)),
		appcookiePrefix: pfx,
	}
}

func (acdh *appCookieDeviceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !httputil.AllowMethod(w, r.Method, "GET", "PUT", "POST") {
		return
	}
	start := time.Now()
	fields := strings.Split(r.URL.Path[len(acdh.appcookiePrefix):], "/")

	switch r.Method {
	// GET /v2/appcookiedevice/:appID/:cookieID
	// response: {"unique_id":"uuu"}
	case "GET":
		if len(fields) != 2 {
			httputil.WriteHTTPError(w, api.ErrPathNotFound)
			return
		}

		appID := fields[0]
		cookieId := fields[1]
		defer in.PrometheusForAppCookieDevice.HTTPGetDuration(start)

		// This include both sender info and inapp_data
		did, err := acdh.acd.GetDeviceID(appID, cookieId)
		if err != nil {
			log.Error(api.ErrCookieNotFound)
			httputil.WriteHTTPError(w, api.ErrCookieNotFound)
			return
		}
		en := json.NewEncoder(w)
		resp := DeviceInfo{
			UniqueId: did,
		}
		log.Debugf("appCookieDeviceHandler; GetUniqueID based on appId and cookieId succeed, uniqeId: %s, appId: %s, cookieID: %s\n", did, appID, cookieId)
		if err := en.Encode(resp); err != nil {
			log.Errorf("appCookieDeviceHandler; api: failed to encode response to %s", r.RemoteAddr)
			httputil.WriteHTTPError(w, api.ErrInternalServer)
			return
		}
		return
		// TODO: should we produce event here? Is this user triggered?

	// PUT /v2/appcookiedevice/:appID/:cookieID  {"unique_id":"uuu"}
	case "PUT":
		if len(fields) != 2 {
			httputil.WriteHTTPError(w, api.ErrPathNotFound)
			return
		}

		appID := fields[0]
		cookieId := fields[1]
		defer in.PrometheusForAppCookieDevice.HTTPPutDuration(start)
		decoder := json.NewDecoder(r.Body)
		req := new(DeviceInfo)
		if err := decoder.Decode(req); err != nil {
			log.Error(api.ErrBadJSONBody)
			httputil.WriteHTTPError(w, api.ErrBadJSONBody)
			return
		}

		err := acdh.acd.BindAppCookieToDevice(appID, cookieId, req.UniqueId)
		if err != nil {
			log.Error(api.ErrInternalServer)
			httputil.WriteHTTPError(w, api.ErrInternalServer)
			return
		}

		//Get paired cookie, and bind appid+pairedCookie -> uniqueID
		pairCookie, err := acdh.acd.getPairedCookie(cookieId)
		if err != nil {
			log.Error(api.ErrInternalServer)
			httputil.WriteHTTPError(w, api.ErrInternalServer)
			return
		}
		if pairCookie != "" {
			err := acdh.acd.BindAppCookieToDevice(appID, pairCookie, req.UniqueId)
			if err != nil {
				log.Error(api.ErrInternalServer)
				httputil.WriteHTTPError(w, api.ErrInternalServer)
				return
			}
		}
		log.Debugf("appCookieDeviceHandler; BindAppCookieToDevice succeed, uniqueID: %s, appId: %s, cookieID: %s, pairCookieID: %s\n", req.UniqueId, appID, cookieId, pairCookie)
		return

	case "POST":
		switch {
		// POST /v2/appcookiedevice/paircookies/appID {"cookie1":"c1","cookie2":"c2"}
		case len(fields) == 2 && fields[0] == PairCookiesPath:
			appID := fields[1]
			decoder := json.NewDecoder(r.Body)
			req := new(PostPairCookieBody)
			if err := decoder.Decode(req); err != nil {
				log.Error("appCookieDeviceHandler; Failed to decode cookie pair request, err:", err)
				httputil.WriteHTTPError(w, api.ErrBadJSONBody)
				return
			}
			if req.Cookie1 == "" || req.Cookie2 == "" {
				log.Error("appCookieDeviceHandler; try to pair empty cookies, cookies:", req.Cookie1, req.Cookie2)
				httputil.WriteHTTPError(w, api.ErrBadRequestBody)
				return
			}
			if err := acdh.acd.pairCookies(req.Cookie1, req.Cookie2); err != nil {
				log.Error("appCookieDeviceHandler; pairCookies failed, err:", err)
				httputil.WriteHTTPError(w, api.ErrInternalServer)
				return
			}

			{ // if appid+cookie1 -> uniqueID, then bind appid+cookie2 -> uniqueID
				deviceID, err := acdh.acd.GetDeviceID(appID, req.Cookie1)
				if err != nil {
					log.Error("appCookieDeviceHandler; GeGetDeviceID failed, err:", err)
					httputil.WriteHTTPError(w, api.ErrInternalServer)
					return
				}
				if deviceID != "" {
					if err := acdh.acd.BindAppCookieToDevice(appID, req.Cookie2, deviceID); err != nil {
						log.Error("appCookieDeviceHandler; BindAppCookieToDevice failed, err:", err)
						httputil.WriteHTTPError(w, api.ErrInternalServer)
						return
					}
				}
			}

			{ // if appid+cookie2 -> uniqueID, then bind appid+cookie1 -> uniqueID
				deviceID, err := acdh.acd.GetDeviceID(appID, req.Cookie2)
				if err != nil {
					log.Error("appCookieDeviceHandler; GetDeviceID failed, err:", err)
					httputil.WriteHTTPError(w, api.ErrInternalServer)
					return
				}
				if deviceID != "" {
					if err := acdh.acd.BindAppCookieToDevice(appID, req.Cookie1, deviceID); err != nil {
						log.Error("appCookieDeviceHandler; BindAppCookieToDevice failed, err:", err)
						httputil.WriteHTTPError(w, api.ErrInternalServer)
						return
					}
				}
			}
			return

		// POST /v2/appcookiedevice/refreshcookie=true  {"cookie":"c1","new_cookie":"c2"}
		case len(fields) == 1 && fields[0] == RefreshCookiePath:
			decoder := json.NewDecoder(r.Body)
			req := new(PostRefreshCookieBody)
			if err := decoder.Decode(req); err != nil {
				log.Error("appCookieDeviceHandler; Failed to decode refresh cookie request, err:", err)
				httputil.WriteHTTPError(w, api.ErrBadJSONBody)
				return
			}
			if req.Cookie == "" || req.NewCookie == "" {
				log.Error("appCookieDeviceHandler; try to refresh empty cookie:", req.Cookie, req.NewCookie)
				httputil.WriteHTTPError(w, api.ErrBadRequestBody)
				return
			}
			if err := acdh.acd.RefreshCookie(req.Cookie, req.NewCookie); err != nil {
				log.Error("appCookieDeviceHandler; RefreshCookie failed, err:", err)
				httputil.WriteHTTPError(w, api.ErrInternalServer)
			}
			return
		}
	}

	log.Debug("appCookieDeviceHandler; Unsupported path:", r.URL.Path, "method:", r.Method)
	httputil.WriteHTTPError(w, api.ErrPathNotFound)
	return
}
