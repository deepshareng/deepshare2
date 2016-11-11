package appinsstatus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"strings"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/appcookiedevice"
	"github.com/MISingularity/deepshare2/deepshared/appinfo"
	"github.com/MISingularity/deepshare2/deepshared/devicecookier"
	"github.com/MISingularity/deepshare2/deepshared/match"
	"github.com/MISingularity/deepshare2/deepshared/token"
	"github.com/MISingularity/deepshare2/frontend/inappdata"
	"github.com/MISingularity/deepshare2/frontend/sharelink"
	"github.com/MISingularity/deepshare2/frontend/urlgenerator"
	"github.com/MISingularity/deepshare2/pkg/httputil"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/testutil"
)

const (
	appID1 = "1652E90881C1FAE8"
	appID2 = "9aeea8a7f4c56ff4"
)

func TestAppInsStatusIos9Wechat(t *testing.T) {
	fakeUA := `Mozilla/5.0 (iPhone; CPU iPhone OS 9_0 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Mobile/13A344 MicroMessenger/6.3.1 NetType/WIFI Language/zh_CN`
	// register apps
	appinfoServer := httptest.NewServer(appinfo.NewAppInfoTestHandler(api.AppInfoPrefix))
	RegisterApps(appinfoServer.URL)

	//0. generate shorturl
	tokenServer := httptest.NewServer(token.NewTokenTestHandler(api.TokenPrefix))
	tokenUrl := tokenServer.URL + api.TokenPrefix
	genurlServer := httptest.NewServer(urlgenerator.NewGenerateUrlTestHandler(tokenUrl))
	genurlUrl := genurlServer.URL + api.GenerateUrlPrefix
	genurlReq := &urlgenerator.GenURLPostBody{
		InAppDataReq: "test_inapp_data",
		IsShort:      true,
	}
	genurlResp := &urlgenerator.GenURLResponseBody{}
	latency, err := testutil.PostJson(genurlUrl+appID1, genurlReq, fakeUA, "", genurlResp)
	if err != nil {
		t.Fatal("Failed to generate url, err:", err)
	}
	log.Debug("genurl latency:", latency, "genurl resp:", genurlResp)

	//1. click shorturl of app1 in wechat, should be AppInstallUnClear
	matchServer := httptest.NewServer(match.NewMatchTestHandler(api.MatchPrefix))
	matchUrl := matchServer.URL + api.MatchPrefix
	cookieServer := httptest.NewServer(devicecookier.NewDeviceCookieTestHandler(api.DeviceCookiePrefix))
	cookieUrl := cookieServer.URL + api.DeviceCookiePrefix
	appcookiedeviceServer := httptest.NewServer(appcookiedevice.NewAppCookieDeviceTestHandler(api.AppCookieDevicePrefix))
	appcookiedeviceUrl := appcookiedeviceServer.URL + api.AppCookieDevicePrefix
	appinfoUrl := appinfoServer.URL + api.AppInfoPrefix
	urlBase, _ := url.Parse("http://example.com")
	sharelinkServer := httptest.NewServer(sharelink.NewShareLinkTestHandler(matchUrl, cookieUrl, appcookiedeviceUrl, appinfoUrl, tokenUrl, urlBase, api.ShareLinkPrefix))
	shorturlApp1 := strings.Replace(genurlResp.Url, "http://example.com", sharelinkServer.URL, -1)
	h, b, latency := testutil.GetHttpWithDSCookie(t, shorturlApp1, fakeUA, "www")
	if !strings.Contains(string(b), "AppInsStatus: '2',") {
		t.Error("AppInsStatus is not 2 (unclear)")
	}
	log.Debug("header:", h, "latency:", latency)

	//2. open app1 by universal link, bind wcookie<->uniqueID appid1+wcookie->uniqueID
	inappdataServer := httptest.NewServer(inappdata.NewInAppDataTestHandler(matchUrl, cookieUrl, appcookiedeviceUrl, genurlUrl, appinfoUrl, api.GetInAppDataPrefix))
	inappdataUrl := inappdataServer.URL + api.GetInAppDataPrefix
	inappdataReq := &inappdata.InAppDataPostBody{
		UniqueID: "uuu",
		WCookie:  "www",
	}
	inappdataResp := &inappdata.InAppDataResponseBody{}
	latency, err = testutil.PostJson(inappdataUrl+appID1, inappdataReq, fakeUA, "", inappdataResp)
	log.Debug("~~~~latency:", latency, "err:", err)
	log.Debug(inappdataResp)

	//3. click shorturl of app1 in wechat, should be installed
	h, b, latency = testutil.GetHttpWithDSCookie(t, shorturlApp1, fakeUA, "www")
	if !strings.Contains(string(b), "AppInsStatus: '1',") {
		t.Error("AppInsStatus is not 1 (installed)")
	}
	log.Debug("header:", h, "latency:", latency)

	//4. generate shorturl for a new app:app2
	genurlResp = &urlgenerator.GenURLResponseBody{}
	latency, err = testutil.PostJson(genurlUrl+appID2, genurlReq, fakeUA, "", genurlResp)
	if err != nil {
		t.Fatal("Failed to generate url, err:", err)
	}
	log.Debug("genurl latency:", latency, "genurl resp:", genurlResp)
	shorturlApp2 := strings.Replace(genurlResp.Url, "http://example.com", sharelinkServer.URL, -1)
	log.Debug("shorturl for app2:", shorturlApp2)

	//5. click shorturl of app2 in wechat, should be uninstalled
	h, b, latency = testutil.GetHttpWithDSCookie(t, shorturlApp2, fakeUA, "www")
	if !strings.Contains(string(b), "AppInsStatus: '0',") {
		t.Error("AppInsStatus is not 0 (not installed)")
	}
	log.Debug("header:", h, "latency:", latency)

	//6. open app2 by universal link, bind wcookie<->uniqueID appid2+wcookie->uniqueID
	inappdataReq = &inappdata.InAppDataPostBody{
		UniqueID: "uuu",
		WCookie:  "www",
	}
	inappdataResp = &inappdata.InAppDataResponseBody{}
	latency, err = testutil.PostJson(inappdataUrl+appID2, inappdataReq, fakeUA, "", inappdataResp)
	log.Debug("~~~~latency:", latency, "err:", err)
	log.Debug(inappdataResp)

	//7. click shorturl of app2 in wechat, should be installed
	h, b, latency = testutil.GetHttpWithDSCookie(t, shorturlApp2, fakeUA, "www")
	if !strings.Contains(string(b), "AppInsStatus: '1',") {
		t.Error("AppInsStatus is not 1 (installed)")
	}
	log.Debug("header:", h, "latency:", latency)

}

