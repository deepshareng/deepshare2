package inappdata

import (
	"net/http"
	"testing"

	"net/url"

	"log"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/pkg/testutil"
)

func TestGetInAppDataByMatch(t *testing.T) {
	tests := []struct {
		code                              int
		requestUrl                        string
		remoteAddr                        string
		header                            map[string]string
		requestBody                       string
		responseBody                      string
		mockResponseCode                  map[string]int
		mockResponseBody                  map[string]string
		mockMatchRequestMatchPath         string
		mockMatchRequestQueryReceiverInfo string
		mockMatchRequestQueryTracking     string
		mockMatchRequestQueryClientIP     string
		mockMatchRequestQueryClientUA     string
	}{
		//Test a new install of IOS device, without 100% match
		{ // 0
			http.StatusOK,
			`http://fds.so/v2/inappdata/7713337217A6E150`,
			"ip1:port1",
			map[string]string{
				"User-Agent":      "test_useragent1",
				"X-Forwarded-For": "testip1",
			},
			` {
			  				"is_newuser" : true,
			  				"app_version_name" : "2",
			  				"has_telephone" : true,
			  				"sdk_info" : "ios1.1.2",
			  				"brand" : "Apple",
			 				"carrier_name" : "中国联通",
			 				"screen_dpi" : 326,
			  				"os" : "iOS",
			  				"is_emulator" : false,
			 				"bluetooth_version" : "Bluetooth 4.2",
			 				"unique_id" : "320C9F2E-4876-4FD4-A6F5-B3DA804F32C7",
			  				"is_wifi_connected" : true,
			 				"has_nfc" : true,
			  				"screen_height" : 1334,
			  				"app_version_build" : "1.2",
			  				"model" : "iPhone7,2",
			  				"os_version" : "9.1",
			  				"screen_width" : 750,
			  				"hardware_id":"hard1"
							}` + "\n",
			`{"inapp_data":"{\"key1\":\"value1\"}","channels":["1","2"]}` + "\n",
			map[string]int{
				"*": http.StatusOK,
				`/v2/matches/7713337217A6E150`: http.StatusOK,
				`/v2/appinfo/7713337217A6E150`: http.StatusOK,
			},
			map[string]string{
				"*": `` + "\n",
				`/v2/matches/7713337217A6E150`: `{"inapp_data":"{\"key1\":\"value1\"}","channels":["1","2"]}` + "\n",
				`/v2/appinfo/7713337217A6E150`: `{}` + "\n",
			},
			`/v2/matches/7713337217A6E150`,
			`{"unique_id":"320C9F2E-4876-4FD4-A6F5-B3DA804F32C7","app_version_name":"2","app_version_code":0,"app_version_build":"1.2","sdk_info":"ios1.1.2","is_wifi_connected":true,"carrier_name":"中国联通","blueTooth_enable":false,"model":"iPhone7,2","brand":"Apple","is_emulator":false,"os":"iOS","os_version":"9.1","has_nfc":true,"has_telephone":true,"bluetooth_version":"Bluetooth 4.2","screen_dpi":326,"screen_width":750,"screen_height":1334,"uri_scheme":"","hardware_id":"hard1"}`,
			RequestTrackingValueInstall,
			`testip1`,
			`test_useragent1`,
		},
		//Test a open of IOS device, without 100% match
		{ // 1
			http.StatusOK,
			`http://fds.so/v2/inappdata/7713337217A6E150`,
			"ip1:port1",
			map[string]string{
				"User-Agent":      "test_useragent1",
				"X-Forwarded-For": "testip1",
			},
			` {
			 				 "is_newuser" : false,
			 				 "app_version_name" : "2",
			 				 "has_telephone" : true,
			 				 "sdk_info" : "ios1.1.2",
			 				 "brand" : "Apple",
			 				 "carrier_name" : "中国联通",
			 				 "screen_dpi" : 326,
			  				 "os" : "iOS",
			 				 "is_emulator" : false,
			 				 "bluetooth_version" : "Bluetooth 4.2",
			 				 "unique_id" : "320C9F2E-4876-4FD4-A6F5-B3DA804F32C7",
			 				 "is_wifi_connected" : true,
			 				 "click_id" : "l8sT0zzTq",
			 				 "has_nfc" : true,
			 				 "screen_height" : 1334,
			 				 "app_version_build" : "1.2",
							 "model" : "iPhone7,2",
			 				 "os_version" : "9.1",
			  				 "screen_width" : 750,
			  				 "hardware_id":"hard1"
							}` + "\n",
			`{"inapp_data":"{\"key1\":\"value1\"}","channels":["1","2"]}` + "\n",
			map[string]int{
				"*": http.StatusOK,
				"/v2/matches/7713337217A6E150/l8sT0zzTq": http.StatusOK,
				`/v2/appinfo/7713337217A6E150`:           http.StatusOK,
			},
			map[string]string{
				"*": `` + "\n",
				"/v2/matches/7713337217A6E150/l8sT0zzTq": `{"inapp_data":"{\"key1\":\"value1\"}","channels":["1","2"]}` + "\n",
				`/v2/appinfo/7713337217A6E150`:           `{}` + "\n",
			},
			`/v2/matches/7713337217A6E150/l8sT0zzTq`,
			`{"unique_id":"320C9F2E-4876-4FD4-A6F5-B3DA804F32C7","app_version_name":"2","app_version_code":0,"app_version_build":"1.2","sdk_info":"ios1.1.2","is_wifi_connected":true,"carrier_name":"中国联通","blueTooth_enable":false,"model":"iPhone7,2","brand":"Apple","is_emulator":false,"os":"iOS","os_version":"9.1","has_nfc":true,"has_telephone":true,"bluetooth_version":"Bluetooth 4.2","screen_dpi":326,"screen_width":750,"screen_height":1334,"uri_scheme":"","hardware_id":"hard1"}`,
			RequestTrackingValueOpen,
			`testip1`,
			`test_useragent1`,
		},
		//Test if click_id is polluted
		{ // 2
			http.StatusOK,
			`http://fds.so/v2/inappdata/7713337217A6E150`,
			"ip1:port1",
			map[string]string{
				"User-Agent":      "test_useragent1",
				"X-Forwarded-For": "testip1",
			},
			` {
 				 "is_newuser" : false,
 				 "app_version_name" : "2",
 				 "has_telephone" : true,
 				 "sdk_info" : "ios1.1.2",
 				 "brand" : "Apple",
 				 "carrier_name" : "中国联通",
 				 "screen_dpi" : 326,
  				 "os" : "iOS",
 				 "is_emulator" : false,
 				 "bluetooth_version" : "Bluetooth 4.2",
 				 "unique_id" : "320C9F2E-4876-4FD4-A6F5-B3DA804F32C7",
 				 "is_wifi_connected" : true,
 				 "click_id" : "l8sT0zzTq?plg_nld=1&plg_uin=1&plg_auth=1&plg_nld=1&plg_usr=1&plg_vkey=1&plg_dev=1",
 				 "has_nfc" : true,
 				 "screen_height" : 1334,
 				 "app_version_build" : "1.2",
				 "model" : "iPhone7,2",
 				 "os_version" : "9.1",
  				 "screen_width" : 750,
  				 "hardware_id":"hard1"
				}` + "\n",
			`{"inapp_data":"{\"key1\":\"value1\"}","channels":["1","2"]}` + "\n",
			map[string]int{
				"*": http.StatusOK,
				"/v2/matches/7713337217A6E150/l8sT0zzTq": http.StatusOK,
				`/v2/appinfo/7713337217A6E150`:           http.StatusOK,
			},
			map[string]string{
				"*": `` + "\n",
				"/v2/matches/7713337217A6E150/l8sT0zzTq": `{"inapp_data":"{\"key1\":\"value1\"}","channels":["1","2"]}` + "\n",
				`/v2/appinfo/7713337217A6E150`:           `{}` + "\n",
			},
			`/v2/matches/7713337217A6E150/l8sT0zzTq`,
			`{"unique_id":"320C9F2E-4876-4FD4-A6F5-B3DA804F32C7","app_version_name":"2","app_version_code":0,"app_version_build":"1.2","sdk_info":"ios1.1.2","is_wifi_connected":true,"carrier_name":"中国联通","blueTooth_enable":false,"model":"iPhone7,2","brand":"Apple","is_emulator":false,"os":"iOS","os_version":"9.1","has_nfc":true,"has_telephone":true,"bluetooth_version":"Bluetooth 4.2","screen_dpi":326,"screen_width":750,"screen_height":1334,"uri_scheme":"","hardware_id":"hard1"}`,
			RequestTrackingValueOpen,
			`testip1`,
			`test_useragent1`,
		},
	}

	for i, tt := range tests {
		serverMock, clientMock, requestHistory, _ := testutil.MockResponse(tt.mockResponseCode, tt.mockResponseBody)
		handler := newInAppDataHandler(clientMock, serverMock.URL+api.MatchPrefix, serverMock.URL+api.DeviceCookiePrefix, "", serverMock.URL+api.GenerateUrlPrefix, serverMock.URL+api.AppInfoPrefix, nil, api.GetInAppDataPrefix)
		defer serverMock.Close()
		w := testutil.HandleWithRequestInfo(handler, "POST", tt.requestUrl, tt.requestBody, tt.header, tt.remoteAddr)
		if w.Code != tt.code {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.code)
		}

		if tt.code == http.StatusOK {
			resBody := string(w.Body.Bytes())
			if resBody != tt.responseBody {
				t.Errorf("#%d: HTTP response body = %s, want = %s", i, resBody, tt.responseBody)
			}

			//find a request history for match
			reqMatchUrlStr := (*requestHistory)[tt.mockMatchRequestMatchPath]
			log.Printf("#%d: request match url = %s", i, reqMatchUrlStr)
			reqMatchUrl, err := url.Parse(reqMatchUrlStr)
			if err != nil {
				t.Fatalf("Request match url is wrong: %v", err)
			}

			if reqMatchUrl.Path != tt.mockMatchRequestMatchPath {
				t.Errorf("#%d: HTTP Match request match path = %s, want = %s", i, reqMatchUrl.Path, tt.mockMatchRequestMatchPath)
			}
			query := reqMatchUrl.Query()
			receiverInfos, ok := query[RequestReceiverInfoKey]
			if !ok {
				t.Fatal("Request Receiver Info is nil", err)
			}
			if receiverInfos[0] != tt.mockMatchRequestQueryReceiverInfo {
				t.Errorf("#%d: HTTP Match request ReceiverInfo = %s, want = %s", i, receiverInfos[0], tt.mockMatchRequestQueryReceiverInfo)
			}
			trackings, ok := query[RequestTrackingKey]
			if !ok {
				t.Fatal("Request trackings is nil", err)
			}
			if trackings[0] != tt.mockMatchRequestQueryTracking {
				t.Errorf("#%d: HTTP Match request Trackings = %s, want = %s", i, trackings[0], tt.mockMatchRequestQueryTracking)
			}

			clientIps, ok := query[MatchRequestClientIPKey]
			if clientIps[0] != tt.mockMatchRequestQueryClientIP {
				t.Errorf("#%d: HTTP Match request client Ips = %s, want = %s", i, clientIps[0], tt.mockMatchRequestQueryClientIP)
			}

			clientUAs, ok := query[MatchRequestClientUAKey]
			if clientUAs[0] != tt.mockMatchRequestQueryClientUA {
				t.Errorf("#%d: HTTP Match request client UAs = %s, want = %s", i, clientUAs[0], tt.mockMatchRequestQueryClientUA)
			}
		}
	}
}

