package sharelink

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"path"

	"net/http/httptest"

	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/shorturl"
	"github.com/MISingularity/deepshare2/deepshared/token"
	"github.com/MISingularity/deepshare2/pkg/testutil"
	"golang.org/x/net/context"
)

//mock url shortener
type mockUrlShortener struct {
}

func NewMockUrlShortener() shorturl.UrlShortener {
	return &mockUrlShortener{}
}

func (urlShortener *mockUrlShortener) ToShortURL(rawUrl *url.URL, namespace string, isPermanent bool, shortURLLifeTime time.Duration, useShortID bool) (*url.URL, error) {
	return nil, errors.New("Not implement")
}
func (urlShortener *mockUrlShortener) ToRawURL(requestUrl *url.URL, namespace string) (*url.URL, error) {
	if path.Base(requestUrl.Path) == "a0" {
		requestUrl.Path = path.Dir(requestUrl.Path)
		requestUrl.RawQuery = "download_title=aaa&download_msg=bbb&redirect_url=ccc&inapp_data=y"
	}
	return requestUrl, nil
}

func TestShortSharelinkAccess(t *testing.T) {
	tests := []struct {
		code                    int
		requestUrl              string
		remoteAddr              string
		header                  map[string]string
		mockRequestPath         string
		mockAppInfoResponseCode map[string]int
		mockAppInfoResponseBody map[string]string
	}{
		{ // 0
			http.StatusOK,
			`http://fds.so/d/7713337217A6E150/a0`,
			"ip1:port1",
			map[string]string{
				"User-Agent":      "Mozilla/5.0 (iPhone; CPU iPhone OS 7_0_3 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) Version/7.0 Mobile/11B511 Safari/9537.53",
				"X-Forwarded-For": "testip1",
			},
			`/7713337217A6E150`,
			map[string]int{`/7713337217A6E150`: http.StatusOK},
			map[string]string{
				`/7713337217A6E150`: `{"AppID":"7713337217A6E150","Android":{"Scheme":"deepshare","Host":"com.singulariti.deepsharedemo","Pkg":"com.singulariti.deepsharedemo","DownloadUrl":""},"Ios":{"Scheme":"deepsharedemo","DownloadUrl":""},"YYBUrl":"","YYBEnable":false}`,
				`/v2/tokens/cookie`: `{"token":"11111"}`,
			},
		},
		{ // 1
			http.StatusOK,
			`http://fds.so/d/7713337217A6E150/a0?k1=v1`,
			"ip2:port2",
			map[string]string{
				"User-Agent":      "Mozilla/5.0 (iPhone; CPU iPhone OS 7_0_3 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) Version/7.0 Mobile/11B511 Safari/9537.53",
				"X-Forwarded-For": "testip2",
			},
			`/7713337217A6E150`,
			map[string]int{`7713337217A6E150`: http.StatusOK},
			map[string]string{`*`: `{"AppID":"7713337217A6E150","Android":{"Scheme":"deepshare","Host":"com.singulariti.deepsharedemo","Pkg":"com.singulariti.deepsharedemo","DownloadUrl":""},"Ios":{"Scheme":"deepsharedemo","DownloadUrl":""},"YYBUrl":"","YYBEnable":false}`},
		},
		{ // 2
			http.StatusNotFound,
			`http://fds.so/d/7713337217A6E150/a1`,
			"ip3:port3",
			map[string]string{
				"User-Agent":      "Mozilla/5.0 (iPhone; CPU iPhone OS 7_0_3 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) Version/7.0 Mobile/11B511 Safari/9537.53",
				"X-Forwarded-For": "testip3",
			},
			``,
			map[string]int{`/v2/tokens/cookie`: http.StatusOK, "*": http.StatusNotFound},
			map[string]string{
				`/v2/tokens/cookie`: `{"token":"11111"}`,
				"*":                 ``,
			},
		},
		{ // 3
			http.StatusInternalServerError,
			`http://fds.so/d/7713337217A6E150?download_title=aaa&download_msg=bbb&redirect_url=ccc&inapp_data=y`,
			"ip4:port4",
			map[string]string{
				"User-Agent":      "Mozilla/5.0 (iPhone; CPU iPhone OS 7_0_3 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) Version/7.0 Mobile/11B511 Safari/9537.53",
				"X-Forwarded-For": "testip4",
			},
			``,
			map[string]int{"*": http.StatusOK},
			map[string]string{
				`/v2/tokens/cookie`: `{"token":"11111"}`,
				"*":                 `{"AppID":"7713337217A6E150":"com.singulariti.deepsharedemo"}`,
			},
		},
		{ // 4
			http.StatusOK,
			`http://fds.so/d/7713337217A6E150?download_title=aaa&download_msg=bbb&redirect_url=ccc&inapp_data=y`,
			"ip5:port5",
			map[string]string{
				"User-Agent":      "testua1",
				"X-Forwarded-For": "testip5",
			},
			``,
			map[string]int{"*": 0},
			map[string]string{"*": ``},
		},
		{ // 5
			http.StatusOK,
			`http://fds.so/d/7713337217A6E150/0`,
			"ip5:port5",
			map[string]string{
				"User-Agent":      "Mozilla/5.0 (iPhone; CPU iPhone OS 7_0_3 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) Version/7.0 Mobile/11B511 Safari/9537.53",
				"X-Forwarded-For": "testip5",
			},
			``,
			map[string]int{"*": 0},
			map[string]string{"*": ``},
		},
	}

	s := NewMockUrlShortener()
	urlBase, err := url.Parse("http://fds.so/")
	if err != nil {
		t.Fatalf("Parse Base url failed: %v", err)
	}
	setupServerEnv()
	for i, tt := range tests {
		mockCode := tt.mockAppInfoResponseCode
		if _, ok := mockCode["*"]; !ok {
			mockCode["*"] = http.StatusOK
		}
		mockResp := tt.mockAppInfoResponseBody
		if _, ok := mockResp["*"]; !ok {
			mockResp["*"] = ""
		}
		serverMock, clientMock, requestHistory, _ := testutil.MockResponse(mockCode, mockResp)
		tokenServer := httptest.NewServer(token.NewTokenTestHandler(api.TokenPrefix))
		defer serverMock.Close()
		handler := newShareLinkHandler(s, clientMock, serverMock.URL, "", "", serverMock.URL, tokenServer.URL+api.TokenPrefix, urlBase, "d", nil)
		w := testutil.HandleWithRequestInfo(handler, "GET", tt.requestUrl, "", tt.header, tt.remoteAddr)

		if w.Code != tt.code {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.code)
		}

		if tt.code == http.StatusOK {
			if tt.mockAppInfoResponseCode["*"] == 0 {
				//In this case, the request is from desktop or the short url seg is invalid
				//Should not have request, just return a html.
				if len(*requestHistory) != 0 {
					t.Errorf("#%d: UA is from desktop or a invalid short seg, should just return a html in this case", i)
				}
				return
			}
			reqAppInfoUrlStr := (*requestHistory)[tt.mockRequestPath]
			reqAppInfoUrl, err := url.Parse(reqAppInfoUrlStr)
			if err != nil {
				t.Fatalf("Request App info url is wrong: %v", err)
			}
			if reqAppInfoUrl.Path != tt.mockRequestPath {
				t.Errorf("#%d: HTTP App info request path = %s, want = %s", i, reqAppInfoUrl.Path, tt.mockRequestPath)
			}
		}
	}
}

