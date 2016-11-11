package appcookiedevice

import (
	"encoding/json"
	"net/http"
	"path"
	"testing"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"github.com/MISingularity/deepshare2/pkg/testutil"
)

func TestAppCookieDevice(t *testing.T) {
	tests := []struct {
		code     int
		appId    string
		cookieId string
		deviceId string
		body     string
	}{
		{
			http.StatusOK,
			"appId1",
			"cookieId1",
			"deviceId1",
			`{"unique_id":"deviceId1"}` + "\n",
		},
		{
			http.StatusOK,
			"appId1",
			"cookieId1",
			"",
			`{"unique_id":""}` + "\n",
		},
	}

	appCookieDev := NewAppCookieDevice(storage.NewInMemSimpleKV())
	handler := newAppCookieDeviceHandler(appCookieDev, nil, api.AppCookieDevicePrefix)

	for i, tt := range tests {
		urlPostRequest := "http://" + path.Join("fds.so", api.AppCookieDevicePrefix, tt.appId, tt.cookieId)
		w := testutil.HandleWithBody(handler, "PUT", urlPostRequest, tt.body)

		if w.Code != tt.code {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.code)
		}

		// We can only get valid tokens on successful request.
		if tt.code == http.StatusOK {
			urlGetRequest := "http://" + path.Join("fds.so", api.AppCookieDevicePrefix, tt.appId, tt.cookieId)
			r := testutil.HandleWithBody(handler, "GET", urlGetRequest, "")
			respBody := DeviceInfo{}
			if err := json.Unmarshal(r.Body.Bytes(), &respBody); err != nil {
				t.Fatalf("#%d: Unmarshall failed", i)
			}
			respUniqueId := respBody.UniqueId

			if respUniqueId != tt.deviceId {
				t.Errorf("#%d, deviceId = %s, want = %s", i, respUniqueId, tt.deviceId)
			}
		}
	}
}

func TestNoDevice(t *testing.T) {
	tests := []struct {
		code     int
		appId    string
		cookieId string
	}{
		{
			http.StatusOK,
			"appId1",
			"cookieId1",
		},
	}

	appCookieDev := NewAppCookieDevice(storage.NewInMemSimpleKV())
	handler := newAppCookieDeviceHandler(appCookieDev, nil, api.AppCookieDevicePrefix)

	for i, tt := range tests {

		urlGetRequest := "http://" + path.Join("fds.so", api.AppCookieDevicePrefix, tt.appId, tt.cookieId)
		r := testutil.HandleWithBody(handler, "GET", urlGetRequest, "")
		respBody := DeviceInfo{}
		if err := json.Unmarshal(r.Body.Bytes(), &respBody); err != nil {
			t.Fatalf("#%d: Unmarshall failed", i)
		}
		respUniqueId := respBody.UniqueId

		if respUniqueId != "" {
			t.Errorf("#%d, deviceId = %s, want = %s", i, respUniqueId, "")
		}

	}
}

func TestPairCookiesHandler(t *testing.T) {
	type req struct {
		method string
		path   string
		body   string
	}
	tests := []struct {
		prepareReqs []req
		path        string
		wcode       int
		wbody       string
	}{
		{
			[]req{
				{"PUT", "/appid/c1", `{"unique_id":"u"}` + "\n"},
				{"POST", "/paircookies/appid", `{"cookie1":"c1","cookie2":"c2"}` + "\n"},
			},
			"/appid/c2",
			http.StatusOK,
			`{"unique_id":"u"}` + "\n",
		},
		{
			[]req{
				{"POST", "/paircookies/appid", `{"cookie1":"cc1","cookie2":"cc2"}` + "\n"},
				{"PUT", "/appid/cc1", `{"unique_id":"uu"}` + "\n"},
			},
			"/appid/cc2",
			http.StatusOK,
			`{"unique_id":"uu"}` + "\n",
		},
		{
			[]req{
				{"PUT", "/appid/ccc2", `{"unique_id":"uuu"}` + "\n"},
				{"POST", "/paircookies/appid", `{"cookie1":"ccc1","cookie2":"ccc2"}` + "\n"},
			},
			"/appid/ccc1",
			http.StatusOK,
			`{"unique_id":"uuu"}` + "\n",
		},
		{
			[]req{
				{"POST", "/paircookies/appid", `{"cookie1":"cccc1","cookie2":"cccc2"}` + "\n"},
				{"PUT", "/appid/cccc1", `{"unique_id":"uuuu"}` + "\n"},
			},
			"/appid/cccc1",
			http.StatusOK,
			`{"unique_id":"uuuu"}` + "\n",
		},
	}
	appCookieDev := NewAppCookieDevice(storage.NewInMemSimpleKV())
	handler := newAppCookieDeviceHandler(appCookieDev, nil, api.AppCookieDevicePrefix)

	for i, tt := range tests {
		//prepare requests
		for _, req := range tt.prepareReqs {
			u := "http://" + path.Join("example.com", api.AppCookieDevicePrefix, req.path)
			testutil.MustHandleWithBodyOK(handler, req.method, u, req.body)
		}
		u := "http://" + path.Join("example.com", api.AppCookieDevicePrefix, tt.path)
		r := testutil.HandleWithBody(handler, "GET", u, "")
		if r.Code != tt.wcode {
			t.Errorf("#%d http code = %d, want = %d\n", i, r.Code, tt.wcode)
		}
		body := string(r.Body.String())
		if body != tt.wbody {
			t.Errorf("#%d responsed body = %s, want = %s\n", i, body, tt.wbody)
		}
	}
}

func TestAppCookieDeviceHandler_RefreshCookie(t *testing.T) {
	type req struct {
		method string
		path   string
		body   string
	}
	tests := []struct {
		prepareReqs []req
		path        string
		wcode       int
		wbody       string
	}{
		{
			[]req{
				{"PUT", "/appid/c1", `{"unique_id":"u"}` + "\n"},
				{"POST", "/refreshcookie", `{"cookie":"c1","new_cookie":"c1_new"}` + "\n"},
			},
			"/appid/c1_new",
			http.StatusOK,
			`{"unique_id":"u"}` + "\n",
		},
		{
			[]req{
				{"PUT", "/appid/cc2", `{"unique_id":"uu"}` + "\n"},
				{"POST", "/paircookies/appid", `{"cookie1":"cc1","cookie2":"cc2"}` + "\n"},
				{"POST", "/refreshcookie", `{"cookie":"cc1","new_cookie":"cc1_new"}` + "\n"},
			},
			"/appid/cc1_new",
			http.StatusOK,
			`{"unique_id":"uu"}` + "\n",
		},
	}
	appCookieDev := NewAppCookieDevice(storage.NewInMemSimpleKV())
	handler := newAppCookieDeviceHandler(appCookieDev, nil, api.AppCookieDevicePrefix)

	for i, tt := range tests {
		//prepare requests
		for _, req := range tt.prepareReqs {
			u := "http://" + path.Join("example.com", api.AppCookieDevicePrefix, req.path)
			testutil.MustHandleWithBodyOK(handler, req.method, u, req.body)
		}
		u := "http://" + path.Join("example.com", api.AppCookieDevicePrefix, tt.path)
		r := testutil.HandleWithBody(handler, "GET", u, "")
		if r.Code != tt.wcode {
			t.Errorf("#%d http code = %d, want = %d\n", i, r.Code, tt.wcode)
		}
		body := string(r.Body.String())
		if body != tt.wbody {
			t.Errorf("#%d responsed body = %s, want = %s\n", i, body, tt.wbody)
		}
	}
}
