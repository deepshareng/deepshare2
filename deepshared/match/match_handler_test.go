package match

import (
	"net/http"
	"net/url"
	"path"
	"testing"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"github.com/MISingularity/deepshare2/pkg/testutil"
)

func TestMatchPUT(t *testing.T) {
	tests := []struct {
		cookieID string
		body     string

		wcode int
		wbody string
	}{
		{ //good request
			"mmm1",
			`{"sender_info":{"sender_id":"senderid1","channels":["x","y"]},"inapp_data":"testcontextbytes","client_ip":"testip1","client_ua":"test_useragent1"}` + "\n",
			http.StatusOK,
			"",
		},
		{ //good request
			"mmm2",
			`{"inapp_data":"{ \"key1\":\"test_value1\",\"key2\":2 }","sender_info":{"sender_id":"7E7B2568-B666-4577-A9DE-83A4ED8528B9","channels":[""],"sdk_info":"ios1.1.2"},"client_ip":"127.0.0.1:63518","client_ua":"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B137 Safari/601.1"}` + "\n",
			http.StatusOK,
			"",
		},
		{ //bad header, without client_ip and client_ua in request body
			"mmm3",
			`{"sender_info":{"sender_id":"senderid3","channels":["x","y"]},"inapp_data":"testcontextbytes"}` + "\n",
			http.StatusBadRequest,
			`{"code":302,"message":"Need client_ip and client_ua in put match body"}` + "\n",
		},
		{ //bad header, without client_ua in request body
			// 2015-12-07 modify: client_ua can be empty, to fix the repeat match issue for ios9 universal link
			// when client_ua is empty, should not bind with UA
			"mmm4",
			`{"sender_info":{"sender_id":"senderid4","channels":["x","y"]},"inapp_data":"testcontextbytes","client_ip":"testip4"}` + "\n",
			http.StatusOK,
			"",
		},
		{ //bad header, without client_ip in request body
			"mmm5",
			`{"sender_info":{"sender_id":"senderid5","channels":["x","y"]},"inapp_data":"testcontextbytes","client_ua":"test_useragent5"}` + "\n",
			http.StatusBadRequest,
			`{"code":302,"message":"Need client_ip and client_ua in put match body"}` + "\n",
		},
		{ // bad invaild JSON body
			"mmm6",
			`{"sender_info":{"sender_id":"senderid6","channels":["x","y"]},"inapp_data":"testcontextbytes","client_ip":"testip6","client_ua":"test_useragent6"//` + "\n",
			http.StatusBadRequest,
			`{"code":100,"message":"Body is invalid JSON"}` + "\n",
		},
	}

	handler := newTestMatchHandler()

	for i, tt := range tests {
		url := "http://" + path.Join("example.com", matchPath("AppID", tt.cookieID))
		w := testutil.HandleWithBody(handler, "PUT", url, tt.body)

		if w.Code != tt.wcode {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.wcode)
		}

		resBody := string(w.Body.Bytes())
		if resBody != tt.wbody {
			t.Errorf("#%d: HTTP response body = %s, want = %s", i, resBody, tt.wbody)
		}
	}
}

