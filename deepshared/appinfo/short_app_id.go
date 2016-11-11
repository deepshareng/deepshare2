package appinfo

import (
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/storage"
)

//Get appID by shortID
func GetAppID(db storage.SimpleKV, shortID string) (string, error) {
	v, err := db.HGet([]byte(keyShortIDToAppID), shortID)
	if err != nil {
		log.Error("AppInfo get appid by shortid err:", err)
		return "", err
	}
	return string(v), nil
}

//Get shortID by appID
func GetShortID(db storage.SimpleKV, appID string) (string, error) {
	v, err := db.HGet([]byte(keyAppIDToShortID), appID)
	if err != nil {
		log.Error("AppInfo get shortid by appid err:", err)
		return "", err
	}
	return string(v), nil
}

//Save AppID <-> shortID
func SetAppIDShortIDPair(db storage.SimpleKV, shortID string, appID string) error {
	if err := db.HSet([]byte(keyShortIDToAppID), shortID, []byte(appID)); err != nil {
		return err
	}
	if err := db.HSet([]byte(keyAppIDToShortID), appID, []byte(shortID)); err != nil {
		return err
	}
	return nil
}
