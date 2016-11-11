package sharelink

import (
	"encoding/json"
	"net/http"
	"net/url"

	"golang.org/x/net/context"

	"github.com/MISingularity/deepshare2/deepshared/devicecookier"
	"github.com/MISingularity/deepshare2/deepshared/shorturl"
	"github.com/MISingularity/deepshare2/pkg/httputil"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
)

type Sharelink struct {
	urlShortener         shorturl.UrlShortener
	specificMatchUrl     string
	specificCookieUrl    string
	SpecificAppCookieUrl string
	SpecificAppInfoUrl   string
	SpecificTokenUrl     string
	urlGeneratorBase     *url.URL
}

func NewSharelink(s shorturl.UrlShortener, matchUrl, cookieUrl, appcookieUrl, appInfoUrl, tokenUrl string, urlBase *url.URL) *Sharelink {
	return &Sharelink{
		urlShortener:         s,
		specificMatchUrl:     matchUrl,
		specificCookieUrl:    cookieUrl,
		SpecificAppCookieUrl: appcookieUrl,
		SpecificAppInfoUrl:   appInfoUrl,
		SpecificTokenUrl:     tokenUrl,
		urlGeneratorBase:     urlBase,
	}
}

//parse request to extract info we need
func (slh *Sharelink) parseRequest(ctx context.Context, reqUrl *url.URL) (appID, inAppData, senderID, channels, sdkInfo string, dlInfo DownloadInfo, redirectUrl string) {
	appID = extractAppIDFromRawUrl(reqUrl)
	values := reqUrl.Query()

	contexts, ok := values[InAppData]
	if ok {
		inAppData = contexts[0]
	} else {
		log.Debug("Sharelink does not contain inappData")
	}

	senderIDs, ok := values[SenderID]
	if ok {
		senderID = senderIDs[0]
	} else {
		log.Debug("Sharelink does not contain senderID")
	}

	channelss, ok := values[Channels]
	if ok {
		channels = channelss[0]
	} else {
		log.Debug("Sharelink does not contain channels")
	}

	sdkInfos, ok := values[SDKInfo]
	if ok {
		sdkInfo = sdkInfos[0]
	} else {
		log.Debug("Sharelink does not contain sdkinfo")
	}

	download_titles, ok := values[DownloadTitle]
	if ok {
		dlInfo.DownloadTitle = download_titles[0]
	}
	download_btn_texts, ok := values[DownloadBtnText]
	if ok {
		dlInfo.DownloadBtnText = download_btn_texts[0]
	}
	download_msgs, ok := values[DownloadMsg]
	if ok {
		dlInfo.DownloadMsg = download_msgs[0]
	}
	ios_download_urls, ok := values[DownloadUrlIos]
	if ok {
		dlInfo.DownloadUrlIos = ios_download_urls[0]
	}
	android_download_urls, ok := values[DownloadUrlAndroid]
	if ok {
		dlInfo.DownloadUrlAndroid = android_download_urls[0]
	}
	uninstall_urls, ok := values[UninstallUrl]
	if ok {
		dlInfo.UninstallUrl = uninstall_urls[0]
	}
	force_downloads, ok := values[ForceDownload]
	if ok {
		dlInfo.ForceDownload = force_downloads[0]
	}
	redirect_urls, ok := values[RedirectUrl]
	if ok {
		redirectUrl = redirect_urls[0]
	}

	return
}

func (slh *Sharelink) restoreRawUrl(ctx context.Context, reqUrl *url.URL, e *messaging.Event) (*url.URL, *messaging.Event) {
	log.Debugf("Sharelink Request URL: %s", reqUrl.String())
	appID, shortSeg := extractInfoFromShortUrl(reqUrl)
	e = AddEventShortUrlToken(e, shortSeg)
	rawUrl, err := slh.urlShortener.ToRawURL(reqUrl, appID)
	if err == shorturl.ErrShortSegNotFound {
		e = AddEventShortSegValid(e, false)
	} else if err != nil {
		//the err means Something wrong with DB
		panic(err)
	}
	log.Debugf("Sharelink extract from short url, Raw URL: %s", rawUrl)
	return rawUrl, e

}

