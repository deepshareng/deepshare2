package inappdata

import (
	"bytes"
	"net/http"

	"encoding/json"
	"time"

	"path"

	"net/url"

	"io/ioutil"

	"strings"

	"strconv"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/appcookiedevice"
	"github.com/MISingularity/deepshare2/deepshared/appinfo"
	"github.com/MISingularity/deepshare2/deepshared/devicecookier"
	"github.com/MISingularity/deepshare2/pkg/httputil"
	in "github.com/MISingularity/deepshare2/pkg/instrumentation"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"golang.org/x/net/context"
)

type inAppDataHandler struct {
	client                  *http.Client
	specificMatchUrl        string
	specificCookieUrl       string
	specificAppCookieUrl    string
	specificUrlGeneratorUrl string
	specificAppInfoUrl      string
	mProducer               messaging.Producer
	endpoint                string
}

func AddHandler(mux *http.ServeMux, endpoint string, matchUrl string, cookieUrl string, appcookieUrl, urlgeneratorUrl, appInfoUrl string, mp messaging.Producer) {
	client := httputil.GetNewClient()
	mux.Handle(endpoint, newInAppDataHandler(client, matchUrl, cookieUrl, appcookieUrl, urlgeneratorUrl, appInfoUrl, mp, endpoint))
}

func newInAppDataHandler(client *http.Client, matchUrl, cookieUrl, appcookieUrl, urlgeneratorUrl, appInfoUrl string, mp messaging.Producer, endpoint string) http.Handler {
	inAppDataHandler := &inAppDataHandler{
		client:                  client,
		specificMatchUrl:        matchUrl,
		specificCookieUrl:       cookieUrl,
		specificAppCookieUrl:    appcookieUrl,
		specificUrlGeneratorUrl: urlgeneratorUrl,
		specificAppInfoUrl:      appInfoUrl,
		mProducer:               mp,
		endpoint:                endpoint,
	}
	return inAppDataHandler
}

func NewInAppDataTestHandler(matchUrl, cookieUrl, appcookieUrl, urlgeneratorUrl, appInfoUrl string, pfx string) http.Handler {
	inAppDataHandler := &inAppDataHandler{
		client:                  http.DefaultClient,
		specificMatchUrl:        matchUrl,
		specificCookieUrl:       cookieUrl,
		specificAppCookieUrl:    appcookieUrl,
		specificUrlGeneratorUrl: urlgeneratorUrl,
		specificAppInfoUrl:      appInfoUrl,
		mProducer:               messaging.NewSimpleProducer(bytes.NewBuffer(nil)),
		endpoint:                pfx,
	}
	return inAppDataHandler
}

