package devicecookier

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"strconv"

	"path"

	"github.com/MISingularity/deepshare2/pkg/httputil"
	"github.com/MISingularity/deepshare2/pkg/log"
	"golang.org/x/net/context"
)

func PutDeviceInfo(ctx context.Context, cli *http.Client, specificCookieUrl, kCookieID, deviceID string) error {
	if kCookieID == "" {
		errInfo := "cookieIDshould not be empty"
		log.Error("PutDeviceInfo;", errInfo)
		return errors.New(errInfo)
	}
	if deviceID == "" {
		errInfo := "deviceID should not be empty"
		log.Error("PutDeviceInfo;", errInfo)
		return errors.New(errInfo)
	}

	cookieUrlStr, err := httputil.AppendPath(specificCookieUrl, path.Join(ApiDevicesPrefix, kCookieID))
	if err != nil {
		log.Error("PutDeviceInfo; AppendPath failed, err:", err)
		return err
	}
	b, err := json.Marshal(DeviceInfo{
		DeviceID: deviceID,
	})
	if err != nil {
		log.Fatal("PutDeviceInfo; Marshal DeviceInfo failed", err)
		return err
	}
	log.Debugf("PutDeviceInfo; Request Cookie URL string is %s", cookieUrlStr)
	req, err := http.NewRequest("PUT", cookieUrlStr, bytes.NewReader(b))

	if err != nil {
		log.Error("PutDeviceInfo; Setup New Request to Device cookie service failed", err)
		return err
	}
	resp, err := cli.Do(req)
	if resp == nil {
		errInfo := "response from devicecookier is nil"
		log.Error("PutDeviceInfo;", errInfo)
		return errors.New(errInfo)
	}
	if resp.StatusCode != http.StatusOK {
		errInfo := "response status code from devicecookier is:" + strconv.Itoa(resp.StatusCode)
		log.Error("PutDeviceInfo;", errInfo)
		return errors.New(errInfo)
	}
	if err != nil {
		log.Error("PutDeviceInfo; Request to Device cookie service failed", err)
		return err
	}
	return nil
}

func PutCookieInfo(ctx context.Context, cli *http.Client, specificCookieUrl, kDeviceID, cookieID string) error {
	if kDeviceID == "" {
		errInfo := "deviceID should not be empty"
		log.Error("PutCookieInfo;", errInfo)
		return errors.New(errInfo)
	}
	if cookieID == "" {
		errInfo := "cookieIDshould not be empty"
		log.Error("PutCookieInfo;", errInfo)
		return errors.New(errInfo)
	}

	cookieUrlStr, err := httputil.AppendPath(specificCookieUrl, path.Join(ApiCookiesPrefix, kDeviceID))
	if err != nil {
		log.Error("PutCookieInfo; AppendPath failed, err:", err)
		return err
	}
	b, err := json.Marshal(CookieInfo{
		CookieID: cookieID,
	})
	if err != nil {
		log.Fatal("PutCookieInfo; Marshal CookieInfo failed", err)
		return err
	}
	log.Debugf("PutCookieInfo; Request Cookie URL string is %s", cookieUrlStr)
	req, err := http.NewRequest("PUT", cookieUrlStr, bytes.NewReader(b))

	if err != nil {
		log.Error("PutCookieInfo; Setup New Request to Device cookie service failed", err)
		return err
	}
	resp, err := cli.Do(req)
	if resp == nil {
		errInfo := "response from devicecookier is nil"
		log.Error("PutCookieInfo;", errInfo)
		return errors.New(errInfo)
	}
	if resp.StatusCode != http.StatusOK {
		errInfo := "response status code from devicecookier is:" + strconv.Itoa(resp.StatusCode)
		log.Error("PutCookieInfo;", errInfo)
		return errors.New(errInfo)
	}
	if err != nil {
		log.Error("PutCookieInfo; Request to Device cookie service failed", err)
		return err
	}
	return nil
}

func GetCookieID(ctx context.Context, cli *http.Client, specificCookieUrl, kUniqueID string) (cookieId string, err error) {
	cookieUrlStr, err := httputil.AppendPath(specificCookieUrl, path.Join(ApiCookiesPrefix, kUniqueID))
	if err != nil {
		log.Errorf("GetCookieID; cookieUrl %s is constructed by us, should not in wrong format: %v", specificCookieUrl, err)
		return "", err
	}
	log.Debugf("GetCookieID; Request Cookie URL string is %s", cookieUrlStr)
	req, err := http.NewRequest("GET", cookieUrlStr, nil)
	if err != nil {
		log.Error("GetCookieID; Setup New Request to Device cookie service failed", err)
		return "", err
	}
	resp, err := cli.Do(req)
	if err != nil {
		log.Error("GetCookieID; Request to Device cookie service failed", err)
		return "", err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	dcrb := CookieInfo{}
	err = decoder.Decode(&dcrb)
	if err != nil {
		log.Error("GetCookieID; Decode response from Device cookie service failed", err)
		return "", err
	}
	cookieId = dcrb.CookieID
	log.Debugf("GetCookieID; Cookie ID Got:%s", cookieId)
	return cookieId, nil
}

func GetUniqueID(ctx context.Context, cli *http.Client, specificCookieUrl, kCookieID string) (uniqueID string, err error) {
	cookieUrlStr, err := httputil.AppendPath(specificCookieUrl, path.Join(ApiDevicesPrefix, kCookieID))
	if err != nil {
		log.Errorf("GetUniqueID; cookieUrl %s is constructed by us, should not in wrong format: %v", specificCookieUrl, err)
		return "", err
	}

	log.Debugf("GetUniqueID; Request Cookie URL string is %s", cookieUrlStr)
	req, err := http.NewRequest("GET", cookieUrlStr, nil)
	if err != nil {
		log.Error("GetUniqueID; Setup New Request to Device cookie service failed", err)
		return "", err
	}
	resp, err := cli.Do(req)
	if err != nil {
		log.Error("GetUniqueID; Request to Device cookie service failed", err)
		return "", err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	dcrb := DeviceInfo{}
	err = decoder.Decode(&dcrb)
	if err != nil {
		log.Error("GetUniqueID; Decode response from Device cookie service failed", err)
		return "", err
	}
	uniqueID = dcrb.DeviceID
	log.Debugf("GetUniqueID; UniqueID Got:%s", uniqueID)
	return uniqueID, nil
}
