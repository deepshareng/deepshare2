package testapi

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"reflect"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/appinfo"
	"github.com/MISingularity/deepshare2/frontend/inappdata"
	"github.com/MISingularity/deepshare2/frontend/jsapi"
	"github.com/MISingularity/deepshare2/frontend/urlgenerator"
	"github.com/MISingularity/deepshare2/pkg/testutil"
)

const (
	test100MatchForIOS9 = true

	testUniqueID  = "testdevice1"
	testSenderID  = "test_sender"
	testSenderID2 = "test_sender_2"

	appIDProduction      = "1652E90881C1FAE8"
	shortAppIDProduction = "000"
	appIDStaging         = "9aeea8a7f4c56ff4"
	shortAppIDStaging    = "0000"
	testParam            = "{ \"key1\":\"test_value1\",\"key2\":\"test_value2\",\"source\":\"ios\" }"
	uaIOS                = "Mozilla/5.0 (iPhone; CPU iPhone OS 7_0_3 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) Version/7.0 Mobile/11B511 Safari/9537.53"
	uaIOS9               = "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) Version/7.0 Mobile/11B511 Safari/9537.53"
	uaAndroid            = "Mozilla/5.0 (Linux; U; Android 4.1.1; zh-cn; MI 2A Build/JRO03L) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30 MicroMessenger/5.4.0.51_r798589.480 NetType/WIFI"

	waitDurationForBind = time.Millisecond * time.Duration(200)
)

var (
	serverUrl              string
	genurlAddr             string
	inappDataAddr          string
	binddevicetocookieAddr string
	sharelinkAddr          string
	dsusageAddr            string
	appinfoAddr            string
)

var (
	testRegisterApp                            = false
	testAppID                                  = appIDProduction
	testShortAppID                             = shortAppIDProduction
	testChannels                               = []string{"1", "2", "3"}
	testChannels2                              = []string{"a", "b", "c"}
	metrics         map[string][]time.Duration = make(map[string][]time.Duration)
	clickUrlTests                              = map[string]struct {
		method string
		path   string
		ua     string
	}{
		"click_url_ios": {
			"GET",
			sharelinkAddr + "/d/<appID>/<clickID>", //fulfill after generate url
			uaIOS,
		},
		"click_url_ios9": {
			"GET",
			sharelinkAddr + "/d/<appID>/<clickID>", //fulfill after generate url
			uaIOS9,
		},
	}
	testParamJSApi1 = map[string]interface{}{"name": "name1"}
	testParamJSApi2 = map[string]interface{}{"name": "name2"}
)

var (
	timeStart    int64
	accessedPath = make(map[string]int)
)

func appendTag(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	q := u.Query()
	q.Add("tag", strconv.Itoa(int(timeStart)))
	if v, ok := accessedPath[u.Path]; ok && v > 0 {
		q.Add("sequence", strconv.Itoa(v))
	}
	u.RawQuery = q.Encode()

	accessedPath[u.Path] = accessedPath[u.Path] + 1
	return u.String()
}

func addMetrics(key string, lapse time.Duration) {
	metrics[key] = append(metrics[key], lapse)
}