func TestMatchGETWithUARepeatedly(t *testing.T) {
	handler := newTestMatchHandler()
	//bind first
	testCookieID := "cc1"
	testAppID := "testAppID"
	testReceiverUA := "Mozilla/5.0 (Linux; U; Android 4.1.5; zh-cn; MI 2A Build/JRO03L) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30 MicroMessenger/5.4.0.51_r798589.480 NetType/WIFI"
	testPutBody := `{"sender_info":{"sender_id":"id1","channels":["ch_x","ch_y"]},"inapp_data":"testcontextbytes","client_ip":"testip","client_ua":"` + testReceiverUA + `"}` + "\n"

	urlPut := "http://" + path.Join("example.com", matchPath(testAppID, testCookieID))
	testutil.MustHandleWithBodyOK(handler, "PUT", urlPut, testPutBody)

	//test GET
	tests := []struct {
		ua           string
		receiverInfo string

		wcode int
		wbody string
	}{
		{
			ua:           testReceiverUA,
			receiverInfo: `{"unique_id":"a"}` + "\n",
			wcode:        http.StatusOK,
			wbody:        `{"inapp_data":"testcontextbytes","sender_id":"id1","channels":["ch_x","ch_y"]}` + "\n",
		},
		{
			ua:           testReceiverUA,
			receiverInfo: `{"unique_id":"a"}` + "\n",
			wcode:        http.StatusNotFound,
			wbody:        `{"code":300,"message":"Failed to match the provided information with any existing bindings"}` + "\n",
		},
	}

	for i, tt := range tests {
		uStr := "http://" + path.Join("example.com", matchPath(testAppID, ""))
		u, err := url.Parse(uStr)
		if err != nil {
			t.Fatal(err)
		}
		queries := u.Query()
		queries.Add("client_ip", "testip")
		queries.Add("client_ua", tt.ua)
		queries.Add("receiver_info", tt.receiverInfo)
		u.RawQuery = queries.Encode()
		rUrl := u.String()

		w := testutil.HandleWithBody(handler, "GET", rUrl, tt.wbody)
		if w.Code != tt.wcode {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.wcode)
		}
		if string(w.Body.Bytes()) != tt.wbody {
			t.Errorf("#%d: response body = %s, want = %s", i, w.Body.Bytes(), tt.wbody)
		}
	}
}