func TestParseRequest(t *testing.T) {
	tests := []struct {
		rawUrl      string
		appId       string
		inAppData   string
		senderID    string
		channels    string
		sdkInfo     string
		dlInfo      DownloadInfo
		redirectUrl string
	}{
		{
			"http://fds.so/d/7713337217A6E150?inapp_data=%7B%22k1%22%3A+%22v1%22%2C%22k2%22%3A+%22v2%22%7D&download_title=aaa&download_btn_text=downloadBtnText&download_msg=bbb&download_url_android=dlUrlAndroid&download_url_ios=dlUrlIos&uninstall_url=uninstallUrl&redirect_url=ccc&sdk_info=android1.3.1&sender_id=aabbcc&channels=channel1|channel2",
			"7713337217A6E150",
			`{"k1": "v1","k2": "v2"}`,
			"aabbcc",
			"channel1|channel2",
			"android1.3.1",
			DownloadInfo{"aaa", "downloadBtnText", "bbb", "dlUrlIos", "dlUrlAndroid", "uninstallUrl"},
			"ccc",
		},
		{
			"http://fds.so/d/7713337217A6E150?inapp_data=&download_title=aaa&download_btn_text=downloadBtnText&download_msg=bbb&download_url_android=dlUrlAndroid&download_url_ios=dlUrlIos&uninstall_url=uninstallUrl&redirect_url=ccc&sdk_info=android1.3.1&sender_id=aabbcc&channels=",
			"7713337217A6E150",
			``,
			"aabbcc",
			"",
			"android1.3.1",
			DownloadInfo{"aaa", "downloadBtnText", "bbb", "dlUrlIos", "dlUrlAndroid", "uninstallUrl"},
			"ccc",
		},
		{
			"http://fds.so/d/7713337217A6E150?channels=&download_msg=&download_title=&download_btn_text=&inapp_data=&uninstall_url=&redirect_url=&sdk_info=&sender_id=",
			"7713337217A6E150",
			``,
			"",
			"",
			"",
			DownloadInfo{"", "", "", "", "", ""},
			"",
		},
	}
	for i, tt := range tests {
		requestUrl, err := url.ParseRequestURI(tt.rawUrl)
		if err != nil {
			t.Fatalf("parse url faild %v", err)
		}
		shareLinkHandler := &shareLinkHandler{}
		appID, inAppData, senderID, channels, sdkInfo, dlInfo, redirectUrl := shareLinkHandler.sl.parseRequest(context.TODO(), requestUrl)
		if appID != tt.appId {
			t.Errorf("#%d parse url faild; appID expect:%s, actual:%s", i, appID, tt.appId)
		}
		if inAppData != tt.inAppData {
			t.Errorf("#%d parse url faild; in app data actual:%s, expect:%s", i, inAppData, tt.inAppData)
		}
		if senderID != tt.senderID {
			t.Errorf("#%d parse url faild; senderID actual:%s, expect:%s", i, senderID, tt.senderID)
		}
		if channels != tt.channels {
			t.Errorf("#%d parse url faild; channels actual:%s, expect:%s", i, channels, tt.channels)
		}
		if sdkInfo != tt.sdkInfo {
			t.Errorf("#%d parse url faild; sdkInfo actual:%s, expect:%s", i, sdkInfo, tt.sdkInfo)
		}
		if dlInfo != tt.dlInfo {
			t.Fatalf("#%d parse url faild; download_info actual:%+v, expect:%+v", i, dlInfo, tt.dlInfo)
		}
		if redirectUrl != tt.redirectUrl {
			t.Fatalf("#%d parse url faild; redirect_url actual:%s, expect:%s", i, redirectUrl, tt.redirectUrl)
		}
	}
}