func TestMain(m *testing.M) {
	//parse os args
	filename := flag.String("logfile", "", "Specify the logfile name to write logs to")
	env := flag.String("env", "production", "Specify the server env to test (staging/production)")
	serverAddr := flag.String("server-addr", "https://fds.so", "Specify the server addr to test")
	argGenurlAddr := flag.String("genurl-addr", "", "Specify the internal addr of genurl")
	argInappDataAddr := flag.String("inappdata-addr", "", "Specify the internal addr of genurl")
	argBinddevicetocookieAddr := flag.String("binddevicetocookie-addr", "", "Specify the internal addr of genurl")
	argSharelinkAddr := flag.String("sharelink-addr", "", "Specify the internal addr of genurl")
	argDsusageAddr := flag.String("dsusage-addr", "", "Specify the internal addr of genurl")
	argAppinfoAddr := flag.String("appinfo-addr", "", "Specify the internal addr of genurl")
	registerApp := flag.Bool("register-app", false, "Register app before test, useful only on local env")
	flag.Parse()

	serverUrl = *serverAddr
	if serverUrl == "" {
		log.Fatal("server addr could not be empty")
	}
	genurlAddr = serverUrl
	inappDataAddr = serverUrl
	binddevicetocookieAddr = serverUrl
	sharelinkAddr = serverUrl
	dsusageAddr = serverUrl
	appinfoAddr = serverUrl
	if *argGenurlAddr != "" {
		genurlAddr = *argGenurlAddr
	}
	if *argInappDataAddr != "" {
		inappDataAddr = *argInappDataAddr
	}
	if *argBinddevicetocookieAddr != "" {
		binddevicetocookieAddr = *argBinddevicetocookieAddr
	}
	if *argSharelinkAddr != "" {
		sharelinkAddr = *argSharelinkAddr
	}
	if *argDsusageAddr != "" {
		dsusageAddr = *argDsusageAddr
	}
	if *argAppinfoAddr != "" {
		appinfoAddr = *argAppinfoAddr
	}
	testRegisterApp = *registerApp

	switch *filename {
	case "":
		log.SetOutput(bytes.NewBuffer(nil))
	case "stdout":
		log.SetOutput(os.Stdout)
	default:
		os.Remove(*filename)
		f, err := os.OpenFile(*filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal("Failed to create log file, err:", err)
		}
		log.SetOutput(f)
	}

	log.Println("Testing against server:", *serverAddr)
	log.Println("genurl addr:", genurlAddr)
	log.Println("inappdata addr:", inappDataAddr)
	log.Println("binddevicetocookie addr:", binddevicetocookieAddr)
	log.Println("sharelink addr:", sharelinkAddr)
	log.Println("dsusage addr:", dsusageAddr)
	log.Println("appinfo addr:", appinfoAddr)

	switch *env {
	case "production":
		testAppID = appIDProduction
		testShortAppID = shortAppIDProduction
	case "staging":
		testAppID = appIDStaging
		testShortAppID = shortAppIDStaging
	default:
		fmt.Println("Invalid argument \"env\": \n\tshould be \"staging\" or \"production\"")
		os.Exit(1)
	}

	timeStart = time.Now().Unix()
	retCode := m.Run()

	log.Println("Lapses in milli-seconds:")
	for k, v := range metrics {
		log.Println("	", k, ":", v)
	}

	os.Exit(retCode)
}
func TestRegisterApp(t *testing.T) {
	//register app
	if testRegisterApp {
		key := "register_app"
		log.Println("---.", key)
		tt := struct {
			method string
			path   string
			ua     string
		}{
			"PUT",
			appinfoAddr + "/v2/appinfo/" + testAppID,
			"",
		}
		req := appinfo.AppInfo{
			AppID:   testAppID,
			ShortID: "000",
			Android: appinfo.AppAndroidInfo{
				Scheme: "deepshare",
				Host:   "com.singulariti.testtest",
				Pkg:    "com.singulariti.testtest",
			},
			Ios: appinfo.AppIosInfo{
				Scheme: "ds" + testAppID,
			},
		}
		lapse := testRequestServer(t, tt.method, appendTag(tt.path), "", "", req, nil)
		addMetrics(key, lapse)
	}
}

