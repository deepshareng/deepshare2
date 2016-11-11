package httputil

import (
	"net/http"
	"testing"
)

func TestAppendPath(t *testing.T) {
	testData := []struct {
		oriUrl     string
		newPathSeg string
		newUrl     string
	}{
		{
			"http://fds.so/v2/match",
			"appid1",
			"http://fds.so/v2/match/appid1",
		},
		{
			"https://fds.soo/v2/devicecookie/",
			"appid2",
			"https://fds.soo/v2/devicecookie/appid2",
		},
	}

	for i, testObj := range testData {
		newUrlStr, err := AppendPath(testObj.oriUrl, testObj.newPathSeg)
		if err != nil {
			t.Fatalf("#%d Original Url is illegal = %s", i, testObj.oriUrl)
		}
		if newUrlStr != testObj.newUrl {
			t.Errorf("#%d NewURL = %#v, want = %#v", i, newUrlStr, testObj.newUrl)
		}
	}
}

func TestParseClientIP(t *testing.T) {
	tests := []struct {
		xForwardfor string
		remoteAddr  string
		expectedIP  string
	}{
		{
			xForwardfor: "",
			remoteAddr:  "testip:testport",
			expectedIP:  "testip",
		},
		{

			xForwardfor: "testip",
			remoteAddr:  "proxyip:testport",
			expectedIP:  "testip",
		},
		{
			xForwardfor: "testip, proxy1ip",
			remoteAddr:  "proxy2ip:testport",
			expectedIP:  "testip",
		},
		{
			xForwardfor: "",
			remoteAddr:  "",
			expectedIP:  "",
		},
	}

	for i, tt := range tests {
		r := &http.Request{
			Header: make(http.Header),
		}
		if tt.xForwardfor != "" {
			r.Header.Set("X-Forwarded-For", tt.xForwardfor)
		}
		r.RemoteAddr = tt.remoteAddr

		ip := ParseClientIP(r)
		if ip != tt.expectedIP {
			t.Errorf("#%d parse client ip failed, ip = %s, want = %s\n", i, ip, tt.expectedIP)
		}
	}
}