func (iadh *inAppDataHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !httputil.AllowMethod(w, r.Method, "POST") {
		return
	}

	start := time.Now()
	switch r.Method {
	case "POST":
		defer in.PrometheusForInappData.HTTPPostDuration(start)
		var body []byte
		var err error
		if body, err = ioutil.ReadAll(r.Body); err != nil {
			log.Errorf("InAppDataHandler; Read Request body error: %s", err)
			httputil.WriteHTTPError(w, api.ErrBadRequestBody)
		}
		log.Debugf("InAppDataHandler; Post Body : %s", string(body))
		//Extract app ID
		appID := path.Base(r.URL.Path)
		if appID == "/" {
			log.Errorf("InAppDataHandler; Illegal Request URL: %s", r.URL.String())
			httputil.WriteHTTPError(w, api.ErrPathNotFound)
		}

		iadpb := InAppDataPostBody{}
		if err := json.Unmarshal(body, &iadpb); err != nil {
			log.Errorf("InAppDataHandler; Get In App Data post body decode error: body = %s; err = %v", body, err)
			httputil.WriteHTTPError(w, api.ErrBadJSONBody)
			return
		}

		//pre process, convert "dl_no_value" to "", since SDK side use "dl_no_value" as empty
		if iadpb.ShortSeg == noValueFromSDK {
			iadpb.ShortSeg = ""
		}
		if iadpb.ClickID == noValueFromSDK {
			iadpb.ClickID = ""
		}

		queryPartPrefix := "?"

		// parse wcookie from shortseg, in case SDK did not split wcookie from shortseg
		if iadpb.WCookie == "" && strings.Contains(iadpb.ShortSeg, queryPartPrefix) {
			u, err := url.Parse(iadpb.ShortSeg)
			if err != nil {
				log.Error("Failed to parse shortseg:", iadpb.ShortSeg, "as url, err:", err)
			} else {
				values := u.Query()
				if v, ok := values[api.WechatCookieQueryKey]; ok {
					iadpb.WCookie = v[0]
				}
			}

		}

		//pre process, sometimes clickID and shortSeg may be polluted by adding query strings in specific browser, such as QQ
		if strings.Contains(iadpb.ShortSeg, queryPartPrefix) {
			iadpb.ShortSeg = strings.Split(iadpb.ShortSeg, queryPartPrefix)[0]
		}
		if strings.Contains(iadpb.ClickID, queryPartPrefix) {
			iadpb.ClickID = strings.Split(iadpb.ClickID, queryPartPrefix)[0]
		}

		go func() {
			wcookie, _ := iadh.getWCookieID(context.TODO(), iadpb.UniqueID)
			// wcookie is refreshed by wechat, should inherit all data saved under wcookie to iadpb.WCookie
			if wcookie != "" && wcookie != iadpb.WCookie {
				appcookiedevice.RefreshCookie(context.TODO(), iadh.client, iadh.specificAppCookieUrl, wcookie, iadpb.WCookie)
			}
			// Bind wCookie and uniqueID
			iadh.bindUniqueIDWCookie(context.TODO(), iadpb.UniqueID, iadpb.WCookie)

			// Bind appID, wCookieID to uniqueID when app is opened by wCookie(wechat open app by universal link)
			iadh.bindAppCookieToUniqueID(appID, iadpb.WCookie, iadpb.UniqueID)

			// Bind cookieID and hardwareID
			iadh.bindHardwareIDCookie(context.TODO(), iadpb.ClickID, iadpb.HardwareID)

			// Bind appID, cookieID to uniqueID when app is opened by click ID
			iadh.bindAppCookieToUniqueID(appID, iadpb.ClickID, iadpb.UniqueID)
		}()

		//Call match service
		mrri := api.MatchReceiverInfo{}
		if err := json.Unmarshal(body, &mrri); err != nil {
			log.Errorf("InAppDataHandler; Get In App Data post body decode error: %s", err)
			httputil.WriteHTTPError(w, api.ErrBadJSONBody)
			return
		}
		receiverInfo, err := json.Marshal(mrri)
		if err != nil {
			log.Errorf("InAppDataHandler; Encode request match Data error: %s", err)
			httputil.WriteHTTPError(w, api.ErrBadJSONBody)
			return
		}
		ip := httputil.ParseClientIP(r)
		ua := r.UserAgent()
		log.Debugf("InAppDataHandler; Request UA = %s", ua)
		tracking := ""
		if iadpb.IsNewUser {
			tracking = RequestTrackingValueInstall
		} else {
			tracking = RequestTrackingValueOpen
		}
		inAppData := InAppDataResponseBody{}
		var retErr error
		var respObj api.MatchResponse
		// for ios9 universal link, short_seg is used to extract in-app data
		// use short_seg to get inapp data, no need to request match
		if iadpb.ShortSeg != "" {
			go iadh.bindAppCookieToUniqueIDWhenUniversalLink(appID, iadpb.UniqueID)
			respObj, retErr = iadh.getInAppDataFromUrlGenerator(context.TODO(), appID, iadpb.ShortSeg, string(receiverInfo), tracking)
			if retErr == nil {
				inAppData.InAppData = respObj.InappData
				inAppData.Channels = respObj.Channels
				go produceInAppDataFromRawUrlEvent(appID, iadh.endpoint, iadpb.ShortSeg, tracking, respObj.InappData, respObj.SenderID, ip, ua, respObj.Channels, mrri, iadh.mProducer)
			}
		} else {
			clickId := ""
			deeplinkID := ""
			cookieId := ""
			if iadpb.ClickID == "" {

				//Call cookie service
				if !iadh.shouldUseWechatCookie(appID, iadpb) {
					cookieId, _ = iadh.getCookieID(context.TODO(), iadpb.UniqueID)
					if cookieId != "" {
						clickId = cookieId
					}
				} else if ver, _ := strconv.ParseFloat(iadpb.OSVersion, 64); ver >= 9 {
					cookieId, _ = iadh.getWCookieID(context.TODO(), iadpb.UniqueID)
					if cookieId != "" {
						clickId = cookieId
					}
				}
				// Bind appID, cookieID to uniqueID when app is opened by home screen click or app store open button
				go iadh.bindAppCookieToUniqueID(appID, cookieId, iadpb.UniqueID)

			} else {
				clickId = iadpb.ClickID
				deeplinkID = iadpb.DeeplinkID
			}

			respObj, retErr = iadh.getInAppDataFromMatch(context.TODO(), appID, clickId, deeplinkID, ip, ua, string(receiverInfo), tracking)
			if retErr == nil {
				inAppData.InAppData = respObj.InappData
				inAppData.Channels = respObj.Channels
				go produceInAppDataFromMatchEvent(appID, iadh.endpoint, tracking, respObj.InappData, respObj.SenderID, ip, ua, clickId, cookieId, respObj.Channels, mrri, iadh.mProducer)
			}
		}

		if retErr != nil {
			log.Error("[Error] InAppDataHandler; request Match/urlgenerator service error:", err)
			httputil.WriteHTTPError(w, api.ErrInternalServer)
			go produceInAppDataErrorEvent(appID, iadh.endpoint, tracking, ip, ua, mrri, retErr, iadh.mProducer)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		// set to empty slice instead of nil to avoid crash in SDK
		if inAppData.Channels == nil {
			inAppData.Channels = []string{}
		}
		log.Debugf("InAppDataHandler; inAppData returned = %s, channels = %v\n", inAppData.InAppData, inAppData.Channels)
		if err := encoder.Encode(inAppData); err != nil {
			// TODO: use a logger pkg and change this to debug level
			log.Errorf("api: failed to encode inapp data response to %s", r.RemoteAddr)
			httputil.WriteHTTPError(w, api.ErrInternalServer)
		}
	}
}

func (iadh *inAppDataHandler) bindHardwareIDCookie(ctx context.Context, cookieID, hardwareId string) {
	if hardwareId != "" && cookieID != "" {
		devicecookier.PutCookieInfo(ctx, iadh.client, iadh.specificCookieUrl, devicecookier.HardwareIDPrefix+hardwareId, cookieID)
	}
}

func (iadh *inAppDataHandler) bindUniqueIDWCookie(ctx context.Context, uniqueID, wCookieID string) {
	if uniqueID != "" && wCookieID != "" && wCookieID != noValueFromSDK {
		devicecookier.PutCookieInfo(ctx, iadh.client, iadh.specificCookieUrl, devicecookier.UniqueIDPrefixWCookie+uniqueID, wCookieID)
		devicecookier.PutDeviceInfo(ctx, iadh.client, iadh.specificCookieUrl, devicecookier.WCookieIDPrefix+wCookieID, uniqueID)
	}
}

func (iadh *inAppDataHandler) bindAppCookieToUniqueID(appId, cookieId, uniqueId string) {
	if cookieId == "" || cookieId == noValueFromSDK {
		return
	}
	if uniqueId == "" {
		return
	}
	if iadh.specificAppCookieUrl == "" {
		return
	}

	iadh.putAppCookieDevice(appId, cookieId, uniqueId)
}

func (iadh *inAppDataHandler) putAppCookieDevice(appId, cookieId, uniqueId string) {
	requestAppCookieUrlStr, err := httputil.AppendPath(iadh.specificAppCookieUrl, appId)
	if err != nil {
		log.Errorf("InAppDataHandler; appcookieUrl %s is constructed by us, should not in wrong format: %v", iadh.specificAppCookieUrl, err)
		return
	}
	requestAppCookieUrlStr, err = httputil.AppendPath(requestAppCookieUrlStr, cookieId)
	if err != nil {
		log.Errorf("InAppDataHandler; appcookieUrl %s is constructed by us, should not in wrong format: %v", requestAppCookieUrlStr, err)
		return
	}
	b, err := json.Marshal(appcookiedevice.DeviceInfo{UniqueId: uniqueId})
	if err != nil {
		log.Fatal("InAppDataHandler; Marshal device info failed", err)
		return
	}
	log.Debugf("InAppDataHandler; Request AppCookie URL string is %s", requestAppCookieUrlStr)
	req, err := http.NewRequest("PUT", requestAppCookieUrlStr, bytes.NewReader(b))

	if err != nil {
		log.Error("InAppDataHandler; Setup New Request to App cookie Device service failed", err)
		return
	}
	_, err = iadh.client.Do(req)
	if err != nil {
		log.Error("InAppDataHandler; Request to App cookie Device service failed", err)
	}
}

func (iadh *inAppDataHandler) bindAppCookieToUniqueIDWhenUniversalLink(appId, uniqueId string) {
	if iadh.specificAppCookieUrl == "" {
		return
	}

	cookieId, _ := iadh.getCookieID(context.TODO(), uniqueId)
	// Bind appID, cookieID to uniqueID when app is opened by home screen click or app store open button
	iadh.bindAppCookieToUniqueID(appId, cookieId, uniqueId)

}

func (iadh *inAppDataHandler) getCookieID(ctx context.Context, uniqueID string) (cookieId string, err error) {
	return devicecookier.GetCookieID(ctx, iadh.client, iadh.specificCookieUrl, devicecookier.UniqueIDPrefix+uniqueID)
}

func (iadh *inAppDataHandler) getWCookieID(ctx context.Context, uniqueID string) (cookieId string, err error) {
	return devicecookier.GetCookieID(ctx, iadh.client, iadh.specificCookieUrl, devicecookier.UniqueIDPrefixWCookie+uniqueID)
}

func (iadh *inAppDataHandler) getInAppDataFromMatch(ctx context.Context, appID, cookieID, deeplinkID, ip, ua, receiverInfo, tracking string) (mr api.MatchResponse, err error) {
	var req *http.Request
	matchUrlStr, err := httputil.AppendPath(iadh.specificMatchUrl, appID)
	if err != nil {
		log.Errorf("InAppDataHandler; matchUrl %s is constructed by us, should not in wrong format: %v", iadh.specificMatchUrl, err)
		panic(err)
	}
	matchUrl, err := url.Parse(matchUrlStr)
	if err != nil {
		log.Errorf("InAppDataHandler; matchUrl %s is constructed by us, should not in wrong format: %v", matchUrlStr, err)
		panic(err)
	}

	if cookieID != "" && deeplinkID != "" {
		// /v2/matches/cookieID_deeplinkID
		pathMatch := matchUrl.Path
		newPathMatch := path.Join(pathMatch, cookieID+"_"+deeplinkID)
		matchUrl.Path = newPathMatch
	} else if cookieID != "" {
		// /v2/matches/cookieID
		pathMatch := matchUrl.Path
		newPathMatch := path.Join(pathMatch, cookieID)
		matchUrl.Path = newPathMatch
	}

	queries := matchUrl.Query()
	queries.Add(RequestReceiverInfoKey, receiverInfo)
	queries.Add(MatchRequestClientIPKey, ip)
	queries.Add(MatchRequestClientUAKey, ua)
	queries.Add(RequestTrackingKey, tracking)
	matchUrl.RawQuery = queries.Encode()
	requestMatchUrlStr := matchUrl.String()
	log.Debugf("InAppDataHandler; Request Match URL string is %s", requestMatchUrlStr)
	req, err = http.NewRequest("GET", requestMatchUrlStr, nil)
	if err != nil {
		log.Error("InAppDataHandler; Setup New Request to Match service failed", err)
		return mr, err
	}

	resp, err := iadh.client.Do(req)
	if err != nil {
		log.Error("InAppDataHandler; Request to Match service failed", err)
		return mr, err
	}
	defer resp.Body.Close()
	// We need to make sure the response body of match get should be in the same format
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&mr)
	if err != nil {
		log.Error("InAppDataHandler; Decode response from Match service failed", err)
	}
	return mr, nil
}

