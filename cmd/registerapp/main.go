package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/appinfo"
	"github.com/MISingularity/deepshare2/pkg/httputil"
)

const (
	serverUrl = "http://127.0.0.1:8080"
)

func main() {

	RegisterApps()
	return
}

func RegisterApps() {
	log.Println("Register existing apps!")

	appInfos := []*appinfo.AppInfo{
		&appinfo.AppInfo{
			AppID:   "e107dd56b99058c0",
			AppName: "超能部",
			Android: appinfo.AppAndroidInfo{
				Scheme:             "dse107dd56b99058c0",
				Host:               "com.haomee.superpower",
				Pkg:                "com.haomee.superpower",
				DownloadUrl:        "http://test.chaonengbu.haomee.cn/app-release.apk",
				IsDownloadDirectly: true,
			},
			Ios: appinfo.AppIosInfo{
				BundleID:            "com.haomee.superpower",
				Scheme:              "tinggo",
				UniversalLinkEnable: true,
				DownloadUrl:         "https://itunes.apple.com/cn/app/ting-guo-wei-xin-wei-bo-you/id515901779?mt=8",
			},
			YYBEnable: true,
			YYBUrl:    "http://a.app.qq.com/o/simple.jsp?pkgname=com.haomee.superpower",
		},
		&appinfo.AppInfo{
			AppID: "f709f09576216199",
			Android: appinfo.AppAndroidInfo{
				Scheme: "dsf709f09576216199",
				Host:   "com.singulariti.deepsharedemo",
				Pkg:    "com.singulariti.deepsharedemo",
			},
			Ios: appinfo.AppIosInfo{
				Scheme: "dsf709f09576216199",
			},
			YYBEnable: false,
		},
		&appinfo.AppInfo{
			AppID: "47B486A671F0BE82",
			Android: appinfo.AppAndroidInfo{
				Scheme:      "tinggo",
				Host:        "open",
				Pkg:         "com.xianguo.tingguo",
				DownloadUrl: "http://fruitlab.net/download_android",
			},
			Ios: appinfo.AppIosInfo{
				Scheme:      "tinggo",
				DownloadUrl: "https://itunes.apple.com/cn/app/ting-guo-wei-xin-wei-bo-you/id515901779?mt=8",
			},
			YYBEnable: true,
			YYBUrl:    "http://a.app.qq.com/o/simple.jsp?pkgname=com.xianguo.tingguo",
		},
		&appinfo.AppInfo{
			AppID: "EE26D0331DD04D54",
			Android: appinfo.AppAndroidInfo{
				Scheme:      "deepshare",
				Host:        "com.incn.yida",
				Pkg:         "com.incn.yida",
				DownloadUrl: "http://yeeda-video.oss-cn-beijing.aliyuncs.com/android/YEEDA.apk",
			},
			Ios: appinfo.AppIosInfo{
				Scheme:              "dsEE26D0331DD04D54",
				DownloadUrl:         "https://itunes.apple.com/cn/app/yi-da/id994848419?mt=8",
				UniversalLinkEnable: true,
			},
			YYBEnable: false,
			YYBUrl:    "http://a.app.qq.com/o/simple.jsp?pkgname=com.xianguo.tingguo",
		},
		&appinfo.AppInfo{
			AppID: "38CCA4C77072DDC9",
			Android: appinfo.AppAndroidInfo{
				Scheme: "ds38CCA4C77072DDC9",
				Host:   "com.croath.LRB",
				Pkg:    "com.croath.LRB",
			},
			Ios: appinfo.AppIosInfo{
				Scheme:              "ds38CCA4C77072DDC9",
				DownloadUrl:         "https://itunes.apple.com/us/app/linux-shou-ce/id593579235?mt=8",
				UniversalLinkEnable: true,
				TeamID:              "AMV535E2BX",
				BundleID:            "com.singularity.linuxcommand",
			},
			YYBEnable: false,
		},
		&appinfo.AppInfo{
			AppID:   "1652E90881C1FAE8",
			ShortID: "sss",
			Android: appinfo.AppAndroidInfo{
				Scheme: "ds1652E90881C1FAE8",
				Host:   "com.singularity.apitest",
				Pkg:    "com.singularity.apitest",
			},
			Ios: appinfo.AppIosInfo{
				Scheme:              "ds1652E90881C1FAE8",
				DownloadUrl:         "https://itunes.apple.com/us/app/linux-shou-ce/id593579235?mt=8",
				UniversalLinkEnable: true,
				TeamID:              "testTeam",
				BundleID:            "com.singularity.apitest",
			},
			YYBEnable: false,
		},
	}

	for _, appinfo := range appInfos {
		b, err := json.Marshal(appinfo)
		if err != nil {
			log.Fatal(err)
		}
		client := httputil.GetNewClient()
		req, err := http.NewRequest("PUT", serverUrl+api.AppInfoPrefix+appinfo.AppID, strings.NewReader(string(b)))
		if err != nil {
			log.Print("[Error], Register APP Info; Setup New Request to Register APP Info failed", err)
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Panic("[Error],Register APP Info do request failed:", err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Panic("[Error], Register app failed, response code is:", resp.StatusCode)
		}
		defer resp.Body.Close()

		fmt.Printf("Register succeed! \n")
	}
	return
}
