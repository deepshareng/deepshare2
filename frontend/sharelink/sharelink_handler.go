package sharelink

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"html/template"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/appcookiedevice"
	"github.com/MISingularity/deepshare2/deepshared/appinfo"
	"github.com/MISingularity/deepshare2/deepshared/shorturl"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg/cookieutil"
	"github.com/MISingularity/deepshare2/pkg/httputil"
	in "github.com/MISingularity/deepshare2/pkg/instrumentation"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"golang.org/x/net/context"
)

var (
	curDir string
)

type shareLinkHandler struct {
	sl          *Sharelink
	endPoint    string
	msgProducer messaging.Producer
	client      *http.Client
}

func setupServerEnv() {
	//get absolute path, which will be used to locate html and js files
	_, filename, _, _ := runtime.Caller(0)
	dir, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		log.Error("error when get filepath", err)
		panic(err)
	}
	curDir = dir
}

func newShareLinkHandler(s shorturl.UrlShortener, c *http.Client, matchUrl, cookieUrl, appcookieUrl, appInfoUrl, tokenUrl string, urlBase *url.URL, endpoint string, mp messaging.Producer) http.Handler {
	shareLinkHandler := &shareLinkHandler{
		sl:          NewSharelink(s, matchUrl, cookieUrl, appcookieUrl, appInfoUrl, tokenUrl, urlBase),
		endPoint:    endpoint,
		msgProducer: mp,
		client:      c,
	}
	return shareLinkHandler
}

func NewShareLinkTestHandler(matchUrl, cookieUrl, appcookieUrl, appInfoUrl, tokenUrl string, urlBase *url.URL, endpoint string) http.Handler {
	s := shorturl.NewUrlShortener(http.DefaultClient, storage.NewInMemSimpleKV(), tokenUrl)
	shareLinkHandler := &shareLinkHandler{
		sl:          NewSharelink(s, matchUrl, cookieUrl, appcookieUrl, appInfoUrl, tokenUrl, urlBase),
		endPoint:    endpoint,
		msgProducer: messaging.NewSimpleProducer(bytes.NewBuffer(nil)),
		client:      http.DefaultClient,
	}
	setupServerEnv()
	return shareLinkHandler
}

func AddHandler(mux *http.ServeMux, endpoint string, urlGeneratorBase string, db storage.SimpleKV, matchUrl, cookieUrl, appcookieUrl, appInfoUrl, tokenUrl string, mp messaging.Producer) {
	client := httputil.GetNewClient()
	s := shorturl.NewUrlShortener(client, db, tokenUrl)
	urlBase, err := url.Parse(urlGeneratorBase)
	if err != nil {
		log.Fatalf("[Error]. Share link front; url Base %s is passed by docker, should not in wrong format: %v", urlGeneratorBase, err)
	}
	mux.Handle(endpoint, newShareLinkHandler(s, client, matchUrl, cookieUrl, appcookieUrl, appInfoUrl, tokenUrl, urlBase, endpoint, mp))
	setupServerEnv()
	mux.Handle(api.JSServerPrefix, http.StripPrefix(api.JSServerPrefix, http.FileServer(http.Dir(curDir+jsFileServerDir))))
}