func (iadh *inAppDataHandler) getInAppDataFromUrlGenerator(ctx context.Context, appID, shortSeg, receiverInfo, tracking string) (mr api.MatchResponse, err error) {
	var req *http.Request
	urlgeneratorUrlStr, err := httputil.AppendPath(iadh.specificUrlGeneratorUrl, appID)
	if err != nil {
		log.Errorf("InAppDataHandler; urlgeneratorUrl %s is constructed by us, should not in wrong format: %v", iadh.specificUrlGeneratorUrl, err)
		panic(err)
	}
	urlgeneratorUrl, err := url.Parse(urlgeneratorUrlStr)
	if err != nil {
		log.Errorf("InAppDataHandler; urlgeneratorUrl %s is constructed by us, should not in wrong format: %v", urlgeneratorUrlStr, err)
		panic(err)
	}
	queries := urlgeneratorUrl.Query()
	queries.Add(RequestReceiverInfoKey, receiverInfo)
	queries.Add(RequestTrackingKey, tracking)
	urlgeneratorUrl.RawQuery = queries.Encode()
	urlgeneratorUrl.Path = path.Join(urlgeneratorUrl.Path, shortSeg)
	requestUrlGeneratorUrlStr := urlgeneratorUrl.String()

	log.Debugf("InAppDataHandler; Request UrlGenerator URL string is %s", requestUrlGeneratorUrlStr)
	req, err = http.NewRequest("GET", requestUrlGeneratorUrlStr, nil)
	if err != nil {
		log.Error("InAppDataHandler; Setup New Request to urlgenerator service failed", err)
		return mr, err
	}

	resp, err := iadh.client.Do(req)
	if err != nil {
		log.Error("InAppDataHandler; Request to urlgenerator service failed", err)
		return mr, err
	}
	defer resp.Body.Close()
	// We need to make sure the response body of urlgenerator get should be in the same format

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("InAppDataHandler; Readall from url urlgenerator failed, err:", err)
		return mr, err
	}
	if err := json.Unmarshal(b, &mr); err != nil {
		log.Error("InAppDataHandler; Decode response from urlgenerator service failed", err, "data:", string(b))
		return mr, err
	}

	log.Debugf("InAppDataHandler; get inappdata response from urlgenerator succeed, data: %#v\n", mr)
	return mr, nil
}

