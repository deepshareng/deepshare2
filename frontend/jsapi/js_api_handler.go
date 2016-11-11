package jsapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/appinfo"
	"github.com/MISingularity/deepshare2/deepshared/shorturl"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/frontend/sharelink"
	"github.com/MISingularity/deepshare2/frontend/urlgenerator"
	"github.com/MISingularity/deepshare2/pkg"
	"github.com/MISingularity/deepshare2/pkg/cookieutil"
	"github.com/MISingularity/deepshare2/pkg/httputil"
	in "github.com/MISingularity/deepshare2/pkg/instrumentation"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/MISingularity/deepshare2/pkg/storage"
)

type jsApiHandler struct {
	endpoint         string
	mProducer        messaging.Producer
	client           *http.Client
	ug               *urlgenerator.UrlGenerator
	sl               *sharelink.Sharelink
	specificTokenUrl string
}

func newJsApiHandler(s shorturl.UrlShortener, c *http.Client, matchUrl, cookieUrl, appcookieUrl, appInfoUrl, tokenUrl string, urlBase *url.URL, endpoint string, mp messaging.Producer) http.Handler {
	jsApiHandler := &jsApiHandler{
		endpoint:  endpoint,
		mProducer: mp,
		client:    c,
		ug:        urlgenerator.NewUrlGenerator(s, urlBase),
		sl:        sharelink.NewSharelink(s, matchUrl, cookieUrl, appcookieUrl, appInfoUrl, tokenUrl, urlBase),
	}
	return jsApiHandler
}

func NewJsApiTestHandler(matchUrl, cookieUrl, appcookieUrl, appInfoUrl, tokenUrl string, urlBase *url.URL, endpoint string) http.Handler {
	s := shorturl.NewUrlShortener(http.DefaultClient, storage.NewInMemSimpleKV(), tokenUrl)
	handler := &jsApiHandler{
		sl:        sharelink.NewSharelink(s, matchUrl, cookieUrl, appcookieUrl, appInfoUrl, tokenUrl, urlBase),
		ug:        urlgenerator.NewUrlGenerator(s, urlBase),
		endpoint:  endpoint,
		mProducer: messaging.NewSimpleProducer(bytes.NewBuffer(nil)),
		client:    http.DefaultClient,
	}
	return handler
}

func CreateHandler(endpoint string, urlGeneratorBase string, db storage.SimpleKV, matchUrl, cookieUrl, appcookieUrl, appInfoUrl, tokenUrl string, mp messaging.Producer) http.Handler {
	client := httputil.GetNewClient()
	s := shorturl.NewUrlShortener(client, db, tokenUrl)
	urlBase, err := url.Parse(urlGeneratorBase)
	if err != nil {
		log.Fatalf("[Error]. Share link front; url Base %s is passed by docker, should not in wrong format: %v", urlGeneratorBase, err)
	}
	return newJsApiHandler(s, client, matchUrl, cookieUrl, appcookieUrl, appInfoUrl, tokenUrl, urlBase, endpoint, mp)
}

func (jsl *jsApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !httputil.AllowMethod(w, r.Method, "POST") {
		return
	}

	start := time.Now()
	switch r.Method {
	case "POST":
		defer in.PrometheusForJsApi.HTTPPostDuration(start)

		values := r.URL.Query()
		clicked := false
		contexts, ok := values[clickedTag]
		if ok && contexts[0] == "true" {
			clicked = true
		}
		ip := httputil.ParseClientIP(r)
		uaStr := r.UserAgent()
		appID := path.Base(r.URL.Path)
		log.Debugf("jsApiHandler; Request UA = %s, clicked = %v, appid = %s", uaStr, clicked, appID)
		ua := uainfo.ExtractUAInfoFromUAString(ip, uaStr)
		if !(ua.IsAndroid() || ua.IsIos()) {
			log.Debug("jsApiHandler access with desktop browser")
			return
		}

		if clicked {
			jsl.handlerSecondCall(w, r, appID, ua)
		} else {
			jsl.handlerFirstCall(w, r, appID, ua)
		}

	}
}

