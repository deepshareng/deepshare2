package appinfo

import (
	"reflect"
	"testing"

	"github.com/MISingularity/deepshare2/pkg/storage"
)

var (
	mockData = map[string]string{
		"app:7713337217A6E150": `{"AppID":"7713337217A6E150","UserConf":{"BgWeChatAndroidTipUrl":"www.androidwechattip","BgWeChatIosTipUrl":"www.ioswechattip"},"Android":{"Scheme":"deepshare","Host":"com.singulariti.deepsharedemo","Pkg":"com.singulariti.deepsharedemo","DownloadUrl":"","IsDownloadDirectly":true,"YYBEnable":false},"Ios":{"Scheme":"deepsharedemo","DownloadUrl":"", "BundleID":"bundleID1", "TeamID":"teamID1", "YYBEnableBelow9":true, "YYBEnableAbove9":true},"YYBUrl":"","YYBEnable":false,"Theme":"0"}`,
	}
	wantApp = map[string]AppInfo{
		"7713337217A6E150": AppInfo{
			AppID: "7713337217A6E150",
			Android: AppAndroidInfo{
				Scheme:             "deepshare",
				Host:               "com.singulariti.deepsharedemo",
				Pkg:                "com.singulariti.deepsharedemo",
				IsDownloadDirectly: true,
				YYBEnable:          false,
			},
			Ios: AppIosInfo{
				Scheme:          "deepsharedemo",
				BundleID:        "bundleID1",
				TeamID:          "teamID1",
				YYBEnableBelow9: true,
				YYBEnableAbove9: true,
			},
			UserConf: UserConfig{
				BgWeChatAndroidTipUrl: "www.androidwechattip",
				BgWeChatIosTipUrl:     "www.ioswechattip",
			},
			YYBEnable: false,
			Theme:     "0",
		},
	}
)

func TestGetAppInfo(t *testing.T) {
	appInfoDB := storage.NewMockSimpleKV()
	appInfoService := NewAppInfoService(appInfoDB)
	appID := "7713337217A6E150"
	app, _ := appInfoService.GetAppInfo(appID)
	if !reflect.DeepEqual(*app, wantApp[appID]) {
		t.Errorf("Get Wrong app info with Appid: %s, actual=%#v, want=%#v", appID, app, wantApp[appID])
	}
}

func TestGetNotExistAppInfo(t *testing.T) {
	appInfoDB := storage.NewMockSimpleKV()
	appInfoService := NewAppInfoService(appInfoDB)
	appID := "aaa"
	app, _ := appInfoService.GetAppInfo(appID)
	if app != nil {
		t.Errorf("Get Wrong app info with Appid: %s, actual=%#v, want nil", appID, app)
	}
}
