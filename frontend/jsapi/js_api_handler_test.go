package jsapi

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"fmt"

	"regexp"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/appcookiedevice"
	"github.com/MISingularity/deepshare2/deepshared/appinfo"
	"github.com/MISingularity/deepshare2/deepshared/devicecookier"
	"github.com/MISingularity/deepshare2/deepshared/match"
	"github.com/MISingularity/deepshare2/deepshared/token"
	"github.com/MISingularity/deepshare2/pkg/testutil"
)

type mockTokenGenerator struct {
	countShortSeg int
	countCookie   int
}

func (mtg *mockTokenGenerator) Generate(namespace string) (string, error) {
	token := ""
	if namespace == "appid" {
		token = fmt.Sprintf("a%d", mtg.countShortSeg)
		mtg.countShortSeg++
	}
	if namespace == "cookie" {
		token = fmt.Sprintf("b%d", mtg.countCookie)
		mtg.countCookie++
	}
	return token, nil
}

func newTestJsApiHandler() http.Handler {
	appinfoServer := httptest.NewServer(appinfo.NewAppInfoTestHandler(api.AppInfoPrefix))
	appinfoUrl := appinfoServer.URL + api.AppInfoPrefix
	matchServer := httptest.NewServer(match.NewMatchTestHandler(api.MatchPrefix))
	matchUrl := matchServer.URL + api.MatchPrefix
	cookieServer := httptest.NewServer(devicecookier.NewDeviceCookieTestHandler(api.DeviceCookiePrefix))
	cookieUrl := cookieServer.URL + api.DeviceCookiePrefix
	appcookiedeviceServer := httptest.NewServer(appcookiedevice.NewAppCookieDeviceTestHandler(api.AppCookieDevicePrefix))
	appcookiedeviceUrl := appcookiedeviceServer.URL + api.AppCookieDevicePrefix
	tokenServer := httptest.NewServer(token.NewTokenTestHandler(api.TokenPrefix))
	tokenUrl := tokenServer.URL + api.TokenPrefix

	urlBase, _ := url.Parse("http://example.com")
	handler := NewJsApiTestHandler(matchUrl, cookieUrl, appcookiedeviceUrl, appinfoUrl, tokenUrl, urlBase, api.ShareLinkPrefix)
	return handler
}

func TestJsApiPost(t *testing.T) {
	tests := []struct {
		path         string
		ua           string
		body         string
		wcode        int
		wbodyPattern string
	}{
		{
			"/v2/jsapi/appid",
			"Mozilla/5.0 (Linux; Android 5.1.1; Nexus 6 Build/LYZ28E) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.20 Mobile Safari/537.36",
			`[{"deeplink_id":"1","inapp_data":{"name":"n1"},"sender_id":"s1","channels":["ch1_x","ch1_y"],"download_title":"download_title1","download_btn_text":"download_btn_text1","download_msg":"download_msg1","download_url_ios":"","download_url_android":""},{"deeplink_id":"2","inapp_data":{"name":"n2"},"sender_id":"s2","channels":["ch2_x","ch2_y"],"download_title":"download_title2","download_btn_text":"download_btn_text2","download_msg":"download_msg2","download_url_ios":"","download_url_android":""}]` + "\n",
			http.StatusOK,
			`{"app_id":"appid","ds_urls":{"1":"http://example.com/d/appid/\w{10}","2":"http://example.com/d/appid/\w{10}"},"chrome_major":44,"is_android":true,"is_ios":false,"ios_major":0,"is_wechat":false,"is_weibo":false,"is_qq":false,"is_facebook":false,"is_twitter":false,"is_firefox":false,"is_qq_browser":false,"is_uc":false,"cannot_deeplink":false,"cannot_get_win_event":false,"cannot_go_market":false,"force_use_scheme":false,"app_name":"","icon_url":"","scheme":"","host":"","bundle_id":"","pkg":"","url":"","is_download_directly":false,"is_universal_link":false,"is_yyb_enable_ios_below_9":false,"is_yyb_enable_ios_above_9":false,"is_yyb_enable_android":false,"yyb_url":"","match_id":"\w{10}","timestamp":\d{10},"ds_tag":"","app_ins_status":0}` + "\n",
		},
	}
	handler := newTestJsApiHandler()
	for i, tt := range tests {
		w := testutil.HandleWithRequestInfo(handler, "POST", "http://example.com"+tt.path, tt.body, map[string]string{"User-Agent": tt.ua}, "")
		if w.Code != tt.wcode {
			t.Errorf("#%d jsapi response code = %d, want = %d\n", i, w.Code, tt.wcode)
		}
		reg := regexp.MustCompile(tt.wbodyPattern)
		if !reg.Match(w.Body.Bytes()) {
			t.Errorf("#%d jsapi response body = %s, want patteren = %s\n", i, w.Body.String(), tt.wbodyPattern)
		}
	}
}

func TestJsApiPostClicked(t *testing.T) {
	tests := []struct {
		path  string
		ua    string
		body  string
		wcode int
		wbody string
	}{
		{
			"/v2/jsapi/appid?clicked=true",
			"Mozilla/5.0 (Linux; Android 5.1.1; Nexus 6 Build/LYZ28E) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.20 Mobile Safari/537.36",
			`{"deeplink_id":"1","match_id":"b0","sender_id":"s1","channels":["ch1_x","ch1_y"],"inapp_data":{"name":"n1"}}` + "\n",
			http.StatusOK,
			`{"ok":true}` + "\n",
		},
	}
	handler := newTestJsApiHandler()
	for i, tt := range tests {
		w := testutil.HandleWithRequestInfo(handler, "POST", "http://example.com"+tt.path, tt.body, map[string]string{"User-Agent": tt.ua}, "")
		if w.Code != tt.wcode {
			t.Errorf("#%d jsapi(clicked) response code = %d, want = %d\n", i, w.Code, tt.wcode)
		}
		if w.Body.String() != tt.wbody {
			t.Errorf("#%d jsapi(clicked) response body = %s, want = %s\n", i, w.Body.String(), tt.wbody)
		}
	}
}
