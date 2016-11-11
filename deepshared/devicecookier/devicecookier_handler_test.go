package devicecookier

import (
	"net/http"
	"testing"

	"path"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"github.com/MISingularity/deepshare2/pkg/testutil"
)

func TestDeviceCookieHandler(t *testing.T) {
	tests := []struct {
		method string
		path   string
		body   string
		wcode  int
		wbody  string
	}{
		{ // 0
			"PUT",
			ApiCookiesPrefix + "/" + HardwareIDPrefix + "hhh",
			`{"cookie_id":"ccc"}` + "\n",
			http.StatusOK,
			"",
		},
		{ // 1
			"PUT",
			ApiCookiesPrefix + "/" + UniqueIDPrefix + "uuu",
			`{"cookie_id":"ccc"}` + "\n",
			http.StatusOK,
			"",
		},
		{ // 2
			"PUT",
			ApiCookiesPrefix + "/" + UniqueIDPrefixWCookie + "uuu",
			`{"cookie_id":"www"}` + "\n",
			http.StatusOK,
			"",
		},
		{ // 3
			"PUT",
			ApiDevicesPrefix + "/" + CookieIDPrefix + "ccc",
			`{"device_id":"uuu"}` + "\n",
			http.StatusOK,
			"",
		},
		{ // 4
			"PUT",
			ApiDevicesPrefix + "/" + WCookieIDPrefix + "www",
			`{"device_id":"uuu"}` + "\n",
			http.StatusOK,
			"",
		},
		{ // 5
			"GET",
			ApiCookiesPrefix + "/" + HardwareIDPrefix + "hhh",
			"",
			http.StatusOK,
			`{"cookie_id":"ccc"}` + "\n",
		},
		{ // 6
			"GET",
			ApiCookiesPrefix + "/" + UniqueIDPrefix + "uuu",
			"",
			http.StatusOK,
			`{"cookie_id":"ccc"}` + "\n",
		},
		{ // 7
			"GET",
			ApiDevicesPrefix + "/" + CookieIDPrefix + "ccc",
			"",
			http.StatusOK,
			`{"device_id":"uuu"}` + "\n",
		},
		{ // 8
			"GET",
			ApiDevicesPrefix + "/" + WCookieIDPrefix + "www",
			"",
			http.StatusOK,
			`{"device_id":"uuu"}` + "\n",
		},
	}

	dc := NewDeviceCookier(storage.NewInMemSimpleKV())
	dch := newDeviceCookieHandler(dc, nil, api.DeviceCookiePrefix)

	for i, tt := range tests {
		u := "http://" + path.Join("example.com", api.DeviceCookiePrefix, tt.path)
		r := testutil.HandleWithBody(dch, tt.method, u, tt.body)
		if r.Code != tt.wcode {
			t.Errorf("#%d http returns wrong code = %d, want = %d\n", i, r.Code, tt.wcode)
		}
		if r.Body.String() != tt.wbody {
			t.Errorf("#%d http returns wrong body = %s, want = %s\n", i, r.Body.String(), tt.wbody)
		}
	}

}
