/*
* Package supports binding between cookie id and device id.
 */
package appcookiedevice

import (
	"reflect"
	"testing"

	"github.com/MISingularity/deepshare2/pkg/storage"
)

func TestCookieIDNormal(t *testing.T) {
	tests := []struct {
		appID    string
		cookieID string
		deviceID string
	}{
		{
			"app1",
			"cookie1",
			"device1",
		},
	}

	s := NewAppCookieDevice(storage.NewInMemSimpleKV())
	for i, datum := range tests {
		err := s.BindAppCookieToDevice(datum.appID, datum.cookieID, datum.deviceID)
		if err != nil {
			t.Fatalf("#%d Parse rawUrl failed: %v, please check test cases", i, err)
		}
	}

	for i, datum := range tests {
		deviceID, err := s.GetDeviceID(datum.appID, datum.cookieID)
		if err != nil {
			t.Fatalf("#%d Has not key: %s", i, datum.deviceID)
		}

		if !reflect.DeepEqual(deviceID, datum.deviceID) {
			t.Errorf("#%d: payload = %s, want = %s", i, deviceID, datum.deviceID)
		}
	}
}

func TestCookieIDFalsePositive(t *testing.T) {
	// Test false positive.
	s := NewAppCookieDevice(storage.NewInMemSimpleKV())
	deviceId, _ := s.GetDeviceID("Should", "HaveNoResult")
	if deviceId != "" {
		t.Errorf("Does not expect value for ShouldHaveNoResult %s", deviceId)
	}
}

func TestPairCookies(t *testing.T) {
	s := NewAppCookieDevice(storage.NewInMemSimpleKV())
	if err := s.pairCookies("c1", "c2"); err != nil {
		t.Fatal(err)
	}
	c2, err := s.getPairedCookie("c1")
	if err != nil {
		t.Fatal(err)
	}
	if c2 != "c2" {
		t.Error("got = ", c2, "want = c2")
	}

	c1, err := s.getPairedCookie("c2")
	if err != nil {
		t.Fatal(err)
	}
	if c1 != "c1" {
		t.Error("got = ", c1, "want = c1")
	}
}

func TestAppCookieDevice_RefreshCookie(t *testing.T) {
	s := NewAppCookieDevice(storage.NewInMemSimpleKV())
	if err := s.pairCookies("cc", "CC"); err != nil {
		t.Fatal(err)
	}
	if err := s.BindAppCookieToDevice("app1", "cc", "d1"); err != nil {
		t.Fatal(err)
	}
	if err := s.BindAppCookieToDevice("app2", "cc", "d2"); err != nil {
		t.Fatal(err)
	}

	//Refresh cookie: cc -> cc_new
	if err := s.RefreshCookie("cc", "cc_new"); err != nil {
		t.Fatal(err)
	}

	//GetDeviceID with the new cookie
	if d, err := s.GetDeviceID("app1", "cc_new"); err != nil {
		t.Fatal(err)
	} else if d != "d1" {
		t.Errorf("deviceID should not change after RefreshCookie, got = %s, want = %s\n", d, "d1")
	}
	if d, err := s.GetDeviceID("app2", "cc_new"); err != nil {
		t.Fatal(err)
	} else if d != "d2" {
		t.Errorf("deviceID should not change after RefreshCookie, got = %s, want = %s\n", d, "d2")
	}

	//getPairedCookie with the new cookie
	if c, err := s.getPairedCookie("cc_new"); err != nil {
		t.Fatal(err)
	} else if c != "CC" {
		t.Errorf("pairedCookie of the new cookie = %s, want = %s\n", c, "CC")
	}
	if c, err := s.getPairedCookie("CC"); err != nil {
		t.Fatal(err)
	} else if c != "cc_new" {
		t.Errorf("pairedCookie of CC = %s, want = %s\n", c, "cc_new")
	}
}
