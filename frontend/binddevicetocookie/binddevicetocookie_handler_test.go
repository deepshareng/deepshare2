package binddevicetocookie

import (
	"net/http"
	"testing"

	"net/url"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/pkg/testutil"
)

func TestInAppDataAccess(t *testing.T) {
	tests := []struct {
		code            int
		requestUrl      string
		remoteAddr      string
		header          map[string]string
		mockRequestPath string
		mockRequestBody string
	}{
		{
			http.StatusOK,
			`http://fds.so/v2/binddevicetocookie/uid1`,
			"ip1:port1",
			map[string]string{
				"User-Agent":      "test_useragent1",
				"X-Forwarded-For": "testip1",
				"Cookie":          "dscookie=dsck1",
			},
			`/v2/devicecookie/cookies/uid_uid1`,
			`{"cookie_id":"dsck1"}`,
		},
	}
	for i, tt := range tests {
		serverMock, clientMock, requestUrlHistory, requestBodyHistory := testutil.MockResponse(map[string]int{
			"*": http.StatusOK,
		}, map[string]string{
			"*": `` + "\n",
		})
		handler := newBindDeviceToCookieHandler(clientMock, serverMock.URL+api.DeviceCookiePrefix, "", api.BindDeviceToCookiePrefix)
		defer serverMock.Close()
		w := testutil.HandleWithRequestInfo(handler, "GET", tt.requestUrl, "", tt.header, tt.remoteAddr)
		if w.Code != tt.code {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.code)
		}
		if tt.code == http.StatusOK {
			reqUrlStr := (*requestUrlHistory)[tt.mockRequestPath]
			reqBodyStr := (*requestBodyHistory)[tt.mockRequestPath]
			reqUrl, err := url.Parse(reqUrlStr)
			if err != nil {
				t.Fatalf("Request match url is wrong: %v", err)
			}
			if reqUrl.Path != tt.mockRequestPath {
				t.Errorf("#%d: HTTP Device Cookie request path = %s, want = %s", i, reqUrl.Path, tt.mockRequestPath)
			}
			if reqBodyStr != tt.mockRequestBody {
				t.Errorf("#%d: HTTP Device Cookie request body = %s, want = %s", i, reqBodyStr, tt.mockRequestBody)
			}
		}
	}
}