func TestRobot_UniversalLink(t *testing.T) {

	tests := map[string]struct {
		method       string
		path         string
		ua           string
		receiverInfo api.MatchReceiverInfo
	}{
		"register_app": {
			"PUT",
			appinfoAddr + "/v2/appinfo/" + testAppID,
			"",
			api.MatchReceiverInfo{},
		},
		"bind_device_to_cookie": {
			"GET",
			binddevicetocookieAddr + "/v2/binddevicetocookie/" + testUniqueID,
			uaIOS9,
			api.MatchReceiverInfo{},
		},
		"inapp_data_100_match_install": {
			"POST",
			inappDataAddr + "/v2/inappdata/" + testAppID,
			"",
			api.MatchReceiverInfo{
				OS:        "iOS",
				OSVersion: "9.1",
			},
		},
		"inapp_data_by_shortseg_open": {
			"POST",
			inappDataAddr + "/v2/inappdata/" + testAppID,
			uaIOS9,
			api.MatchReceiverInfo{
				OS:        "iOS",
				OSVersion: "9.1",
			},
		},
		"inapp_data_by_polluted_shortseg_open": {
			"POST",
			inappDataAddr + "/v2/inappdata/" + testAppID,
			uaIOS9,
			api.MatchReceiverInfo{
				OS:        "iOS",
				OSVersion: "9.1",
			},
		},
	}

	//generate url
	sharelink, shortSeg := generateUrl(t, uaAndroid)
	//click the url
	cookieID := clickUrlWithDSCookie("click_url_ios9", sharelink, "ccccc_ios9", t)
	log.Println("device1 click_url, got cookieID:", cookieID)
	//bind device to cookie
	if test100MatchForIOS9 {
		key := "bind_device_to_cookie"
		log.Println("---.", key)
		tt := tests[key]
		_, _, lapse := testutil.GetHttpWithDSCookie(t, appendTag(tt.path), uaIOS9, cookieID)
		addMetrics(key, lapse)
	}

	//sleep to wait for the bind to finish
	time.Sleep(waitDurationForBind)

	// get in-app data by ios9 100% match(install)
	if test100MatchForIOS9 {
		key := "inapp_data_100_match_install"
		log.Println("---.", key)
		tt := tests[key]
		req := &inappdata.InAppDataPostBody{
			UniqueID:  testUniqueID,
			IsNewUser: true,
			OS:        tt.receiverInfo.OS,
			OSVersion: tt.receiverInfo.OSVersion,
		}
		resp := &inappdata.InAppDataResponseBody{}
		addMetrics(key, testRequestServer(t, tt.method, appendTag(tt.path), tt.ua, "", req, resp))
		if resp.InAppData != testParam {
			log.Printf("[ERROR] get in-app data 100 match (install) failed, inappdata = %s, want = %s\n", resp.InAppData, testParam)
			t.Fatalf("get in-app data 100 match (install) failed, inappdata = %s, want = %s\n", resp.InAppData, testParam)
		}
		if !reflect.DeepEqual(resp.Channels, testChannels) {
			log.Printf("[ERROR] get in-app data 100 match (install) failed, channels = %s, want = %s\n", resp.Channels, testChannels)
			t.Fatalf("get in-app data 100 match (install) failed, channels = %s, want = %s\n", resp.Channels, testChannels)
		}
	}

	//get in-app data by shortseg
	{
		key := "inapp_data_by_shortseg_open"
		log.Println("---.", key)
		tt := tests[key]
		req := &inappdata.InAppDataPostBody{
			UniqueID:  "testdevice1",
			IsNewUser: false,
			ShortSeg:  shortSeg,
			OS:        tt.receiverInfo.OS,
			OSVersion: tt.receiverInfo.OSVersion,
		}
		resp := &inappdata.InAppDataResponseBody{}
		addMetrics(key, testRequestServer(t, tt.method, appendTag(tt.path), tt.ua, "", req, resp))
		if resp.InAppData != testParam {
			log.Printf("[ERROR] get in-app data by shortseg(open) failed, inappdata = %s, want = %s\n", resp.InAppData, testParam)
			t.Fatalf("get in-app data by shortseg(open) failed, inappdata = %s, want = %s\n", resp.InAppData, testParam)
		}
		if !reflect.DeepEqual(resp.Channels, testChannels) {
			log.Printf("[ERROR] get in-app data by shortseg(open) failed, channels = %s, want = %s\n", resp.Channels, testChannels)
			t.Fatalf("get in-app data by shortseg(open) failed, channels = %s, want = %s\n", resp.Channels, testChannels)
		}
	}

	//get in-app data by polluted shortseg
	{
		key := "inapp_data_by_polluted_shortseg_open"
		log.Println("---.", key)
		tt := tests[key]
		req := &inappdata.InAppDataPostBody{
			UniqueID:  "testdevice1",
			IsNewUser: false,
			ShortSeg:  shortSeg + "?plg_nld=1&plg_uin=1&plg_auth=1&plg_nld=1&plg_usr=1&plg_vkey=1&plg_dev=1",
			OS:        tt.receiverInfo.OS,
			OSVersion: tt.receiverInfo.OSVersion,
		}
		resp := &inappdata.InAppDataResponseBody{}
		addMetrics(key, testRequestServer(t, tt.method, appendTag(tt.path), tt.ua, "", req, resp))
		if resp.InAppData != testParam {
			log.Printf("[ERROR] get in-app data by shortseg(open) failed, inappdata = %s, want = %s\n", resp.InAppData, testParam)
			t.Fatalf("get in-app data by shortseg(open) failed, inappdata = %s, want = %s\n", resp.InAppData, testParam)
		}
		if !reflect.DeepEqual(resp.Channels, testChannels) {
			log.Printf("[ERROR] get in-app data by shortseg(open) failed, channels = %s, want = %s\n", resp.Channels, testChannels)
			t.Fatalf("get in-app data by shortseg(open) failed, channels = %s, want = %s\n", resp.Channels, testChannels)
		}
	}

}