//TODO write tests for appInsStatus in native browser
func TestAppInsStatusBrowser(t *testing.T) {
	//first, open a link in wechat, should be AppNotInstalled
}

func RegisterApps(serverAddr string) {
	log.Debug("Register existing apps!")

	appInfos := []*appinfo.AppInfo{
		&appinfo.AppInfo{
			AppID:   "1652E90881C1FAE8",
			AppName: "testapp1",
			Android: appinfo.AppAndroidInfo{
				Scheme:             "ds1652E90881C1FAE8",
				Host:               "com.singulariti.testapp1",
				Pkg:                "com.singulariti.testapp1",
				DownloadUrl:        "http://baidu.com",
				IsDownloadDirectly: true,
			},
			Ios: appinfo.AppIosInfo{
				BundleID:            "com.singulariti.test",
				Scheme:              "tinggo",
				UniversalLinkEnable: true,
				DownloadUrl:         "https://itunes.apple.com/cn/app/ting-guo-wei-xin-wei-bo-you/id515901779?mt=8",
			},
			YYBEnable: true,
			YYBUrl:    "http://a.app.qq.com/o/simple.jsp?pkgname=com.haomee.superpower",
		},
		&appinfo.AppInfo{
			AppID:   "9aeea8a7f4c56ff4",
			AppName: "testapp2",
			Android: appinfo.AppAndroidInfo{
				Scheme: "ds9aeea8a7f4c56ff4",
				Host:   "com.singulariti.testapp2",
				Pkg:    "com.singulariti.testapp2",
			},
			Ios: appinfo.AppIosInfo{
				Scheme: "ds9aeea8a7f4c56ff4",
			},
			YYBEnable: false,
		},
	}

	for _, appinfo := range appInfos {
		b, err := json.Marshal(appinfo)
		if err != nil {
			log.Fatal(err)
		}
		client := httputil.GetNewClient()
		req, err := http.NewRequest("PUT", serverAddr+api.AppInfoPrefix+appinfo.AppID, strings.NewReader(string(b)))
		if err != nil {
			log.Debug("[Error], Register APP Info; Setup New Request to Register APP Info failed", err)
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("[Error],Register APP Info do request failed:", err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Fatal("[Error], Register app failed, response code is:", resp.StatusCode)
		}
		defer resp.Body.Close()

		fmt.Printf("Register succeed! \n")
	}
	return
}