func TestExtractInfoFromShortUrl(t *testing.T) {
	tests := []struct {
		shortUrl string
		appID    string
		shortSeg string
	}{
		{
			`http://127.0.0.1:8080/d/7713337217A6E150/MzMwNjc0MjgxNzQ4ODg5Ng`,
			`7713337217A6E150`,
			`MzMwNjc0MjgxNzQ4ODg5Ng`,
		},
	}
	for i, tt := range tests {
		urlShort, err := url.Parse(tt.shortUrl)
		if err != nil {
			t.Fatalf("#%d parse url %s faild %v", i, tt.shortUrl, err)
		}
		appId, shortSeg := extractInfoFromShortUrl(urlShort)
		if appId != tt.appID {
			t.Errorf("#%d Extract AppID From ShortUrl faild; appId actual:%s, expect:%s", i, appId, tt.appID)
		}
		if shortSeg != tt.shortSeg {
			t.Errorf("#%d Extract ShortSeg From ShortUrl faild; appId actual:%s, expect:%s", i, shortSeg, tt.shortSeg)
		}
	}
}

func TestExtractAppIDFromRawUrl(t *testing.T) {
	tests := []struct {
		rawUrl string
		appID  string
	}{
		{
			`http://127.0.0.1:8080/d/7713337217A6E150?channels=&download_msg=&download_title=&inapp_data=%7B+%22key1%22%3A%22test_value1%22%2C%22key2%22%3A2+%7D&redirect_url=&sdk_info=ios1.1.2&sender_id=7E7B2568-B666-4577-A9DE-83A4ED8528B9`,
			`7713337217A6E150`,
		},
	}
	for i, tt := range tests {
		urlRaw, err := url.Parse(tt.rawUrl)
		if err != nil {
			t.Fatalf("#%d parse url %s faild %v", i, tt.rawUrl, err)
		}
		appId := extractAppIDFromRawUrl(urlRaw)
		if appId != tt.appID {
			t.Errorf("#%d Extract AppID From ShortUrl faild; appId actual:%s, expect:%s", i, appId, tt.appID)
		}
	}
}

func TestExtractDSTag(t *testing.T) {
	tests := []struct {
		url   string
		dsTag string
	}{
		{
			`https://fds.so/d/7713337217A6E150?ds_tag=123456`,
			`123456`,
		},
	}
	for i, tt := range tests {
		url, err := url.Parse(tt.url)
		if err != nil {
			t.Fatalf("#%d parse url %s faild %v", i, tt.url, err)
		}
		tag := extractDSTag(url)
		if tag != tt.dsTag {
			t.Errorf("#%d Extract DSTag From ShortUrl faild; dsTag actual:%s, expect:%s", i, tag, tt.dsTag)
		}
	}
}