func TestRobot(t *testing.T) {
	tests := map[string]struct {
		method       string
		path         string
		ua           string
		receiverInfo api.MatchReceiverInfo
	}{
		"inapp_data_by_ua_open": {
			"POST",
			inappDataAddr + "/v2/inappdata/" + testAppID,
			uaIOS,
			api.MatchReceiverInfo{
				OS:        "iOS",
				OSVersion: "7.0.3",
			},
		},
		"inapp_data_by_clickid_open": {
			"POST",
			inappDataAddr + "/v2/inappdata/" + testAppID,
			"",
			api.MatchReceiverInfo{
				OS:        "iOS",
				OSVersion: "7.0.3",
			},
		},
	}
	sharelink, _ := generateUrl(t, uaAndroid)
	cookieID := clickUrlWithDSCookie("click_url_ios", sharelink, "ccccc_ios7", t)

	//sleep to wait for the bind to finish
	time.Sleep(waitDurationForBind)

	//get in-app data by cookieid(clickid) (open)
	{
		key := "inapp_data_by_clickid_open"
		log.Println("---.", key)
		tt := tests[key]
		req := &inappdata.InAppDataPostBody{
			UniqueID:  "testdevice2",
			IsNewUser: false,
			ClickID:   cookieID,
			OS:        tt.receiverInfo.OS,
			OSVersion: tt.receiverInfo.OSVersion,
		}
		resp := &inappdata.InAppDataResponseBody{}
		addMetrics(key, testRequestServer(t, tt.method, appendTag(tt.path), tt.ua, "", req, resp))
		if resp.InAppData != testParam {
			log.Printf("[ERROR] get in-app data by clickid(open) failed, inappdata = %s, want = %s\n", resp.InAppData, testParam)
			t.Fatalf("get in-app data by clickid(open) failed, inappdata = %s, want = %s\n", resp.InAppData, testParam)
		}
		if !reflect.DeepEqual(resp.Channels, testChannels) {
			log.Printf("[ERROR] get in-app data by clickid(open) failed, channels = %v, want = %v\n", resp.Channels, testChannels)
			t.Fatalf("get in-app data by clickid(open) failed, channels = %s, want = %s\n", resp.Channels, testChannels)
		}
	}
	//get in-app data by ua (open)
	clickUrlWithDSCookie("click_url_ios", sharelink, "ccccc_android", t)
	time.Sleep(waitDurationForBind)
	{
		key := "inapp_data_by_ua_open"
		log.Println("---.", key)
		log.Println("	get in-app data by ua - first time")
		tt := tests[key]
		req := &inappdata.InAppDataPostBody{
			UniqueID:  "testdevice2",
			IsNewUser: false,
			ClickID:   "dl_no_value",
			OS:        tt.receiverInfo.OS,
			OSVersion: tt.receiverInfo.OSVersion,
		}
		resp := &inappdata.InAppDataResponseBody{}
		addMetrics(key, testRequestServer(t, tt.method, appendTag(tt.path), tt.ua, "", req, resp))
		if resp.InAppData != testParam {
			log.Printf("[ERROR] get in-app data by ua(open) failed, inappdata = %s, want = %s\n", resp.InAppData, testParam)
			t.Fatalf("get in-app data by ua(open) failed, inappdata = %s, want = %s\n", resp.InAppData, testParam)
		}
		if !reflect.DeepEqual(resp.Channels, testChannels) {
			log.Printf("[ERROR] get in-app data by ua(open) failed, channels = %s, want = %s\n", resp.Channels, testChannels)
			t.Fatalf("get in-app data by ua(open) failed, channels = %s, want = %s\n", resp.Channels, testChannels)
		}
		time.Sleep(waitDurationForBind)
		log.Println("	get in-app data by ua - repeatedly")
		req = &inappdata.InAppDataPostBody{
			UniqueID:  "testdevice2",
			IsNewUser: false,
			ClickID:   "dl_no_value",
			OS:        tt.receiverInfo.OS,
			OSVersion: tt.receiverInfo.OSVersion,
		}
		resp = &inappdata.InAppDataResponseBody{}
		addMetrics(key, testRequestServer(t, tt.method, appendTag(tt.path), tt.ua, "", req, resp))
		if resp.InAppData != "" {
			log.Printf("[ERROR] get in-app data by ua repeatedly(open) failed, inappdata = %s, want = %s\n", resp.InAppData, "")
			t.Fatalf("get in-app data by ua repeatedly(open) failed, inappdata = %s, want = %s\n", resp.InAppData, "")
		}
	}
}

