/*
* Package supports binding between cookie id and device id.
 */
package appcookiedevice

import (
	"time"

	in "github.com/MISingularity/deepshare2/pkg/instrumentation"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/storage"
)

const (
	cookiePairRKeyPrefix    = "cookie_pair:"
	installedAppsPKeyPrefix = "installed_apps:"
)

// Server allow bind device to a cookie and retrieve cookie with deviceID later.
type AppCookieDevice interface {
	BindAppCookieToDevice(appID, cookieID, deviceID string) error
	RefreshCookie(cookieID, newCookieID string) error
	GetDeviceID(appID, cookieID string) (string, error)
	IsAppInstalled(appID, cookieID string) (bool, error)
	pairCookies(cookie1, cookie2 string) error
	getPairedCookie(cookie string) (string, error)
}

func NewAppCookieDevice(skv storage.SimpleKV) AppCookieDevice {
	return &appCookieDevice{
		skv: skv,
	}
}

type appCookieDevice struct {
	skv storage.SimpleKV
}

func getKey(appID, cookieID string) string {
	return "acd:" + appID + ":" + cookieID
}

func (s *appCookieDevice) BindAppCookieToDevice(appID, cookieID, deviceID string) error {
	start := time.Now()

	if err := s.skv.Set([]byte(getKey(appID, cookieID)), []byte(deviceID)); err != nil {
		return err
	}

	in.PrometheusForAppCookieDevice.StoragePutDuration(start)

	// save appIDs under wcookie, to update appID:wcookie -> uniqueID bindings when wcookie is refreshed
	if err := s.skv.SAdd([]byte(installedAppsPKeyPrefix+cookieID), appID); err != nil {
		return err
	}
	return nil
}

func (s *appCookieDevice) GetDeviceID(appID, cookieID string) (string, error) {
	start := time.Now()

	v, err := s.skv.Get([]byte(getKey(appID, cookieID)))
	if err != nil {
		return "", err
	}

	in.PrometheusForAppCookieDevice.StorageGetDuration(start)
	return string(v), nil
}

func (s *appCookieDevice) IsAppInstalled(appID, cookieID string) (bool, error) {
	start := time.Now()

	v, err := s.skv.Get([]byte(getKey(appID, cookieID)))
	if err != nil {
		return false, err
	}

	isInstalled := false
	if string(v) != "" {
		isInstalled = true
	}
	in.PrometheusForAppCookieDevice.StorageGetDuration(start)
	return isInstalled, nil
}

//save dscookie value under wechat dscookie value
func (s *appCookieDevice) pairCookies(cookie1, cookie2 string) error {
	log.Debug("pair cookies:", cookie1, cookie2)
	if err := s.skv.Set([]byte(cookiePairRKeyPrefix+cookie1), []byte(cookie2)); err != nil {
		return err
	}
	if err := s.skv.Set([]byte(cookiePairRKeyPrefix+cookie2), []byte(cookie1)); err != nil {
		return err
	}
	return nil
}

//get dscookie value given wechat dscookie value
func (s *appCookieDevice) getPairedCookie(cookie string) (string, error) {
	log.Debug("Get paired cookie given:", cookie)
	if v, err := s.skv.Get([]byte(cookiePairRKeyPrefix + cookie)); err != nil {
		return "", err
	} else {
		return string(v), nil
	}
}

//update appID:cookieID -> uniqueID bindings and cookie pair when cookieID is refreshed
func (s *appCookieDevice) RefreshCookie(cookieID, newCookieID string) error {
	log.Debug("RefreshCookie:", cookieID, "->", newCookieID)
	appIDs, err := s.skv.SMembers([]byte(installedAppsPKeyPrefix + cookieID))
	if err != nil {
		log.Error(err)
		return err
	}
	log.Debug("RefreshCookie, affected AppIDs:", appIDs, cookieID, "->", newCookieID)

	//update appID:cookieID -> uniqueID bindings
	for _, appID := range appIDs {
		k := getKey(appID, cookieID)
		v, err := s.skv.Get([]byte(k))
		if err != nil {
			log.Error("RefreshCookie; Failed to get value under key:", k)
			continue
		}
		if len(v) > 0 {
			k := getKey(appID, newCookieID)
			err := s.skv.Set([]byte(k), v)
			if err != nil {
				log.Error("RefreshCookie; Failed to set value for key:", k)
				continue
			}
		}
	}

	//update cookie pair
	c, err := s.getPairedCookie(cookieID)
	if err != nil {
		log.Error(err)
		return err
	}
	if c != "" {
		log.Error(err)
		return s.pairCookies(c, newCookieID)
	}

	return nil
}