func TestGetInAppDataByShortSeg(t *testing.T) {
	tests := []struct {
		code                            int
		requestUrl                      string
		remoteAddr                      string
		header                          map[string]string
		requestBody                     string
		responseBody                    string
		mockResponseCode                map[string]int
		mockResponseBody                map[string]string
		mockUrlRequestUrlPath           string
		mockUrlRequestQueryReceiverInfo string
		mockUrlRequestQueryTracking     string
		mockUrlRequestQueryClientIP     string
		mockUrlRequestQueryClientUA     string
	}{
		{
			http.StatusOK,
			`http://fds.so/v2/inappdata/7713337217A6E150`,
			"ip1:port1",
			map[string]string{
				"User-Agent":      "test_useragent1",
				"X-Forwarded-For": "testip1",
			},
			` {
  				"is_newuser" : true,
  				"short_seg" : "aabbcc",
  				"app_version_name" : "2",
  				"has_telephone" : true,
  				"sdk_info" : "ios1.1.2",
  				"brand" : "Apple",
 				"carrier_name" : "中国联通",
 				"screen_dpi" : 326,
  				"os" : "iOS",
  				"is_emulator" : false,
 				"bluetooth_version" : "Bluetooth 4.2",
 				"unique_id" : "320C9F2E-4876-4FD4-A6F5-B3DA804F32C7",
  				"is_wifi_connected" : true,
 				"has_nfc" : true,
  				"screen_height" : 1334,
  				"app_version_build" : "1.2",
  				"model" : "iPhone7,2",
  				"os_version" : "9.1",
  				"screen_width" : 750,
  				"hardware_id":"hard1"
				}` + "\n",
			`{"inapp_data":"{\"key1\":\"value1\"}","channels":["1","2"]}` + "\n",
			map[string]int{
				`/v2/url/7713337217A6E150/aabbcc`: http.StatusOK,
			},
			map[string]string{
				`/v2/url/7713337217A6E150/aabbcc`: `{"inapp_data":"{\"key1\":\"value1\"}","channels":["1","2"]}` + "\n",
			},
			`/v2/url/7713337217A6E150/aabbcc`,
			`{"unique_id":"320C9F2E-4876-4FD4-A6F5-B3DA804F32C7","app_version_name":"2","app_version_code":0,"app_version_build":"1.2","sdk_info":"ios1.1.2","is_wifi_connected":true,"carrier_name":"中国联通","blueTooth_enable":false,"model":"iPhone7,2","brand":"Apple","is_emulator":false,"os":"iOS","os_version":"9.1","has_nfc":true,"has_telephone":true,"bluetooth_version":"Bluetooth 4.2","screen_dpi":326,"screen_width":750,"screen_height":1334,"uri_scheme":"","hardware_id":"hard1"}`,
			RequestTrackingValueInstall,
			`testip1`,
			`test_useragent1`,
		},
	}

	for i, tt := range tests {
		serverMock, clientMock, requestHistory, _ := testutil.MockResponse(tt.mockResponseCode, tt.mockResponseBody)
		handler := newInAppDataHandler(clientMock, serverMock.URL+api.MatchPrefix, serverMock.URL+api.DeviceCookiePrefix, "", serverMock.URL+api.GenerateUrlPrefix, "", nil, api.GetInAppDataPrefix)
		defer serverMock.Close()
		w := testutil.HandleWithRequestInfo(handler, "POST", tt.requestUrl, tt.requestBody, tt.header, tt.remoteAddr)
		if w.Code != tt.code {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.code)
		}

		if tt.code == http.StatusOK {
			resBody := string(w.Body.Bytes())
			if resBody != tt.responseBody {
				t.Errorf("#%d: HTTP response body = %s, want = %s", i, resBody, tt.responseBody)
			}
			reqGenUrlStr := (*requestHistory)[tt.mockUrlRequestUrlPath]
			log.Printf("#%d: request match url = %s", i, reqGenUrlStr)
			reqGenUrl, err := url.Parse(reqGenUrlStr)
			if err != nil {
				t.Fatalf("Request match url is wrong: %v", err)
			}

			if reqGenUrl.Path != tt.mockUrlRequestUrlPath {
				t.Errorf("#%d: HTTP Match request match path = %s, want = %s", i, reqGenUrl.Path, tt.mockUrlRequestUrlPath)
			}
			query := reqGenUrl.Query()
			receiverInfos, ok := query[RequestReceiverInfoKey]
			if !ok {
				t.Fatal("Request Receiver Info is nil", err)
			}
			if receiverInfos[0] != tt.mockUrlRequestQueryReceiverInfo {
				t.Errorf("#%d: HTTP Match request ReceiverInfo = %s, want = %s", i, receiverInfos[0], tt.mockUrlRequestQueryReceiverInfo)
			}
			trackings, ok := query[RequestTrackingKey]
			if !ok {
				t.Fatal("Request trackings is nil", err)
			}
			if trackings[0] != tt.mockUrlRequestQueryTracking {
				t.Errorf("#%d: HTTP Match request Trackings = %s, want = %s", i, trackings[0], tt.mockUrlRequestQueryTracking)
			}
		}
	}
}

