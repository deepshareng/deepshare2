package urlgenerator

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"testing"

	"strings"

	"net/http/httptest"

	"regexp"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/appinfo"
	"github.com/MISingularity/deepshare2/deepshared/shorturl"
	"github.com/MISingularity/deepshare2/deepshared/token"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"github.com/MISingularity/deepshare2/pkg/testutil"
)

func TestShortURLPOST(t *testing.T) {
	tests := []struct {
		code int
		body string
	}{
		{
			http.StatusOK,
			`{"data":"{\"k1\": \"v1\",\"k2\": \"v2\"}","download_title":"dTitle","download_msg":"dMsg","redirect_url":"http://www.baidu.com/aaa?abc=eee","sdk_info":"android1.3.1","sender_id":"aabbcc"}` + "\n",
		},
		{
			http.StatusOK,
			`{"channels":["chan1_chantype1"],"download_title":"dTitle","download_msg":"dMsg","redirect_url":"http://www.baidu.com/aaa?abc=eee","sdk_info":"android1.3.1","forwarded_sender_id":"aabbcc"}` + "\n",
		},
	}

	tokenServer := httptest.NewServer(token.NewTokenTestHandler(api.TokenPrefix))
	handler, q, _ := newTestGenerateUrl(tokenServer.URL + api.TokenPrefix)
	urlRequest := "http://" + path.Join("fds.so", api.GenerateUrlPrefix, "7713337217A6E150")
	for i, tt := range tests {
		q.Reset()
		w := testutil.HandleWithBody(handler, "POST", urlRequest, tt.body)

		if w.Code != tt.code {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.code)
		}

		// We can only get valid tokens on successful request.
		if tt.code == http.StatusOK {
			respBody := GenURLResponseBody{}
			if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
				t.Fatalf("#%d: Unmarshall failed", i)
			}
			respUrlStr := respBody.Url
			respUrl, err := url.Parse(respUrlStr)
			if err != nil {
				t.Errorf("#%d, short url is illegal; url = %s", i, respUrlStr)
			}
			isShort, _ := shorturl.IsLegalShortFormat(respUrl)
			if !isShort {
				t.Errorf("#%d, short url is not in short format; url = %s", i, respUrlStr)
			}
		}
	}
}

func TestShortURLGET(t *testing.T) {
	tests := []struct {
		shortSeg string
		code     int
		body     string
	}{
		{
			"",
			http.StatusOK,
			`{"inapp_data":"{\"k1\": \"v1\",\"k2\": \"v2\"}","sender_id":"sss","channels":["1","2"]}` + "\n",
		},
	}

	tokenServer := httptest.NewServer(token.NewTokenTestHandler(api.TokenPrefix))
	handler, q, _ := newTestGenerateUrl(tokenServer.URL + api.TokenPrefix)
	urlPost := "http://" + path.Join("fds.so", api.GenerateUrlPrefix, "7713337217A6E150")
	for i, tt := range tests {
		//prepare: post to generate short seg
		q.Reset()
		r := testutil.MustHandleWithBodyOK(handler, "POST", urlPost, tt.body)
		urlBody := GenURLResponseBody{}
		if err := json.Unmarshal(r.Body.Bytes(), &urlBody); err != nil {
			t.Fatal(err)
		}

		u, err := url.Parse(urlBody.Url)
		if err != nil {
			t.Fatal(err)
		}
		isShort, _ := shorturl.IsLegalShortFormat(u)
		if !isShort {
			t.Fatalf("#%d, url is not in short format; url = %s", i, urlBody.Url)
		}
		fields := strings.Split(urlBody.Url, "/")
		if len(fields) < 2 {
			t.Fatal("wrong num of fields")
		}
		tt.shortSeg = fields[len(fields)-1]

		//test GET
		urlGet, err := url.Parse(urlPost)
		if err != nil {
			t.Fatal(err)
		}
		urlGet.Path = path.Join(urlGet.Path, tt.shortSeg)
		r = testutil.HandleWithBody(handler, "GET", urlGet.String(), "")

		if r.Code != tt.code {
			t.Errorf("#%d, response code invalid, code = %d; want = %d for url: %s", i, r.Code, tt.code, urlGet.String())
		}
		if r.Body.String() != tt.body {
			t.Errorf("#%d, response body invalid, body = %s; want = %s for url: %s", i, r.Body.String(), tt.body, urlGet.String())
		}
	}

}