func TestMatchGET(t *testing.T) {
	handler := newTestMatchHandler()

	tests := []struct {
		cookieID        string
		ua              string
		body            string
		params          map[string]string
		getWithCookieID bool

		wbody string
		wcode int
	}{
		{ //OK
			cookieID: "m1",
			ua:       "",
			body:     `{"sender_info":{"sender_id":"id1","channels":["ch_x","ch_y"]},"inapp_data":"testcontextbytes","client_ip":"testip1","client_ua":"Mozilla/5.0 (Linux; U; Android 4.1.1; zh-cn; MI 2A Build/JRO03L) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30 MicroMessenger/5.4.0.51_r798589.480 NetType/WIFI"}` + "\n",
			params: map[string]string{
				"receiver_info": `{"unique_id":"ddd1","nfc":true}`,
				"tracking":      "install",
				"client_ip":     "testip1",
				"client_ua":     "Mozilla/5.0 (Linux; U; Android 4.1.1; zh-cn; MI 2A Build/JRO03L) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30 MicroMessenger/5.4.0.51_r798589.480 NetType/WIFI",
			},
			getWithCookieID: false,

			wbody: `{"inapp_data":"testcontextbytes","sender_id":"id1","channels":["ch_x","ch_y"]}` + "\n",
			wcode: http.StatusOK,
		},
		{ //Bad parameters: without client_ua
			cookieID: "m2",
			ua:       "",
			body:     `{"sender_info":{"sender_id":"id2","channels":["ch_x","ch_y"]},"inapp_data":"testcontextbytes","client_ip":"testip2","client_ua":"test_useragent2"}` + "\n",
			params: map[string]string{
				"receiver_info": `{"unique_id":"ddd2","nfc":true}`,
				"tracking":      "install",
				"client_ip":     "testip2",
			},
			getWithCookieID: false,
			wbody:           `{"code":301,"message":"Need client_ip and client_ua in parameters"}` + "\n",
			wcode:           http.StatusBadRequest,
		},
		{ //Bad parameters: without client_ip
			cookieID: "m3",
			ua:       "",
			body:     `{"sender_info":{"sender_id":"id3","channels":["ch_x","ch_y"]},"inapp_data":"testcontextbytes","client_ip":"testip3","client_ua":"test_useragent3"}` + "\n",
			params: map[string]string{
				"receiver_info": `{"unique_id":"ddd3","nfc":true}`,
				"tracking":      "install",
				"client_ua":     "test_useragent3",
			},
			getWithCookieID: false,
			wbody:           `{"code":301,"message":"Need client_ip and client_ua in parameters"}` + "\n",
			wcode:           http.StatusBadRequest,
		},
		{ //Bad parameters: without client_ip and client_ua
			cookieID: "m4",
			ua:       "",
			body:     `{"sender_info":{"sender_id":"id4","channels":["ch_x","ch_y"]},"inapp_data":"testcontextbytes","client_ip":"testip4","client_ua":"test_useragent4"}` + "\n",
			params: map[string]string{
				"receiver_info": `{"unique_id":"ddd4","nfc":true}`,
				"tracking":      "install",
			},
			getWithCookieID: false,
			wbody:           `{"code":301,"message":"Need client_ip and client_ua in parameters"}` + "\n",
			wcode:           http.StatusBadRequest,
		},
		{ //OK: with full receiverInfo which can be parsed to a valid UAInfo
			cookieID: "m5",
			ua:       "",
			body:     `{"sender_info":{"sender_id":"id5","channels":["ch_x","ch_y"]},"inapp_data":"testcontextbytes","client_ip":"testip5","client_ua":"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B137 Safari/601.1"}` + "\n",
			params: map[string]string{
				"receiver_info": `{"unique_id":"ddd5","nfc":true,"os":"iOS","os_version":"9.1"}`,
				"tracking":      "install",
				"client_ip":     "testip5",
				"client_ua":     "testua5",
			},
			getWithCookieID: false,
			wbody:           `{"inapp_data":"testcontextbytes","sender_id":"id5","channels":["ch_x","ch_y"]}` + "\n",
			wcode:           http.StatusOK,
		},
		{ //OK: with cookieID, other params can be empty
			cookieID:        "m6",
			ua:              "",
			body:            `{"sender_info":{"sender_id":"id6","channels":["ch_x","ch_y"]},"inapp_data":"testcontextbytes","client_ip":"testip6","client_ua":"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B137 Safari/601.1"}` + "\n",
			params:          map[string]string{},
			getWithCookieID: true,
			wbody:           `{"inapp_data":"testcontextbytes","sender_id":"id6","channels":["ch_x","ch_y"]}` + "\n",
			wcode:           http.StatusOK,
		},
	}

	// Prepare phase:
	// - do some PUT to bind.
	// - record corresponding cookieID.
	// We can later use the cookieID to test GET.
	for _, p := range tests {
		url := "http://" + path.Join("example.com", matchPath("AppID", p.cookieID))
		testutil.MustHandleWithBodyOK(handler, "PUT", url, p.body)
	}

	for i, tt := range tests {
		cookieID := ""
		if tt.getWithCookieID {
			cookieID = tt.cookieID
		}
		uStr := "http://" + path.Join("example.com", matchPath("AppID", cookieID))
		u, err := url.Parse(uStr)
		if err != nil {
			t.Fatal(err)
		}
		queries := u.Query()
		for k, v := range tt.params {
			queries.Add(k, v)
		}
		u.RawQuery = queries.Encode()
		rUrl := u.String()

		w := testutil.HandleWithBody(handler, "GET", rUrl, tt.wbody)
		if w.Code != tt.wcode {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.wcode)
		}
		if string(w.Body.Bytes()) != tt.wbody {
			t.Errorf("#%d: response body = %s, want = %s", i, w.Body.Bytes(), tt.wbody)
		}
	}
}

// newTestMatchHandler sets up a simple matcher handler with inmem storage.
func newTestMatchHandler() http.Handler {
	s := NewSimpleMatcher(storage.NewInMemSimpleKV(), uaMatchExpireAfterSecDefault)
	return newMatchHandler(s, messaging.NewSimpleProducer(nil), api.MatchPrefix)
}

func matchPath(appID, cookieID string) string {
	if cookieID != "" {
		return path.Join(api.MatchPrefix, appID, cookieID)
	}
	return path.Join(api.MatchPrefix, appID)
}