func (slh *shareLinkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !httputil.AllowMethod(w, r.Method, "GET") {
		return
	}

	start := time.Now()
	switch r.Method {
	case "GET":
		defer in.PrometheusForSharelink.HTTPGetDuration(start)
		ctx := context.TODO()
		ip := httputil.ParseClientIP(r)
		ua := r.UserAgent()
		log.Debugf("shareLinkHandler; Request UA = %s", ua)
		uaInfo := uainfo.ExtractUAInfoFromUAString(ip, ua)
		if !(uaInfo.IsAndroid() || uaInfo.IsIos()) {
			//TODO: we need to add sending sms to mobile for url opened in desktop browser
			log.Debug("shareLinkHandler access with desktop browser")
			dstHtml := "sharelink_response_desktop.html"
			wirteDstHtml(w, dstHtml, "")
			return
		}

		requestUrl := new(url.URL)
		requestUrl.Scheme = slh.sl.urlGeneratorBase.Scheme
		requestUrl.Host = slh.sl.urlGeneratorBase.Host
		requestUrl.Path = r.URL.Path
		requestUrl.RawQuery = r.URL.RawQuery
		dsTag := extractDSTag(requestUrl)
		wCookie := extractWechatCookie(requestUrl)

		isShort, isLegal := shorturl.IsLegalShortFormat(requestUrl)
		if !isLegal {
			dstHtml := "sharelink_response_illegal.html"
			wirteDstHtml(w, dstHtml, "")
			return
		}
		e := &messaging.Event{}
		var rawUrl *url.URL
		if isShort {
			appID, shortSeg := extractInfoFromShortUrl(requestUrl)
			if shortSeg == shortSegForInvalid {
				log.Debugf("shareLinkHandler; Visit the invalidate page")
				//				dstHtml := api.JSServerPrefix + "sharelink_response_closeme.html"
				//				http.Redirect(w, r, dstHtml, http.StatusOK)
				dstHtml := "sharelink_response_closeme.html"
				_, err := httputil.AppendPath(slh.sl.SpecificAppInfoUrl, appID)
				if err != nil {
					log.Errorf("shareLinkHandler; appInfoUrl %s is constructed by us, should not in wrong format: %v", slh.sl.SpecificAppInfoUrl, err)
					panic(err)
				}
				appInfo, err := appinfo.GetAppInfoByUrl(slh.client, appID, slh.sl.SpecificAppInfoUrl)
				if err != nil {
					httputil.WriteHTTPError(w, api.ErrInternalServer)
					return
				} else if appInfo == nil {
					//TODO: maybe we should show this to user?
					wirteDstHtml(w, dstHtml, "")
					return
				} else {
					wirteDstHtml(w, dstHtml, appInfo.Theme)
				}
				return
			}
			rawUrl, e = slh.sl.restoreRawUrl(ctx, requestUrl, e)
		} else {
			rawUrl = requestUrl
		}

		appID, inAppData, senderID, channels, sdkInfo, dlInfo, redirectUrl := slh.sl.parseRequest(ctx, rawUrl)
		cookie, isNewCookie, err := cookieutil.GetCookie(r, slh.client, slh.sl.SpecificTokenUrl)
		if err != nil {
			httputil.WriteHTTPError(w, api.ErrInternalServer)
			//The error means system clock is moving backwards, which means the system works abnormal
			//So we should kill this image
			log.Error("shareLinkHandler generate cookie", err)
			panic(err)
		}
		log.Debugf("shareLinkHandler; cookie value = %s", cookie.Value)
		http.SetCookie(w, cookie)

		e = AddEventInfo(e, appID, slh.endPoint, channels, senderID, cookie.Value, ip, ua, dsTag)
		//Call appInfo service
		appInfo, err := appinfo.GetAppInfoByUrl(slh.client, appID, slh.sl.SpecificAppInfoUrl)
		if err != nil {
			httputil.WriteHTTPError(w, api.ErrInternalServer)
			return
		} else if appInfo == nil {
			httputil.WriteHTTPError(w, api.ErrAppIDNotFound)
			return
		} else {
			appInstallStatus := AppInstallUnClear
			switch {
			//in wechat browser
			case uaInfo.IsWechat == true:
				if isNewCookie == false {
					appInstallStatus = slh.sl.GetAppInsStatusWechat(slh.client, appID, cookie.Value)
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
					appInstallStatus = AppNotInstalled
				} else {
					appInstallStatus = slh.sl.GetAppInsStatusBrowser(slh.client, appID, cookie.Value)
					log.Debugf("shareLinkHandler; app install status = %d", appInstallStatus)
				}
			}

			writeResponse(ctx, w, r, appInfo, appID, cookie.Value, uaInfo, dlInfo, redirectUrl, dsTag, int(appInstallStatus))
			e = AddEventOfAppInfo(e, appInfo.YYBEnable, appInfo.Ios.UniversalLinkEnable)
		}

		if uaInfo.IsIos() {
			if appInfo.Ios.UniversalLinkEnable && uaInfo.IosMajorVersion() >= 9 && !appInfo.Ios.YYBEnableAbove9 {
				//Under IOS9 & universallink enable situation, UA should not be used to bind match. Because:
				//1, For install, 100% match should take effect, and cookie is used to match.
				//2. For open, url would return the data directly.
				// but when appInfo.Ios.YYBEnableAbove9 is true, UA is needed to match
				//And if we bind will match in this case, app would receive repeated data when it is opened normally.
				ua = ""
			}
		}
		//Call match service
		go slh.sl.RequestMatch(slh.client, appID, inAppData, senderID, channels, sdkInfo, cookie.Value, "", ip, ua)

		// open in safari with wCookie (previously opened with wechat)
		if !uaInfo.IsWechat && wCookie != "" {
			go slh.bindCookieWithWCookie(appID, wCookie, cookie.Value)
		}
		go FireShareLinkEvent(e, slh.msgProducer)
	}
}