func TestRawURLPOST(t *testing.T) {
	tests := []struct {
		code        int
		body        string
		encodeQuery string
	}{
		{
			http.StatusOK,
			`{"inapp_data":"{\"k1\": \"v1\",\"k2\": \"v2\"}","download_title":"dTitle","download_btn_text":"dBtnText",
"download_msg":"dMsg","download_url_android":"dlUrlAndroid","download_url_ios":"dlUrlIos","uninstall_url":"uninstallUrl","redirect_url":"http://www.baidu.com/aaa?abc=eee","sdk_info":"android1.3.1","sender_id":"aabbcc","is_short":false,"channels":["channel1","channel2"]}` + "\n",
			"channels=channel1%7Cchannel2&download_btn_text=dBtnText&download_msg=dMsg&download_title=dTitle&download_url_android=dlUrlAndroid&download_url_ios=dlUrlIos&inapp_data=%7B%22k1%22%3A+%22v1%22%2C%22k2%22%3A+%22v2%22%7D&redirect_url=http%3A%2F%2Fwww.baidu.com%2Faaa%3Fabc%3Deee&sdk_info=android1.3.1&sender_id=aabbcc&uninstall_url=uninstallUrl",
		},
		{
			http.StatusOK,
			`{"inapp_data":"{\"k1\": \"v1\",\"k2\": \"v2\"}","download_title":"dTitle","download_btn_text":"dBtnText",
"download_msg":"dMsg","download_url_android":"dlUrlAndroid","download_url_ios":"dlUrlIos","uninstall_url":"uninstallUrl","redirect_url":"http://www.baidu.com/aaa?abc=eee","sdk_info":"android1.3.1","sender_id":"aabbcc","is_short":false}` + "\n",
			"channels=&download_btn_text=dBtnText&download_msg=dMsg&download_title=dTitle&download_url_android=dlUrlAndroid&download_url_ios=dlUrlIos&inapp_data=%7B%22k1%22%3A+%22v1%22%2C%22k2%22%3A+%22v2%22%7D&redirect_url=http%3A%2F%2Fwww.baidu.com%2Faaa%3Fabc%3Deee&sdk_info=android1.3.1&sender_id=aabbcc&uninstall_url=uninstallUrl",
		},
		{
			http.StatusOK,
			`{"is_short":false}` + "\n",
			"channels=&download_btn_text=&download_msg=&download_title=&download_url_android=&download_url_ios=&inapp_data=&redirect_url=&sdk_info=&sender_id=&uninstall_url=",
		},
	}

	tokenServer := httptest.NewServer(token.NewTokenTestHandler(api.TokenPrefix))
	handler, q, _ := newTestGenerateUrl(tokenServer.URL + api.TokenPrefix)
	requestUrl := "http://" + path.Join("fds.so", api.GenerateUrlPrefix)
	for i, tt := range tests {
		q.Reset()
		w := testutil.HandleWithBody(handler, "POST", requestUrl, tt.body)

		if w.Code != tt.code {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.code)
		}

		// We can only get valid tokens on successful request.
		if tt.code == http.StatusOK {
			respBody := GenURLResponseBody{}
			if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
				t.Fatalf("#%d: Unmarshall failed", i)
			}
			urlStr := respBody.Url
			if resultUrl, err := url.Parse(urlStr); err != nil {
				t.Fatalf("#%d: Parse failed, url returned should in right format", i)
			} else {
				isShort, _ := shorturl.IsLegalShortFormat(resultUrl)
				if isShort {
					t.Errorf("#%d, url is not in raw format; url = %s", i, urlStr)
				}
				queryStr := resultUrl.RawQuery
				if queryStr != tt.encodeQuery {
					t.Errorf("#%d, query=%v, want=%v", i, queryStr, tt.encodeQuery)
				}
			}
		}
	}
}