// second call, bind (cookieID_deeplinkID) -> inappData
// multiple deeplinkIDs, multiple inappDatas
func (jsl *jsApiHandler) handlerFirstCall(w http.ResponseWriter, r *http.Request, appID string, ua *uainfo.UAInfo) {
	decoder := json.NewDecoder(r.Body)
	reqs := JsApiPost{}
	if err := decoder.Decode(&reqs); err != nil {
		log.Errorf("js api post body decode error: %v, r.ContentLength: %v, r.Body: %#v, r", err, r.ContentLength, r.Body)
		httputil.WriteHTTPError(w, api.ErrBadJSONBody)
		return
	}
	log.Debugf("js api request: %#v\n", reqs)

	cookie, isNewCookie, err := cookieutil.GetCookie(r, jsl.client, jsl.sl.SpecificTokenUrl)
	if err != nil {
		httputil.WriteHTTPError(w, api.ErrInternalServer)
		//The error means system clock is moving backwards, which means the system works abnormal
		//So we should kill this image
		log.Error("jsApiHandler generate cookie", err)
		panic(err)
	}
	log.Debugf("jsApiHandler; cookie value = %s", cookie.Value)

	// for CORS request, need the following headers to set cookie
	o := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", o)
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, If-Modified-Since")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	http.SetCookie(w, cookie)
	log.Debugf("Set header and cookie: %#v\n", w.Header())

	//Call appInfo service
	resp := JsApiResp{}
	appInfo, err := appinfo.GetAppInfoByUrl(jsl.client, appID, jsl.sl.SpecificAppInfoUrl)
	if err != nil {
		httputil.WriteHTTPError(w, api.ErrInternalServer)
		return
	} else if appInfo == nil {
		httputil.WriteHTTPError(w, api.ErrAppIDNotFound)
		return
	} else {
		appInstallStatus := sharelink.AppInstallUnClear
		switch {
		//in wechat browser
		case ua.IsWechat == true:
			if isNewCookie == false {
				appInstallStatus = jsl.sl.GetAppInsStatusWechat(jsl.client, appID, cookie.Value)
			}
		//in native browser
		default:
			if isNewCookie {
				//isNewCookie==true could be
				//	1. first time clicking a deepshare url, or
				//	2. browser cookie cleared
				//so we can not judge if the app is installed
				// but we'd like user jumping to appstore, to save the step to make a choice on landing page
				// so we make it AppNotInstalled
				appInstallStatus = sharelink.AppNotInstalled
			} else {
				appInstallStatus = jsl.sl.GetAppInsStatusBrowser(jsl.client, appID, cookie.Value)
				log.Debugf("jsApiHandler; app install status = %d", appInstallStatus)
			}
		}

		resp.AppID = appID
		resp.ChromeMajor = ua.ChromeMajor
		resp.IsAndroid = (ua.Os == "android")
		resp.IsIos = (ua.Os == "ios")
		resp.IosMajor = ua.IosMajorVersion()
		resp.IsWechat = ua.IsWechat
		resp.IsWeibo = ua.IsWeibo
		resp.IsQQ = ua.IsQQ
		resp.IsFacebook = ua.IsFacebook
		resp.IsTwitter = ua.IsTwitter
		resp.IsFirefox = ua.IsFirefox
		resp.IsQQBrowser = ua.IsQQBrowser
		resp.IsUC = strings.Contains(ua.Browser, "UC Browser")
		resp.CannotDeeplink = ua.CannotDeeplink
		resp.CannotGetWinEvent = ua.CannotGetWindowsEvent
		resp.CannotGoMarket = ua.CannotGoMarket
		resp.ForceUseScheme = ua.ForceUseScheme

		//From App info
		resp.AppName = appInfo.AppName
		resp.IconUrl = appInfo.IconUrl

		//Init below fields in case of the value is changed to "NO_VALUE" when Execute the html template
		resp.Scheme = ""
		resp.Host = ""
		resp.BundleID = ""
		resp.Pkg = ""
		resp.Url = ""
		resp.IsDownloadDirectly = false
		resp.IsUniversallink = false
		resp.YYBUrl = ""
		resp.IsYYBEnableAndroid = false
		resp.IsYYBEnableIosAbove9 = false
		resp.IsYYBEnableIosBelow9 = false

		if ua.IsAndroid() {
			resp.Scheme = appInfo.Android.Scheme   /*"deepshare"*/
			resp.Host = appInfo.Android.Host       /*"com.singulariti.deepsharedemo"*/
			resp.Pkg = appInfo.Android.Pkg         /*"com.singulariti.deepsharedemo"*/
			resp.Url = appInfo.Android.DownloadUrl /*"https://play.google.com/store/apps/details?id=com.xianguo.tingguo"*/
			resp.IsDownloadDirectly = appInfo.Android.IsDownloadDirectly
			resp.IsYYBEnableAndroid = appInfo.Android.YYBEnable
		} else if ua.IsIos() {
			resp.BundleID = appInfo.Ios.BundleID
			resp.Scheme = appInfo.Ios.Scheme   /*"ds7713337217A6E150"*/
			resp.Url = appInfo.Ios.DownloadUrl /*"itms-apps://itunes.apple.com/artist/seligman-ventures-ltd/id515901779"*/
			resp.IsUniversallink = appInfo.Ios.UniversalLinkEnable
			resp.IsYYBEnableIosBelow9 = appInfo.Ios.YYBEnableBelow9
			resp.IsYYBEnableIosAbove9 = appInfo.Ios.YYBEnableAbove9
		}
		if appInfo.YYBEnable {
			resp.YYBUrl = appInfo.YYBUrl
		}

		//From sender info
		resp.MatchId = cookie.Value
		resp.Timestamp = int64(time.Now().UTC().Unix())
		//			resp.DsTag = dsTag//????????
		resp.AppInsStatus = int(appInstallStatus)

	}

	for _, req := range reqs {
		if req.InAppDataReq == nil {
			req.InAppData = ""
		} else {
			if b, err := json.Marshal(req.InAppDataReq); err != nil {
				log.Error("js api; marshal inappdata failed, err:", err)
				httputil.WriteHTTPError(w, api.ErrInternalServer)
				return
			} else {
				req.InAppData = string(b)
			}
		}
		e := &messaging.Event{}
		e = sharelink.AddEventInfo(e, appID, jsl.endpoint, strings.Join(req.Channels, "|"), req.SenderID, cookie.Value, ua.Ip, ua.Ua, "")
		e = sharelink.AddEventOfAppInfo(e, appInfo.YYBEnable, appInfo.Ios.UniversalLinkEnable)
		if e.KVs == nil {
			e.KVs = make(map[string]interface{})
		}
		e.KVs["deeplink_id"] = req.DeepLinkID

		// generate shorturl for:
		// 1. for ios9 devices, serves as universal link
		// 2. for other devices, if callbacks are not set, need redirect to the shorturl to show ds tips or pages
		genUrlBody := &urlgenerator.GenURLPostBody{
			InAppData:          req.InAppData,
			DownloadTitle:      req.DownloadTitle,
			DownloadBtnText:    req.DownloadBtnText,
			DownloadMsg:        req.DownloadMsg,
			DownloadUrlIos:     req.DownloadUrlIos,
			DownloadUrlAndroid: req.DownloadUrlAndroid,
			IsShort:            true,
			SenderID:           req.SenderID,
			Channels:           req.Channels,
		}
		u, err := jsl.ug.GenerateUrl(appID, genUrlBody, ua, jsl.endpoint, jsl.mProducer)
		if err != nil {
			log.Errorf("Failed to generate shorturl, appid: %s, err: %v\n", appID, err)
			httputil.WriteHTTPError(w, api.ErrInternalServer)
		} else {
			if resp.DSUrls == nil {
				resp.DSUrls = make(map[string]string)
			}
			//for ios9, wechat browser, append wcookie in the path (will be sent to server with POST InappData request when app opened via universal link)
			// http://fds.so/appid/shortseg?wcookie=www
			if ua.IsWechat && ua.IosMajorVersion() >= 9 {
				u = jsl.ug.AppendWCookie(w, r, jsl.client, jsl.specificTokenUrl, u)
			}

			resp.DSUrls[req.DeepLinkID] = u
			i := strings.LastIndex(u, "/")
			if i > 0 {
				shortSeg := u[i:]
				e = sharelink.AddEventShortUrlToken(e, shortSeg)
			}
		}

		//Call match service, first call, only bind cookie->inappData, don't bind UA->inappData (will bind in second call)
		go jsl.sl.RequestMatch(jsl.client, appID, req.InAppData, req.SenderID, pkg.EncodeStringSlice(req.Channels), "", cookie.Value, req.DeepLinkID, ua.Ip, "")
		go sharelink.FireShareLinkEvent(e, jsl.mProducer)
	}

	en := json.NewEncoder(w)
	if err := en.Encode(resp); err != nil {
		log.Error("Failed to encode response for JS API, err:", err)
		httputil.WriteHTTPError(w, api.ErrInternalServer)
	}
}

