package sharelink

import (
	"fmt"
	"html/template"
	"net/http"

	"time"

	"strings"

	"github.com/MISingularity/deepshare2/deepshared/appinfo"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg/log"
	"golang.org/x/net/context"
)

func injectInfo(appInfo *appinfo.AppInfo, appID string, matchID string, uaObj *uainfo.UAInfo, dlInfo DownloadInfo, redirectUrl, dsTag string, appInstallStatus int) map[string]interface{} {
	info := make(map[string]interface{})
	//From UA info
	info["Chrome_major"] = uaObj.ChromeMajor
	info["AppID"] = appID
	os := uaObj.Os
	info["isAndroid"] = (os == "android")
	info["isIOS"] = (os == "ios")
	info["isWechat"] = uaObj.IsWechat
	info["isWeibo"] = uaObj.IsWeibo
	info["isFacebook"] = uaObj.IsFacebook
	info["isTwitter"] = uaObj.IsTwitter

	info["isQQ"] = uaObj.IsQQ
	info["isFirefox"] = uaObj.IsFirefox
	info["isQQBrowser"] = uaObj.IsQQBrowser
	info["Ios_major"] = uaObj.IosMajorVersion()
	info["CannotDeeplink"] = uaObj.CannotDeeplink
	info["CannotGetWinEvent"] = uaObj.CannotGetWindowsEvent
	info["CannotGoMarket"] = uaObj.CannotGoMarket
	info["ForceUseScheme"] = uaObj.ForceUseScheme
	info["isUC"] = strings.Contains(uaObj.Browser, "UC Browser")

	//From App info
	info["AppName"] = appInfo.AppName
	info["IconUrl"] = appInfo.IconUrl
	if appInfo.Theme == "" {
		info["Theme"] = "0"
	} else {
		info["Theme"] = appInfo.Theme
	}
	info["UserConf_Bg_WechatAndroidTip_url"] = appInfo.UserConf.BgWeChatAndroidTipUrl
	info["UserConf_Bg_WechatIosTip_url"] = appInfo.UserConf.BgWeChatIosTipUrl

	//Init below fields in case of the value is changed to "NO_VALUE" when Execute the html template
	info["Scheme"] = ""
	info["Host"] = ""
	info["BundleID"] = ""
	info["Pkg"] = ""
	info["Url"] = ""
	info["isDownloadDirectly"] = false
	info["isUniversallink"] = false
	info["YYB_url"] = appInfo.YYBUrl
	info["YYB_Enable_Android"] = false
	info["YYB_Enable_Ios_Above_9"] = false
	info["YYB_Enable_Ios_below_9"] = false

	if uaObj.IsAndroid() {
		info["Scheme"] = appInfo.Android.Scheme /*"deepshare"*/
		info["Host"] = appInfo.Android.Host     /*"com.singulariti.deepsharedemo"*/
		info["Pkg"] = appInfo.Android.Pkg       /*"com.singulariti.deepsharedemo"*/
		if dlInfo.DownloadUrlAndroid != "" {
			info["Url"] = dlInfo.DownloadUrlAndroid
		} else {
			info["Url"] = appInfo.Android.DownloadUrl /*"https://play.google.com/store/apps/details?id=com.xianguo.tingguo"*/
		}
		info["isDownloadDirectly"] = appInfo.Android.IsDownloadDirectly
		if appInfo.YYBEnable == true {
			//This is for forward compatbility;
			//At first, there is only appInfo.YYBEnable's value available, and then we change to use appInfo.Android.YYBEnable
			//So we need to use the first value to sync.
			appInfo.Android.YYBEnable = true
		}
		info["YYB_Enable_Android"] = appInfo.Android.YYBEnable
	} else if uaObj.IsIos() {
		info["BundleID"] = appInfo.Ios.BundleID
		info["Scheme"] = appInfo.Ios.Scheme /*"ds7713337217A6E150"*/
		if dlInfo.DownloadUrlIos != "" {
			info["Url"] = dlInfo.DownloadUrlIos
		} else {
			info["Url"] = appInfo.Ios.DownloadUrl /*"itms-apps://itunes.apple.com/artist/seligman-ventures-ltd/id515901779"*/
		}
		info["isUniversallink"] = appInfo.Ios.UniversalLinkEnable
		info["YYB_Enable_Ios_below_9"] = appInfo.Ios.YYBEnableBelow9
		info["YYB_Enable_Ios_Above_9"] = appInfo.Ios.YYBEnableAbove9
		info["Force_download"] = appInfo.Ios.ForceDownload
	}

	//From sender info
	info["Match_id"] = matchID
	info["Download_msg"] = dlInfo.DownloadMsg
	info["Download_btn_text"] = dlInfo.DownloadBtnText
	info["Download_title"] = dlInfo.DownloadTitle
	info["Uninstall_url"] = dlInfo.UninstallUrl
	info["Redirect_url"] = redirectUrl
	info["Timestamp"] = timeStamp()
	info["DsTag"] = dsTag
	info["AppInsStatus"] = appInstallStatus
	return info
}

// 1. get app data
// 2. get html from /template dir
// 3. return html along with other needed information
//use app info to inject the js template, and write the response using the injected js
//if matchId is 0, it means this is a wrong request from a browser which has finished his redirection.
func writeResponse(ctx context.Context, w http.ResponseWriter, r *http.Request, appInfo *appinfo.AppInfo, appID string, matchID string, uaObj *uainfo.UAInfo, dlInfo DownloadInfo, redirectUrl, dsTag string, appInstallStatus int) {
	log.Debugf("Sharelink of the app %s with the matchID %s is clicked", appID, matchID)
	if matchID == "0" {
		log.Info("Sharelink of matchID is 0")
		//TODO close me directly
		fmt.Fprintf(w, "%s", "Succeed!")
		return
	}

	log.Debug("Sharelink access with mobile")
	dstHtml := "sharelink_response_mobile.html"
	info := injectInfo(appInfo, appID, matchID, uaObj, dlInfo, redirectUrl, dsTag, appInstallStatus)

	if r.FormValue("qdsxdebug") == "qdsx" {
		info["debug"] = "true"
	} else {
		info["debug"] = "false"
	}

	log.Debug("Sharelink writeResponse with info", info)
	srcHtmlFile := curDir + jsTemplateDir + dstHtml
	t, err := template.ParseFiles(srcHtmlFile)
	if err != nil {
		log.Error("Sharelink error when parse html files", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	err = t.Execute(w, info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func timeStamp() int64 {
	return int64(time.Now().UTC().Unix())
}
