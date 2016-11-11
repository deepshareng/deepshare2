package appinfo

import (
	"encoding/json"

	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/storage"
)

//TODO: should not expose these interface
type AppInfoService interface {
	GetAppInfo(appID string) (*AppInfo, error)
	SetAppInfo(appID string, appInfo *AppInfo) error
}

type appInfoService struct {
	db storage.SimpleKV
}

func NewAppInfoService(appInfoDB storage.SimpleKV) *appInfoService {
	return &appInfoService{
		db: appInfoDB,
	}
}

//Get App struct by appID
func (appInfoService *appInfoService) GetAppInfo(appID string) (*AppInfo, error) {
	appInfo, err := appInfoService.db.Get([]byte(keyAppIDPrefix + appID))
	if err != nil {
		log.Error("AppInfo get app data error:", err)
		return nil, err
	}

	log.Debugf("Appinfo returned value: %s ; with appID: %s", string(appInfo), appID)
	if appInfo == nil || len(appInfo) == 0 {
		return nil, nil
	}
	app := AppInfo{}
	if err := json.Unmarshal(appInfo, &app); err != nil {
		log.Error("AppInfo Failed to unmarshal app data. appid:", appID, ", json data:", string(appInfo))
		return nil, err
	}
	return &app, nil
}

//Set App struct by appID
func (appInfoService *appInfoService) SetAppInfo(appID string, appinfo *AppInfo) error {
	appInfo, err := json.Marshal(appinfo)
	if err != nil {
		log.Error("AppInfo get app data error:", err)
		return err
	}

	return appInfoService.db.Set([]byte(keyAppIDPrefix+appID), []byte(appInfo))
}