func (sl *Sharelink) getUniqueIDByAppCookie(cli *http.Client, appID, cookie string) (uniqueID string, err error) {
	appCookieUrlStr, err := httputil.AppendPath(sl.SpecificAppCookieUrl, appID)
	if err != nil {
		log.Errorf("shareLinkHandler; cookieUrl %s is constructed by us, should not in wrong format: %v", sl.SpecificAppCookieUrl, err)
		return "", err
	} else {
		appCookieUrlStr, err = httputil.AppendPath(appCookieUrlStr, cookie)
		if err != nil {
			log.Errorf("shareLinkHandler; cookieUrl %s is constructed by us, should not in wrong format: %v", appCookieUrlStr, err)
			return "", err
		} else {
			appCookieUrl, err := url.Parse(appCookieUrlStr)
			if err != nil {
				log.Errorf("shareLinkHandler; cookieUrl %s is constructed by us, should not in wrong format: %v", appCookieUrlStr, err)
				panic(err)
			}
			requestAppCookieUrlStr := appCookieUrl.String()
			log.Debugf("shareLinkHandler; Request Cookie URL string is %s", requestAppCookieUrlStr)
			req, err := http.NewRequest("GET", requestAppCookieUrlStr, nil)
			if err != nil {
				log.Error("shareLinkHandler; Setup New Request to Device cookie service failed", err)
				return "", err
			}
			resp, err := cli.Do(req)
			if err != nil {
				log.Error("shareLinkHandler; Request to App cookie Device service failed", err)
				return "", err
			}
			defer resp.Body.Close()
			decoder := json.NewDecoder(resp.Body)
			acdrb := AppCookieDeviceResponseBody{}
			err = decoder.Decode(&acdrb)
			if err != nil {
				log.Error("shareLinkHandler; Decode response from App cookie Device service failed", err)
				return "", err
			}
			uniqueID = acdrb.UniqueID
			log.Debugf("shareLinkHandler; App Cookie to Unique ID Got:%s", uniqueID)
			return uniqueID, nil
		}
	}
}

func (sl *Sharelink) GetAppInsStatusBrowser(cli *http.Client, appID, cookie string) appInsStatus {
	uniqueID, err := sl.getUniqueIDByAppCookie(cli, appID, cookie)
	if err != nil {
		log.Error("failed to getUniqueIDByAppCookie, err:, err")
		return AppInstallUnClear
	}
	if uniqueID != "" {
		return AppInstalled
	} else {
		uid, err := sl.GetUniqueIDByCookie(context.TODO(), cli, cookie)
		if err != nil {
			log.Error("failed to GetUniqueIDByWCookie(wechat), err:", err)
			return AppInstallUnClear
		} else if uid != "" {
			log.Debug("appInstallStatus = AppNotInstalled")
			return AppNotInstalled
		} else {
			return AppInstallUnClear
		}
	}
}

func (sl *Sharelink) GetAppInsStatusWechat(cli *http.Client, appID, cookie string) (appInstallStatus appInsStatus) {
	uid, err := sl.getUniqueIDByAppCookie(cli, appID, cookie)
	if err != nil {
		log.Error("failed to getUniqueIDByAppCookie(wechat), err:", err)
		return AppInstallUnClear
	}
	if uid != "" {
		return AppInstalled
	} else {
		uid, err := sl.GetUniqueIDByWCookie(context.TODO(), cli, cookie)
		if err != nil {
			log.Error("failed to GetUniqueIDByWCookie(wechat), err:", err)
			return AppInstallUnClear
		} else if uid != "" {
			log.Debug("appInstallStatus = AppNotInstalled")
			return AppNotInstalled
		} else {
			return AppInstallUnClear
		}
	}
}

func (sl *Sharelink) GetUniqueIDByCookie(ctx context.Context, cli *http.Client, cookie string) (string, error) {
	return devicecookier.GetUniqueID(ctx, cli, sl.specificCookieUrl, devicecookier.CookieIDPrefix+cookie)
}

func (sl *Sharelink) GetUniqueIDByWCookie(ctx context.Context, cli *http.Client, wCookie string) (string, error) {
	return devicecookier.GetUniqueID(ctx, cli, sl.specificCookieUrl, devicecookier.WCookieIDPrefix+wCookie)
}