// second call, bind UA -> inappData
func (jsl *jsApiHandler) handlerSecondCall(w http.ResponseWriter, r *http.Request, appID string, ua *uainfo.UAInfo) {
	decoder := json.NewDecoder(r.Body)
	req := JsApiPostClicked{}
	if err := decoder.Decode(&req); err != nil {
		log.Errorf("js api(clicked) post body decode error: %v, r.ContentLength: %v, r.Body: %#v, r", err, r.ContentLength, r.Body)
		httputil.WriteHTTPError(w, api.ErrBadJSONBody)
		return
	}
	if b, err := json.Marshal(req.InAppDataReq); err != nil {
		log.Error("js api(clicked); marshal inappdata failed, err:", err)
		httputil.WriteHTTPError(w, api.ErrInternalServer)
		return
	} else {
		req.InAppData = string(b)
	}
	log.Debugf("js api(clicked) request: %#v\n", req)

	resp := JsApiRespClicked{
		OK: true,
	}
	en := json.NewEncoder(w)
	if err := en.Encode(resp); err != nil {
		log.Error("Failed to encode response for JS API, err:", err)
		httputil.WriteHTTPError(w, api.ErrInternalServer)
	}
	//since match fails when cookie is empty, here we just provide a placeholder to avoid cookie to be empty
	go jsl.sl.RequestMatch(jsl.client, appID, req.InAppData, req.SenderID, pkg.EncodeStringSlice(req.Channels), "", api.CookiePlaceHolder, "", ua.Ip, ua.Ua)
}
