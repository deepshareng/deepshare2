package urlgenerator

import (
	"net/url"
	"path"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/shorturl"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/frontend/sharelink"
	"github.com/MISingularity/deepshare2/pkg"

	"net/http"

	"time"

	"github.com/MISingularity/deepshare2/pkg/cookieutil"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
)

type UrlGenerator struct {
	urlShortener     shorturl.UrlShortener
	urlGeneratorBase *url.URL
}

func NewUrlGenerator(s shorturl.UrlShortener, urlBase *url.URL) *UrlGenerator {
	return &UrlGenerator{
		urlShortener:     s,
		urlGeneratorBase: urlBase,
	}
}

func (gu *UrlGenerator) urlPrefix(appID string) *url.URL {
	urlResult := new(url.URL)
	urlResult.Scheme = gu.urlGeneratorBase.Scheme
	urlResult.Host = gu.urlGeneratorBase.Host
	urlResult.Path = path.Join(api.ShareLinkPrefix, appID)
	return urlResult
}
func (gu *UrlGenerator) GenerateUrl(appID string, genURLPostBody *GenURLPostBody, ua *uainfo.UAInfo, endpoint string, mp messaging.Producer) (gUrl string, err error) {
	urlResult := gu.urlPrefix(appID)
	values := make(url.Values)
	values.Add(sharelink.DownloadTitle, genURLPostBody.DownloadTitle)
	values.Add(sharelink.DownloadBtnText, genURLPostBody.DownloadBtnText)
	values.Add(sharelink.DownloadMsg, genURLPostBody.DownloadMsg)
	values.Add(sharelink.DownloadUrlIos, genURLPostBody.DownloadUrlIos)
	values.Add(sharelink.DownloadUrlAndroid, genURLPostBody.DownloadUrlAndroid)
	values.Add(sharelink.UninstallUrl, genURLPostBody.UninstallUrl)
	values.Add(sharelink.RedirectUrl, genURLPostBody.RedirectUrl)
	values.Add(sharelink.InAppData, genURLPostBody.InAppData)
	values.Add(sharelink.SDKInfo, genURLPostBody.SDKInfo)
	values.Add(sharelink.SenderID, genURLPostBody.SenderID)
	channelsStr := pkg.EncodeStringSlice(genURLPostBody.Channels)
	values.Add(sharelink.Channels, channelsStr)
	values.Add(sharelink.ForceDownload, genURLPostBody.ForceDownload)
	urlResult.RawQuery = values.Encode()
	rawUrlStr, shortUrlStr := urlResult.String(), ""
	urlResultStr := rawUrlStr
	log.Debugf("Url generator: Raw url is %s", rawUrlStr)
	if genURLPostBody.IsShort {
		shortURLLifeTime := shorturl.ShortURLLifeTimeDefault
		if CanSetShortURLLifeTime(appID) {
			shortURLLifeTime = time.Duration(24*30) * time.Hour
		}
		shortUrl, err := gu.urlShortener.ToShortURL(urlResult, appID, genURLPostBody.IsPermanent, shortURLLifeTime, genURLPostBody.UseShortID)
		if err != nil {
			log.Error("Url generator failed, err:", err)
		}
		shortUrlStr = shortUrl.String()
		urlResultStr = shortUrlStr
		log.Debugf("Url generator: short url is %s", shortUrlStr)
	}
	if err := produceGenerateUrlEvent(appID, rawUrlStr, shortUrlStr, endpoint, genURLPostBody, ua, mp); err != nil {
		//TODO should write an error log for alerting
		panic(err)
	}
	return urlResultStr, nil
}

//TODO need add field in appinfo to define the LifeTime limit, need to be binded with payment status in dashboard
// temporarily only set lifetime for the specific appID
func CanSetShortURLLifeTime(appID string) bool {
	if appID == "3ed77201a3187b34" {
		return true
	}
	return false
}

func (gu UrlGenerator) AppendWCookie(w http.ResponseWriter, r *http.Request, client *http.Client, specificTokenUrl string, u string) string {
	cookie, _, err := cookieutil.GetCookie(r, client, specificTokenUrl)
	if err != nil {
		log.Error("AppendWCookie: failed to GetCookie, err:", err)
	} else {
		uu, err := url.Parse(u)
		if err != nil {
			log.Error("AppendWCookie: failed to parse url:", u, "err:", err)
		} else {
			values := make(url.Values)
			values.Add(api.WechatCookieQueryKey, cookie.Value)
			uu.RawQuery = values.Encode()
			u = uu.String()
			http.SetCookie(w, cookie)
		}
	}
	return u
}