func TestShouldUseWechatCookie(t *testing.T) {
	tests := []struct {
		inAppDataPostBody     InAppDataPostBody
		mockResponseBody      map[string]string
		shouldUseWeChatCookie bool
	}{
		{
			inAppDataPostBody: InAppDataPostBody{
				OS:        "iOS",
				OSVersion: "9.1",
			},
			mockResponseBody: map[string]string{
				"*": `{"AppID":"appId","Android":{"YYBEnable":false},"Ios":{"YYBEnableBelow9":true, "YYBEnableAbove9":true}}`,
			},
			shouldUseWeChatCookie: true,
		},
		{
			inAppDataPostBody: InAppDataPostBody{
				OS:        "iOS",
				OSVersion: "9.1",
			},
			mockResponseBody: map[string]string{
				"*": `{"AppID":"appId","Android":{"YYBEnable":false},"Ios":{"YYBEnableBelow9":true, "YYBEnableAbove9":false}}`,
			},
			shouldUseWeChatCookie: false,
		},
		{
			inAppDataPostBody: InAppDataPostBody{
				OS:        "iOS",
				OSVersion: "8.3",
			},
			mockResponseBody: map[string]string{
				"*": `{"AppID":"appId","Android":{"YYBEnable":false},"Ios":{"YYBEnableBelow9":true, "YYBEnableAbove9":false}}`,
			},
			shouldUseWeChatCookie: true,
		},
		{
			inAppDataPostBody: InAppDataPostBody{
				OS:        "iOS",
				OSVersion: "8.3",
			},
			mockResponseBody: map[string]string{
				"*": `{"AppID":"appId","Android":{"YYBEnable":false},"Ios":{"YYBEnableBelow9":false, "YYBEnableAbove9":false}}`,
			},
			shouldUseWeChatCookie: false,
		},
		{
			inAppDataPostBody: InAppDataPostBody{
				OS: "Android",
			},
			mockResponseBody: map[string]string{
				"*": `{"AppID":"appId","Android":{"YYBEnable":false},"Ios":{"YYBEnableBelow9":true, "YYBEnableAbove9":true}}`,
			},
			shouldUseWeChatCookie: false,
		},
		{
			inAppDataPostBody: InAppDataPostBody{
				OS: "Android",
			},
			mockResponseBody: map[string]string{
				"*": `{"AppID":"appId","Android":{"YYBEnable":true},"Ios":{"YYBEnableBelow9":true, "YYBEnableAbove9":true}}`,
			},
			shouldUseWeChatCookie: true,
		},
		{
			inAppDataPostBody: InAppDataPostBody{
				OS: "Android",
			},
			mockResponseBody: map[string]string{
				"*": `{"AppID":"appId","Android":{"YYBEnable":false},"Ios":{"YYBEnableBelow9":true, "YYBEnableAbove9":true},"YYBEnable":true}`,
			},
			shouldUseWeChatCookie: true,
		},
	}
	for i, tt := range tests {
		serverMock, clientMock, _, _ := testutil.MockResponse(map[string]int{
			"*": http.StatusOK,
		}, tt.mockResponseBody)
		handler := newInAppDataHandler(clientMock, serverMock.URL+api.MatchPrefix, serverMock.URL+api.DeviceCookiePrefix, "", serverMock.URL+api.GenerateUrlPrefix, serverMock.URL+api.AppInfoPrefix, nil, api.GetInAppDataPrefix)
		iADH, OK := handler.(*inAppDataHandler)
		if OK {
			shoudlUseWechatCookie := iADH.shouldUseWechatCookie("appId", tt.inAppDataPostBody)
			if tt.shouldUseWeChatCookie != shoudlUseWechatCookie {
				t.Errorf("#%d: Should use wechat cookie = %v, want = %v", i, shoudlUseWechatCookie, tt.shouldUseWeChatCookie)
			}
		} else {
			t.Fatal("inAppDataHandler is in wrong type")
		}

		serverMock.Close()
	}
}