func TestShortURLPOST_UseShortID(t *testing.T) {
	tests := []struct {
		appInfo      appinfo.AppInfo
		body         string
		wcode        int
		wbodyPattern string
	}{
		{
			appinfo.AppInfo{
				AppID:   "testAppID_longlong",
				ShortID: "SSSS",
			},
			`{"use_shortid":true,"inappdata":{"k":"v"}}` + "\n",
			http.StatusOK,
			`{"url":"http://fds.so/d/SSSS/\w{10}","path":"/d/SSSS/\w{10}"}`,
		},
		{
			appinfo.AppInfo{
				AppID:   "testAppID_longlong",
				ShortID: "SSSS",
			},
			`{"use_shortid":false,"inappdata":{"k":"v"}}` + "\n",
			http.StatusOK,
			`{"url":"http://fds.so/d/testAppID_longlong/\w{10}","path":"/d/testAppID_longlong/\w{10}"}`,
		},
	}

	db := storage.NewInMemSimpleKV()
	tokenServer := httptest.NewServer(token.NewTokenTestHandler(api.TokenPrefix))
	s := shorturl.NewUrlShortener(http.DefaultClient, db, tokenServer.URL)
	urlBase, _ := url.Parse("http://fds.so")
	q := new(bytes.Buffer)
	p := messaging.NewSimpleProducer(q)
	handler := newGenerateUrlHandler(http.DefaultClient, s, urlBase, p, api.GenerateUrlPrefix, tokenServer.URL)

	for i, tt := range tests {
		// first register app
		appinfoServer := httptest.NewServer(appinfo.NewAppInfoHandler(db, nil, api.AppInfoPrefix))
		b, err := json.Marshal(tt.appInfo)
		if err != nil {
			t.Fatalf("#%d Failed to marshal json, err: %v", i, err)
		}
		bytes.NewBuffer(b)
		req, err := http.NewRequest("PUT", appinfoServer.URL+api.AppInfoPrefix+tt.appInfo.AppID, bytes.NewReader(b))
		if err != nil {
			t.Fatalf("#%d Failed to new http.Request, err: %v", i, err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp == nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("#%d Failed to request to appinfo server, err: %v, resp.StatusCode: %d", i, err, resp.StatusCode)
		}

		//test handler
		q.Reset()
		requestUrl := "http://" + path.Join("fds.so", api.GenerateUrlPrefix, tt.appInfo.AppID)
		w := testutil.HandleWithBody(handler, "POST", requestUrl, tt.body)
		if w.Code != tt.wcode {
			t.Errorf("#%d HTTP status code = %d, want = %d\n", i, w.Code, tt.wcode)
		}

		reg := regexp.MustCompile(tt.wbodyPattern)
		if !reg.Match(w.Body.Bytes()) {
			t.Errorf("#%d response body = %s, want pattern = %s\n", i, w.Body.String(), tt.wbodyPattern)
		}
	}
}

// newTestDeepshareHandler sets up simple deepshare backend (storage, etc.)
// and returns a handler that we can use to test core logic.
func newTestGenerateUrl(tokenUrl string) (http.Handler, *bytes.Buffer, shorturl.UrlShortener) {
	s := shorturl.NewUrlShortener(http.DefaultClient, storage.NewInMemSimpleKV(), tokenUrl)
	urlBase, _ := url.Parse("http://fds.so")
	q := new(bytes.Buffer)
	p := messaging.NewSimpleProducer(q)
	return newGenerateUrlHandler(http.DefaultClient, s, urlBase, p, api.GenerateUrlPrefix, tokenUrl), q, s
}
