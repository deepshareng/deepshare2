package devicecookier

import (
	"reflect"
	"testing"

	"github.com/MISingularity/deepshare2/pkg/storage"
)

func TestCookieIDNormal(t *testing.T) {
	tests := []struct {
		cookieID string
		deviceID string
	}{
		{
			"7713337217A6E150/MA==",
			"xyz",
		},
	}

	s := NewDeviceCookier(storage.NewInMemSimpleKV())
	for i, datum := range tests {
		err := s.SaveCookieID(datum.deviceID, datum.cookieID)
		if err != nil {
			t.Fatalf("#%d Parse rawUrl failed: %v, please check test cases", i, err)
		}
	}

	for i, datum := range tests {
		cookie, err := s.GetCookieID(datum.deviceID)
		if err != nil {
			t.Fatalf("#%d Has not key: %s", i, datum.deviceID)
		}

		if !reflect.DeepEqual(cookie, datum.cookieID) {
			t.Errorf("#%d: payload = %s, want = %s", i, cookie, datum.cookieID)
		}
	}
}

func TestCookieIDFalsePositive(t *testing.T) {
	// Test false positive.
	s := NewDeviceCookier(storage.NewInMemSimpleKV())
	cookie, _ := s.GetCookieID("ShouldHaveNoResult")
	if cookie != "" {
		t.Errorf("Does not expect value for ShouldHaveNoResult %s", cookie)
	}
}

func TestDeviceIDNormal(t *testing.T) {
	tests := []struct {
		cookieID string
		deviceID string
	}{
		{
			"ddd",
			"ccc",
		},
	}

	s := NewDeviceCookier(storage.NewInMemSimpleKV())
	for i, datum := range tests {
		err := s.SaveDeviceID(datum.cookieID, datum.deviceID)
		if err != nil {
			t.Fatalf("#%d Parse rawUrl failed: %v, please check test cases", i, err)
		}
	}

	for i, tt := range tests {
		deviceID, err := s.GetDeviceID(tt.cookieID)
		if err != nil {
			t.Fatalf("#%d Has not key: %s", i, tt.cookieID)
		}

		if !reflect.DeepEqual(deviceID, tt.deviceID) {
			t.Errorf("#%d: payload = %s, want = %s", i, deviceID, tt.deviceID)
		}
	}
}
