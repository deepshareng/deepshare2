package binddevicetocookie

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/devicecookier"
	"github.com/MISingularity/deepshare2/pkg/cookieutil"
	"github.com/MISingularity/deepshare2/pkg/httputil"
	in "github.com/MISingularity/deepshare2/pkg/instrumentation"
	"github.com/MISingularity/deepshare2/pkg/log"
	"golang.org/x/net/context"
)

type bindDeviceToCookieHandler struct {
	client            *http.Client
	specificCookieUrl string
	specificTokenUrl  string
	endpoint          string
}

func AddHandler(mux *http.ServeMux, endpoint string, cookieUrl, tokenUrl string) {
	client := httputil.GetNewClient()
	mux.Handle(endpoint, newBindDeviceToCookieHandler(client, cookieUrl, tokenUrl, endpoint))
}

func newBindDeviceToCookieHandler(client *http.Client, cookieUrl, tokenUrl, endpoint string) http.Handler {
	inAppDataHandler := &bindDeviceToCookieHandler{
		client:            client,
		specificCookieUrl: cookieUrl,
		specificTokenUrl:  tokenUrl,
		endpoint:          endpoint,
	}
	return inAppDataHandler
}

func (bdtc *bindDeviceToCookieHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !httputil.AllowMethod(w, r.Method, "GET") {
		return
	}
	start := time.Now()
	switch r.Method {
	case "GET":
		defer in.PrometheusForBindDeviceToCookier.HTTPGetDuration(start)

		uniqueId := r.URL.Path[len(bdtc.endpoint):]

		if uniqueId == "" {
			log.Error("BindDeviceToCookieHandler does not contain UniqueID")
			return
		}

		//get cookie given uniqueID
		if r.FormValue("getcookie") == "true" {
			log.Debug("bindDeviceToCookieHandler; get cookieID, uniqueID =", uniqueId)

			cookieID, err := bdtc.getCookieID(context.TODO(), uniqueId)
			if err != nil {
				httputil.WriteHTTPError(w, api.ErrInternalServer)
				panic(err)
			}
			cookieInfo := devicecookier.CookieInfo{CookieID: cookieID}
			en := json.NewEncoder(w)
			if err := en.Encode(cookieInfo); err != nil {
				httputil.WriteHTTPError(w, api.ErrInternalServer)
				panic(err)
			}
			return
		}

		lines, _ := r.Header["Cookie"]
		log.Debugf("BindDeviceToCookieHandler; Request Cookie string is %s", lines)
		cookie, isNew, err := cookieutil.GetCookie(r, bdtc.client, bdtc.specificTokenUrl)
		if isNew {
			log.Debugf("BindDeviceToCookieHandler; New Cookie created")
			http.SetCookie(w, cookie)
		}
		if err != nil {
			httputil.WriteHTTPError(w, api.ErrInternalServer)
			//The error means system clock is moving backwards, which means the system works abnormal
			//So we should kill this image
			log.Error("BindDeviceToCookieHandler generate cookie", err)
			panic(err)
		}

		bdtc.bindDeviceToCookie(context.TODO(), uniqueId, cookie.Value)
	}
}

func (bdtch *bindDeviceToCookieHandler) bindDeviceToCookie(ctx context.Context, uniqueID, cookieID string) {
	if uniqueID != "" && cookieID != "" {
		devicecookier.PutCookieInfo(ctx, bdtch.client, bdtch.specificCookieUrl, devicecookier.UniqueIDPrefix+uniqueID, cookieID)
		devicecookier.PutDeviceInfo(ctx, bdtch.client, bdtch.specificCookieUrl, devicecookier.CookieIDPrefix+cookieID, uniqueID)
	}
}

func (bdtch *bindDeviceToCookieHandler) getCookieID(ctx context.Context, uniqueID string) (cookieID string, err error) {
	return devicecookier.GetCookieID(ctx, bdtch.client, bdtch.specificCookieUrl, devicecookier.UniqueIDPrefix+uniqueID)
}