func (iadh *inAppDataHandler) shouldUseWechatCookie(appId string, iadpb InAppDataPostBody) bool {
	if iadh.specificAppInfoUrl == "" {
		log.Errorf("InAppDataHandler; AppInfoUrl is not assigned")
		panic(nil)
	}
	appInfo, err := appinfo.GetAppInfoByUrl(iadh.client, appId, iadh.specificAppInfoUrl)
	if err != nil {
		log.Errorf("InAppDataHandler; Request to appinfo service failed")
		panic(err)
	}
	if appInfo == nil {
		return false
	}
	if iadpb.OS == "iOS" {
		ver, err := strconv.ParseFloat(iadpb.OSVersion, 64)
		if err != nil {
			log.Errorf("InAppDataHandler; IOs version code is illegal, %s", iadpb.OSVersion)
			//This is because 9 takes majority market
			ver = 9
		}
		if ver >= 9 {
			if appInfo.Ios.YYBEnableAbove9 {
				return true
			}
		} else {
			if appInfo.Ios.YYBEnableBelow9 {
				return true
			}
		}
	} else if iadpb.OS == "Android" {
		if appInfo.Android.YYBEnable {
			return true
		}
		//TODO: delete this in future
		//At first, there is only appInfo.YYBEnable's value available, and then we change to use appInfo.Android.YYBEnable
		//So we need to use the initial value to sync.
		if appInfo.YYBEnable {
			return true
		}
	}

	return false
}