func extractDSTag(u *url.URL) (dsTag string) {
	values := u.Query()

	contexts, ok := values[dsTAG]
	if ok {
		dsTag = contexts[0]
		log.Debug("Sharelink; dstag:" + dsTag)
	} else {
		log.Debug("Sharelink does not contain dstag")
		dsTag = ""
	}
	return
}

func extractWechatCookie(u *url.URL) (wCookie string) {
	values := u.Query()

	contexts, ok := values[api.WechatCookieQueryKey]
	if ok {
		wCookie = contexts[0]
		log.Debug("Sharelink; wechat cookie:" + wCookie)
	} else {
		log.Debug("Sharelink does not contain wechat cookie")
		wCookie = ""
	}
	return
}

func extractInfoFromShortUrl(shortUrl *url.URL) (string, string) {
	pathes := strings.Split(shortUrl.Path, "/")
	appID := pathes[2]
	shortSeg := pathes[3]
	return appID, shortSeg
}

func extractAppIDFromRawUrl(rawurl *url.URL) string {
	pathes := strings.Split(rawurl.Path, "/")
	appID := pathes[2]
	return appID
}

func wirteDstHtml(w http.ResponseWriter, dstHtml string, theme string) {
	srcHtmlFile := curDir + jsTemplateDir + dstHtml
	t, err := template.ParseFiles(srcHtmlFile)
	if err != nil {
		log.Error("Sharelink error when parse html files", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	info := make(map[string]interface{})

	if theme == "" {
		info["Theme"] = "0"
	} else {
		info["Theme"] = theme
	}
	err = t.Execute(w, info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// pair cookie and wCookie
// if cookie has binded uniqueID, then bind it under wCookie too (in appcookiedevice)
func (slh *shareLinkHandler) bindCookieWithWCookie(appID, wCookie, cookie string) {
	log.Debug("put dscookie! cookie:", cookie, "wechat cookie:", wCookie)

	if cookie == "" || wCookie == "" {
		return
	}
	if slh.sl.SpecificAppCookieUrl == "" {
		return
	}

	b, err := json.Marshal(appcookiedevice.PostPairCookieBody{
		Cookie1: cookie,
		Cookie2: wCookie,
	})
	if err != nil {
		log.Fatal("SharelinkHandler; Marshal device info failed", err)
		return
	}
	appCookieUrl, err := httputil.AppendPath(slh.sl.SpecificAppCookieUrl, appcookiedevice.PairCookiesPath)
	if err != nil {
		panic(err)
	}
	appCookieUrl, err = httputil.AppendPath(appCookieUrl, appID)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", appCookieUrl, bytes.NewReader(b))
	if err != nil {
		log.Fatal("SharelinkHandler; new request failed! err:", err)
	}
	if _, err := slh.client.Do(req); err != nil {
		if err != nil {
			log.Error("SharelinkHandler; Request to CookiePair service failed", err)
		}
	}
}