func TestDSUsage(t *testing.T) {
	//wait 1 s for the attribution aggregation to finish
	time.Sleep(time.Second * time.Duration(4))
	urlStr := dsusageAddr + "/v2/dsusages/" + testAppID + "/" + testSenderID
	resp := make(map[string]int)
	addMetrics("get_dsusage", testRequestServer(t, "GET", appendTag(urlStr), "", "", nil, &resp))

	//clear usage after get (no matter if the get result is correct or not, to make sure the next test will not be affected)
	lapse := testRequestServer(t, "DELETE", appendTag(urlStr), "", "", nil, &resp)
	addMetrics("delete_dsusage", lapse)

	//The robots opened 4 times and installed once
	if resp["new_install"] != 1 && resp["new_open"] != 4 {
		log.Printf("[ERROR] installs = %d, want = %d, opens = %d, want = %d\n", resp["new_install"], 1, resp["new_open"], 4)
		t.Fatalf("installs = %d, want = %d, opens = %d, want = %d\n", resp["new_install"], 1, resp["new_open"], 4)
	}

	//after clear the usages should be 0
	addMetrics("get_dsusage", testRequestServer(t, "GET", appendTag(urlStr), "", "", nil, &resp))
	//The robots opened 4 times and installed once
	if resp["new_install"] != 0 && resp["new_open"] != 0 {
		log.Printf("[ERROR] installs = %d, want = %d, opens = %d, want = %d\n", resp["new_install"], 0, resp["new_open"], 0)
		t.Fatalf("installs = %d, want = %d, opens = %d, want = %d\n", resp["new_install"], 0, resp["new_open"], 0)
	}
}

func testRequestServer(t *testing.T, method string, path string, ua string, dscookie string, req interface{}, resp interface{}) time.Duration {
	var lapse time.Duration
	u, err := url.Parse(path)
	if err != nil {
		log.Println(err)
		t.Fatal(err)
	}
	switch method {
	case "GET":
		if n, err := testutil.GetJsonResponse(u.String(), ua, resp); err != nil {
			log.Println(err)
			t.Fatal(err)
		} else {
			lapse = n
		}
	case "POST":
		if n, err := testutil.PostJson(u.String(), req, ua, dscookie, &resp); err != nil {
			log.Println(err)
			t.Fatal(err)
		} else {
			lapse = n
		}
	case "PUT":
		if n, err := testutil.PutJson(u.String(), req, ua, dscookie, &resp); err != nil {
			log.Println(err)
			t.Fatal(err)
		} else {
			lapse = n
		}
	case "DELETE":
		lapse = testutil.DeleteRequest(path)
	}
	return lapse
}

