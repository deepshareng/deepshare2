package uainfo

import (
	"bytes"
	"testing"
)

func TestUAFingerPrinter(t *testing.T) {
	tests := []struct {
		ua  UAInfo
		out []byte
	}{
		{
			ua:  UAInfo{Ua: "test ua", Os: "testos", OsVersion: "testosv", Ip: "testip", Brand: "testbrand"},
			out: []byte("UA:testip_testos_testosv_testbrand"),
		},
		{ // if Ua string is not set, Transform() should return empty string
			ua:  UAInfo{Ua: "", Os: "testos", OsVersion: "testosv", Ip: "testip", Brand: "testbrand"},
			out: []byte(""),
		},
	}

	for i, tt := range tests {
		uf := NewUAFingerPrinter(&tt.ua)
		b := uf.Transform()
		if !bytes.Equal([]byte(b), tt.out) {
			t.Errorf("#%d: got = %s, want = %s", i, string(b), string(tt.out))
		}
	}
}

func TestObjectFucntion(t *testing.T) {
	tests := []struct {
		ua        UAInfo
		isAndroid bool
		isIos     bool
		iosVer    int
	}{
		{
			ua:        UAInfo{Ua: "test ua", Os: "testos", OsVersion: "testosv", Ip: "testip", Brand: "testbrand"},
			isAndroid: false,
			isIos:     false,
			iosVer:    0,
		},
		{
			ua:        UAInfo{Ua: "test ua", Os: "android", OsVersion: "testosv", Ip: "testip", Brand: "testbrand"},
			isAndroid: true,
			isIos:     false,
			iosVer:    0,
		},
		{ // if Ua string is not set, Transform() should return empty string
			ua:        UAInfo{Ua: "", Os: "ios", OsVersion: "testosv", Ip: "testip", Brand: "testbrand"},
			isAndroid: false,
			isIos:     true,
			iosVer:    0,
		},
		{ // if Ua string is not set, Transform() should return empty string
			ua:        UAInfo{Ua: "", Os: "ios", OsVersion: "8.1", Ip: "testip", Brand: "testbrand"},
			isAndroid: false,
			isIos:     true,
			iosVer:    8,
		},
	}
	for i, tt := range tests {
		if tt.isIos != tt.ua.IsIos() {
			t.Errorf("#%d: isIos got = %t, want = %t", i, tt.ua.IsIos(), tt.isIos)
		}
		if tt.isAndroid != tt.ua.IsAndroid() {
			t.Errorf("#%d: isAndroid got = %t, want = %t", i, tt.ua.IsAndroid(), tt.isAndroid)
		}
		if tt.iosVer != tt.ua.IosMajorVersion() {
			t.Errorf("#%d: ios major version got = %d, want = %d", i, tt.ua.IosMajorVersion(), tt.iosVer)
		}
	}
}
