/*
* Package supports binding between cookie id and device id.
 */
package devicecookier

import (
	"time"

	in "github.com/MISingularity/deepshare2/pkg/instrumentation"
	"github.com/MISingularity/deepshare2/pkg/storage"
)

// Server allow bind device to a cookie and retrieve cookie with deviceID later.
type DeviceCookier interface {
	//deviceID -> cookieID
	SaveCookieID(kDeviceID string, vCookieID string) error
	//cookieID -> deviceID
	SaveDeviceID(kCookieID string, vDeviceID string) error
	GetCookieID(kDeviceId string) (string, error)
	GetDeviceID(kCookieId string) (string, error)
}

func NewDeviceCookier(skv storage.SimpleKV) DeviceCookier {
	return &deviceCookier{
		skv: skv,
	}
}

type deviceCookier struct {
	skv storage.SimpleKV
}

func (s *deviceCookier) SaveCookieID(kDeviceID string, vCookieID string) error {
	// deviceID -> cookieID
	start := time.Now()

	if err := s.skv.Set([]byte(kDeviceID), []byte(vCookieID)); err != nil {
		return err
	}

	in.PrometheusForDeviceCookier.StoragePutDuration(start)

	return nil
}

func (s *deviceCookier) SaveDeviceID(kCookieID string, vDeviceID string) error {
	// cookieID -> deviceID
	start := time.Now()

	if err := s.skv.Set([]byte(kCookieID), []byte(vDeviceID)); err != nil {
		return err
	}

	in.PrometheusForDeviceCookier.StoragePutDuration(start)

	return nil
}

func (s *deviceCookier) GetCookieID(deviceID string) (string, error) {
	start := time.Now()

	v, err := s.skv.Get([]byte(deviceID))
	if err != nil {
		return "", err
	}

	in.PrometheusForDeviceCookier.StorageGetDuration(start)

	return string(v), nil
}

func (s *deviceCookier) GetDeviceID(cookieID string) (string, error) {
	start := time.Now()

	v, err := s.skv.Get([]byte(cookieID))
	if err != nil {
		return "", err
	}

	in.PrometheusForDeviceCookier.StorageGetDuration(start)

	return string(v), nil
}