func clickUrlWithDSCookie(key string, sharelink string, dscookie string, t *testing.T) string {
	sharelink = strings.Replace(sharelink, serverUrl, sharelinkAddr, -1)
	log.Println("---.", key)
	tt := clickUrlTests[key]
	tt.path = sharelink
	header, body, lapse := testutil.GetHttpWithDSCookie(t, appendTag(tt.path), tt.ua, dscookie)
	addMetrics(key, lapse)
	cookie := header.Get("Set-Cookie")
	parts := strings.Split(cookie, ";")
	cookieID := ""
	for _, seg := range parts {
		a := strings.Split(seg, "=")
		if len(a) >= 2 && a[0] == "dscookie" {
			cookieID = a[1]
			break
		}
	}

	if len(body) < 100 {
		log.Println("click url failed, html body is too short")
		t.Fatal("click url failed, html body is too short")
	}
	if !strings.Contains(strings.ToLower(string(body)), "<html>") {
		log.Println("click url failed, should contain <html>")
		t.Fatal("click url failed, should contain <html>")
	}
	if !strings.Contains(string(body), "deepshare-redirect.min.js") {
		log.Println("click url faild, should contain deepshare-redirect.min.js")
		t.Fatal("click url faild, should contain deepshare-redirect.min.js")
	}
	return cookieID
}

func generateUrl(t *testing.T, ua string) (string, string) {
	key := "generate_url"
	log.Println("---.", key)
	req := &urlgenerator.GenURLPostBody{
		InAppDataReq:  testParam,
		DownloadTitle: "dTitle",
		DownloadMsg:   "dMsg",
		RedirectUrl:   "rUrl",
		IsShort:       true,
		SDKInfo:       "testsdk0.0",
		SenderID:      testSenderID,
		Channels:      testChannels,
	}
	resp := &urlgenerator.GenURLResponseBody{}
	addMetrics(key, testRequestServer(t, "POST", appendTag(genurlAddr+"/v2/url/"+testAppID), ua, "", req, resp))
	if resp.Url == "" {
		log.Println("generate url failed: url is empty")
		t.Fatal("generate url failed: url is empty")
	}
	if !strings.HasPrefix(resp.Url, serverUrl) {
		log.Println("[ERROR] generate url failed, invalid host")
		t.Fatal("generate url failed, invalid host")
	}
	parts := strings.Split(strings.TrimPrefix(resp.Url, serverUrl), "/")
	if len(parts) < 2 {
		log.Println("[ERROR] generate url failed, url format is not valid:", resp.Url)
		t.Fatal("generate url failed, url format is not valid:", resp.Url)
	}
	appID := parts[len(parts)-2]
	if appID != testAppID && appID != testShortAppID {
		log.Printf("[ERROR] generate url failed, appid not agree, want=%s got=%s", testAppID, appID)
		t.Fatalf("generate url failed, appid not agree, want=%s got=%s", testAppID, appID)
	}
	shortSeg := parts[len(parts)-1]
	if shortSeg == "" {
		log.Println("[ERROR] generate url failed, shortseg should not be empty")
		t.Fatal("generate url failed, shortseg should not be empty")
	}
	return resp.Url, shortSeg
}

