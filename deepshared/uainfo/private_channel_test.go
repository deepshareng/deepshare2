package uainfo

import (
	"reflect"
	"strings"
	"testing"

	"github.com/MISingularity/deepshare2/pkg/path"
	"github.com/MISingularity/uap-go/uaparser"
)

var uas = []struct {
	browser  string
	version  string
	params   map[string]bool
	uastring string
}{
	{
		browser:  "Chrome",
		version:  "33",
		params:   map[string]bool{},
		uastring: "Mozilla/5.0 (Linux; Android 4.4.4; Android SDK built for x86 Build/KK) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/33.0.0.0 Mobile Safari/537.36",
	},
	{
		browser:  "Chrome",
		version:  "40",
		params:   map[string]bool{},
		uastring: "Mozilla/5.0 (iPhone; CPU iPhone OS 8_1_3 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) CriOS/40.0.2214.69 Mobile/12B466 Safari/600.1.4",
	},
	{
		params:   map[string]bool{"wechat": true},
		uastring: "Mozilla/5.0 (Linux; U; Android 4.1.1; zh-cn; MI 2A Build/JRO03L) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30 MicroMessenger/5.4.0.51_r798589.480 NetType/WIFI",
	},
	{
		browser:  "Chrome",
		version:  "45",
		params:   map[string]bool{},
		uastring: "Mozilla/5.0 (iPhone; CPU iPhone OS 8_4_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) CriOS/45.0.2454.89 Mobile/12H321 Safari/600.1.4",
	},
	{
		browser:  "Safari",
		params:   map[string]bool{"wechat": true},
		uastring: "Mozilla/5.0 (iPhone; CPU iPhone OS 8_4_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Mobile/12H321 MicroMessenger/6.2.6 NetType/WIFI Language/zh_CN",
	},
	{
		browser:  "Safari",
		params:   map[string]bool{"weibo": true},
		uastring: "Mozilla/5.0 (iPhone; CPU iPhone OS 8_4_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Mobile/12H321 Weibo (iPhone7,2__weibo__5.4.0__iphone__os8.4.1)",
	},
	{
		browser:  "QQ Browser",
		params:   map[string]bool{"qq": true, "qq browser": true},
		uastring: "Mozilla/5.0 (Linux; U; Android 4.4.4; zh-cn; MI 4LTE Build/KTU84P) AppleWebKit/533.1 (KHTML, like Gecko)Version/4.0 MQQBrowser/5.4 TBS/025477 Mobile Safari/533.1 V1_AND_SQ_5.9.5_288_YYB_D QQ/5.9.5.2575 NetType/WIFI WebP/0.3.0 Pixel/1080",
	},
	{
		browser:  "Safari",
		params:   map[string]bool{"qq": true},
		uastring: "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Mobile/13B143 QQ/5.9.5.451 Pixel/750 NetType/WIFI Mem/101",
	},
	{
		browser:  "Safari",
		params:   map[string]bool{},
		uastring: "HUAWEI G521-L076_TD/S100 Linux/3.4.39 Android/4.3 Release/08.15.2013 Browser/AppleWebkit534.30 Mobile Safari/534.30",
	},
	{
		browser:  "",
		params:   map[string]bool{},
		uastring: "Xiaomi_2014216_TD-LTE/V1 Linux/3.4.0 Android/4.4.4 Release/20.10.2014 Browser/AppleWebKit537.36 Mobile Safari/537.36 System/Android 4.4.4 XiaoMi/MiuiBrowser/2.0.1",
	},
	{
		params:   map[string]bool{},
		uastring: "Dalvik/1.6.0 (Linux; U; Android 4.4.4; MI 4LTE MIUI/5.11.19)",
	},
	{
		browser:  "Chrome",
		version:  "45",
		params:   map[string]bool{"wechat": true},
		uastring: "Mozilla/5.0 (Linux; Android 6.0; Nexus 6 Build/MRA58N; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/45.0.2454.95 Mobile Safari/537.36 MicroMessenger/6.3.7.51_rbb7fa12.660 NetType/WIFI Language/zh_CN",
	},
	{
		browser:  "Firefox",
		params:   map[string]bool{"firefox": true, "windows cannot catch": true},
		uastring: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:40.0) Gecko/20100101 Firefox/40.0",
	},
	{
		browser:  "QQ Browser",
		params:   map[string]bool{"qq browser": true},
		uastring: "Mozilla/5.0 (iPhone 6p; CPU iPhone OS 8_4 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/6.0 MQQBrowser/6.1.1 Mobile/12H143 Safari/8536.25",
	},
	{
		params: map[string]bool{"qq": true},

		uastring: "Mozilla/5.0 (iPhone; CPU iPhone OS 8_4 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Mobile/12H143 QQ/6.1.0.496 Pixel/1080 NetType/WIFI Mem/49",
	},
	{
		params:   map[string]bool{"wechat": true},
		uastring: "Mozilla/5.0 (iPhone; CPU iPhone OS 8_4 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Mobile/12H143 MicroMessenger/6.3.8 NetType/WIFI Language/zh_CN",
	},

	{
		params:   map[string]bool{"weibo": true},
		uastring: "Mozilla/5.0 (iPhone; CPU iPhone OS 8_4 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Mobile/12H143 Weibo (iPhone7,1__weibo__5.7.0__iphone__os8.4)",
	},
	{
		browser:  "QQ Browser",
		params:   map[string]bool{"qq browser": true},
		uastring: "Mozilla/5.0 (Linux; U; Android 6.0; en-gb; HUAWEI NXT-AL10 Build/HUAWEINXT-AL10) AppleWebKit/537.36 (KHTML, like Gecko)Version/4.0 Chrome/37.0.0.0 MQQBrowser/6.3 Mobile Safari/537.36",
	},
	{
		browser:  "Chrome",
		version:  "45",
		params:   map[string]bool{"qq": true},
		uastring: "Mozilla/5.0 (Linux; Android 6.0; HUAWEI NXT-AL10 Build/HUAWEINXT-AL10; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/45.0.2454.95 Mobile Safari/537.36 V1_AND_SQ_6.1.0_312_YYB_D PA QQ/6.1.0.2635 NetType/WIFI WebP/0.4.1 Pixel/1080",
	},
	{
		browser:  "QQ Browser",
		params:   map[string]bool{"qq browser": true, "wechat": true},
		uastring: "Mozilla/5.0 (Linux; U; Android 4.4.4; zh-cn; MI NOTE LTE Build/KTU84P) AppleWebKit/533.1 (KHTML, like Gecko)Version/4.0 MQQBrowser/5.4 TBS/025483 Mobile Safari/533.1 MicroMessenger/6.3.8.50_r251a77a.680 NetType/WIFI Language/zh_CN",
	},
	{
		browser:  "Chrome",
		version:  "45",
		params:   map[string]bool{"weibo": true},
		uastring: "Mozilla/5.0 (Linux; Android 6.0; HUAWEI NXT-AL10 Build/HUAWEINXT-AL10; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/45.0.2454.95 Mobile Safari/537.36 Weibo (HUAWEI-HUAWEI NXT-AL10__weibo__4.3.0__android__android6.0)",
	},
	{
		params:   map[string]bool{"sogou search": true},
		uastring: "Mozilla/5.0 (Linux; Android 6.0; HUAWEI NXT-AL10 Build/HUAWEINXT-AL10; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/45.0.2454.95 Mobile Safari/537.36 SogouSearch Android1.0 version3.0 ",
	},
	{
		params:   map[string]bool{"uc browser": true},
		uastring: "Mozilla/5.0 (Linux; U; Android 6.0; zh-CN; HUAWEI NXT-AL10 Build/HUAWEINXT-AL10) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 UCBrowser/10.9.0.703 U3/0.8.0 Mobile Safari/534.30",
	},
	{
		params:   map[string]bool{"cannot go market": true},
		uastring: "Mozilla/5.0 (Linux; Android 6.0; HUAWEI NXT-AL10 Build/HUAWEINXT-AL10; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/45.0.2454.95 Mobile Safari/537.36 MxBrowser/4.5.6.1000 ",
	},
	{
		params:   map[string]bool{"windows cannot catch": true},
		uastring: "Mozilla/5.0 (Linux; Android 6.0; HUAWEI NXT-AL10 Build/HUAWEINXT-AL10; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/45.0.2454.95 Mobile Safari/537.36 LieBaoFast/3.25.7 ",
	},
	{
		params:   map[string]bool{"windows cannot catch": true},
		uastring: "Mozilla/5.0 (Linux; U; Android 6.0; zh-CN; HUAWEI NXT-AL10 Build/HUAWEINXT-AL10) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 OPR/11.1.1.99254 Mobile Safari/537.36 ",
	},
	{
		params:   map[string]bool{"windows cannot catch": true},
		uastring: "Mozilla/5.0 (Linux; Android 6.0; HUAWEI NXT-AL10 Build/HUAWEINXT-AL10; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/45.0.2454.95 Mobile Safari/537.36; 360 Aphone Browser (6.9.9.39) ",
	},
	{
		params:   map[string]bool{"firefox": true, "windows cannot catch": true},
		uastring: "Mozilla/5.0 (Android 6.0; Mobile; rv:43.0) Gecko/43.0 Firefox/43.0 ",
	},
	{
		params:   map[string]bool{"cannotdeeplink": true},
		uastring: "Mozilla/6.0(Linux; Android 6.0; HUAWEI NXT-AL10 Build/HUAWEINXT-AL10) AppleWebKit/537.36(KHTML,like Gecko) Version/6.0 Mobile Safari/537.36",
	},
	{
		params:   map[string]bool{"cannotdeeplink": true, "windows cannot catch": true},
		uastring: "Mozilla/5.0 (Linux; Android 6.0; HUAWEI NXT-AL10 Build/HUAWEINXT-AL10) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/35.0.1916.138 Mobile Safari/537.36 baidubrowser/6.3.20.0 (Baidu; P1 6.0)",
	},

	{
		params:   map[string]bool{"windows cannot catch": true},
		uastring: "Mozilla/5.0 (Linux; Android 6.0; HUAWEI NXT-AL10 Build/HUAWEINXT-AL10; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/45.0.2454.95 Mobile Safari/537.36 baidubrowser/5.4.5.0 (Baidu; P1 6.0)",
	},
	{
		params:   map[string]bool{},
		uastring: "Mozilla/5.0 (Linux; Android 6.0; HUAWEI NXT-AL10 Build/HUAWEINXT-AL10; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/45.0.2454.95 Mobile Safari/537.36 baidubrowser/4.0.18.7 (Baidu; P1 6.0)",
	},
}

func TestAllPatch(t *testing.T) {
	curdir, err := path.Getcurdir()
	if err != nil {
		t.Fatal(err)
	}
	regexFile := curdir + "/regexes.yaml"
	parser, err := uaparser.New(regexFile)
	if err != nil {
		t.Fatal(err)
	}

	for i, expected := range uas {
		client := parser.Parse(expected.uastring)
		browser := client.UserAgent.Family
		version := client.UserAgent.Major
		params := client.Channel.Params
		if !strings.Contains(browser, expected.browser) {
			t.Errorf("#%d: parse browser failed, expect: %v, got: %v for ua: %v\n", i, expected.browser, browser, expected.uastring)
		}
		if expected.browser == "Chrome" && version != expected.version {
			t.Errorf("#%d: parse version failed, expect: %v, got: %v for ua: %v\n", i, expected.version, version, expected.uastring)
		}
		if !reflect.DeepEqual(params, expected.params) {
			t.Errorf("#%d: parse channel failed, expect: %v, got: %v for ua: %v\n", i, expected.params, params, expected.uastring)
		}
	}
}