func TestRobotJSApi(t *testing.T) {
	tests := map[string]struct {
		method       string
		path         string
		ua           string
		receiverInfo api.MatchReceiverInfo
	}{
		"js_api_1": {
			"POST",
			serverUrl + "/v2/jsapi/" + testAppID,
			uaIOS9,
			api.MatchReceiverInfo{
				OS:        "iOS",
				OSVersion: "9.1",
			},
		},
		"js_api_2": {
			"POST",
			serverUrl + "/v2/jsapi/" + testAppID + "?clicked=true",
			uaIOS9,
			api.MatchReceiverInfo{
				OS:        "iOS",
				OSVersion: "9.1",
			},
		},
		"inapp_data_by_cookie_deeplinkid_jsapi": {
			"POST",
			serverUrl + "/v2/inappdata/" + testAppID,
			"",
			api.MatchReceiverInfo{
				OS:        "iOS",
				OSVersion: "9.1",
			},
		},
		"inapp_data_by_ua_jsapi": {
			"POST",
			serverUrl + "/v2/inappdata/" + testAppID,
			uaIOS9,
			api.MatchReceiverInfo{
				OS:        "iOS",
				OSVersion: "9.1",
			},
		},
	}

	//POST to JSApi, first call
	cookieID := func() string {
		key := "js_api_1"
		tt := tests[key]
		cookieID := postJSApiWithDSCookie(key, appendTag(tt.path), tt.ua, "ccccc_jsapi", t)
		return cookieID
	}()

	//sleep to wait for the bind to finish
	time.Sleep(waitDurationForBind)

	//get in-app data by cookieID(clickID) + deeplinkid
	{
		key := "inapp_data_by_cookie_deeplinkid_jsapi"
		tt := tests[key]
		log.Println("---.", key, "deeplink_id = 1")
		req := &inappdata.InAppDataPostBody{
			UniqueID:   "testdevice_jsapi",
			IsNewUser:  false,
			ClickID:    cookieID,
			DeeplinkID: "1",
			OS:         tt.receiverInfo.OS,
			OSVersion:  tt.receiverInfo.OSVersion,
		}
		resp := &inappdata.InAppDataResponseBody{}
		addMetrics(key, testRequestServer(t, tt.method, appendTag(tt.path), tt.ua, "", req, resp))
		inappData, _ := json.Marshal(resp.InAppData)
		s, _ := strconv.Unquote(string(inappData))
		wInAppdata, _ := json.Marshal(testParamJSApi1)
		if s != string(wInAppdata) {
			log.Printf("[ERROR] %s failed, inappdata = %v, want = %v\n", key, s, string(wInAppdata))
			t.Fatalf("%s failed, inappdata = %v, want = %v\n", key, s, string(wInAppdata))
		}
		if !reflect.DeepEqual(resp.Channels, testChannels) {
			log.Printf("[ERROR] %s failed, channels = %v, want = %v\n", key, resp.Channels, testChannels)
			t.Fatalf("%s failed, channels = %s, want = %s\n", key, resp.Channels, testChannels)
		}
	}

	//get in-app data by cookieID(clickID) + deeplinkid
	{
		key := "inapp_data_by_cookie_deeplinkid_jsapi"
		tt := tests[key]
		log.Println("---.", key, "deeplink_id = 2")
		req := &inappdata.InAppDataPostBody{
			UniqueID:   "testdevice_jsapi",
			IsNewUser:  false,
			ClickID:    cookieID,
			DeeplinkID: "2",
			OS:         tt.receiverInfo.OS,
			OSVersion:  tt.receiverInfo.OSVersion,
		}
		resp := &inappdata.InAppDataResponseBody{}
		addMetrics(key, testRequestServer(t, tt.method, appendTag(tt.path), tt.ua, "", req, resp))
		inappData, _ := json.Marshal(resp.InAppData)
		s, _ := strconv.Unquote(string(inappData))
		wInAppdata, _ := json.Marshal(testParamJSApi2)
		if s != string(wInAppdata) {
			log.Printf("[ERROR] %s failed, inappdata = %v, want = %v\n", key, s, string(wInAppdata))
			t.Fatalf("%s failed, inappdata = %v, want = %v\n", key, s, string(wInAppdata))
		}
		if !reflect.DeepEqual(resp.Channels, testChannels2) {
			log.Printf("[ERROR] %s failed, channels = %v, want = %v\n", key, resp.Channels, testChannels2)
			t.Fatalf("%s failed, channels = %s, want = %s\n", key, resp.Channels, testChannels2)
		}
	}

	//POST to jsapi, second call
	{
		key := "js_api_2"
		tt := tests[key]
		log.Println("---.", key)
		req := jsapi.JsApiPostClicked{
			InAppDataReq: testParamJSApi2,
			SenderID:     testSenderID2,
			Channels:     testChannels2,
		}
		resp := &jsapi.JsApiRespClicked{}
		addMetrics(key, testRequestServer(t, "POST", appendTag(tt.path), tt.ua, "", req, resp))
		if !resp.OK {
			log.Printf("[ERROR] %s failed, resp not OK: %v\n", key, resp.OK)
			t.Fatalf("%s failed, resp not OK: %v\n", key, resp.OK)
		}
	}

	//sleep to wait for the bind to finish
	time.Sleep(waitDurationForBind)

	//get in-app data by UA
	{
		key := "inapp_data_by_ua_jsapi"
		tt := tests[key]
		log.Println("---.", key)
		req := &inappdata.InAppDataPostBody{
			UniqueID:  "testdevice_jsapi",
			IsNewUser: false,
			OS:        tt.receiverInfo.OS,
			OSVersion: tt.receiverInfo.OSVersion,
		}
		resp := &inappdata.InAppDataResponseBody{}
		addMetrics(key, testRequestServer(t, tt.method, appendTag(tt.path), tt.ua, "", req, resp))
		inappData, _ := json.Marshal(resp.InAppData)
		s, _ := strconv.Unquote(string(inappData))
		wInAppdata, _ := json.Marshal(testParamJSApi2)
		if s != string(wInAppdata) {
			log.Printf("[ERROR] %s failed, inappdata = %v, want = %v\n", key, s, string(wInAppdata))
			t.Fatalf("%s failed, inappdata = %v, want = %v\n", key, s, string(wInAppdata))
		}
		if !reflect.DeepEqual(resp.Channels, testChannels2) {
			log.Printf("[ERROR] %s failed, channels = %v, want = %v\n", key, resp.Channels, testChannels2)
			t.Fatalf("%s failed, channels = %s, want = %s\n", key, resp.Channels, testChannels2)
		}
	}

}

func postJSApiWithDSCookie(key string, path string, ua string, dscookie string, t *testing.T) (cookie string) {
	log.Println("---.", key)
	req := &jsapi.JsApiPost{
		{
			DeepLinkID:   "1",
			InAppDataReq: testParamJSApi1,
			SenderID:     testSenderID,
			Channels:     testChannels,
		},
		{
			DeepLinkID:   "2",
			InAppDataReq: testParamJSApi2,
			SenderID:     testSenderID2,
			Channels:     testChannels2,
		},
	}
	resp := &jsapi.JsApiResp{}
	addMetrics(key, testRequestServer(t, "POST", path, ua, "cccc_jsapi", req, resp))
	if resp.AppID != testAppID && resp.AppID != testShortAppID {
		log.Printf("[ERROR] %s failed, appID=%s, want=%s\n", key, resp.AppID, testAppID)
		t.Fatalf("%s failed, appID=%s, want=%s\n", key, resp.AppID, testAppID)
	}
	if resp.DSUrls == nil || resp.DSUrls["1"] == "" || resp.DSUrls["2"] == "" {
		log.Printf("[ERROR] %s failed, ds_urls=%v, want len=2\n", key, resp.DSUrls)
		t.Fatalf("%s failed, ds_urls=%v, want len=2\n", key, resp.DSUrls)
	}
	for k, dsurl := range resp.DSUrls {
		parts := strings.Split(strings.TrimPrefix(dsurl, serverUrl), "/")
		if len(parts) < 2 {
			log.Printf("[ERROR] %s failed, url format is not valid: %s", key, dsurl)
			t.Fatalf("%s failed, url format is not valid: %s", key, dsurl)
		}
		appID := parts[len(parts)-2]
		if appID != testAppID && appID != testShortAppID {
			log.Printf("[ERROR] %s failed, appid not agree, want=%s got=%s", key, testAppID, appID)
			t.Fatalf("%s failed, appid not agree, want=%s got=%s", key, testAppID, appID)
		}
		if k != "1" && k != "2" {
			log.Printf("[ERROR] %s failed, unknown deeplink_id=%s", key, k)
			t.Fatalf("%s failed, unknown deeplink_id=%s", key, k)
		}
		shortSeg := parts[len(parts)-1]
		if shortSeg == "" {
			log.Println("[ERROR] generate url failed, shortseg should not be empty")
			t.Fatal("generate url failed, shortseg should not be empty")
		}
		if resp.MatchId == "" {
			log.Println("[ERROR] generate url failed, matchID should not be empty")
			t.Fatal("generate url failed, matchID should not be empty")
		}
	}
	return resp.MatchId
}
