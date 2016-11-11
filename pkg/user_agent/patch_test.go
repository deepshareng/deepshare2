package user_agent

import (
	"testing"
)

var uas = []struct {
	os                 string
	osversion          string
	brand              string
	iosMajorVersion    int
	isWechat           bool
	isWeibo            bool
	chromeMajorVersion int
	isFireFox          bool
	isQQ               bool
	uastring           string
}{
	{
		os:                 "android",
		osversion:          "4.4.4",
		brand:              "-",
		iosMajorVersion:    0,
		isWechat:           false,
		isWeibo:            false,
		chromeMajorVersion: 33,
		isFireFox:          false,
		uastring:           "Mozilla/5.0 (Linux; Android 4.4.4; Android SDK built for x86 Build/KK) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/33.0.0.0 Mobile Safari/537.36",
	},
	{
		os:                 "ios",
		osversion:          "8.1.3",
		iosMajorVersion:    8,
		isWechat:           false,
		isWeibo:            false,
		chromeMajorVersion: 40,
		isFireFox:          false,
		uastring:           "Mozilla/5.0 (iPhone; CPU iPhone OS 8_1_3 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) CriOS/40.0.2214.69 Mobile/12B466 Safari/600.1.4",
	},
	{
		os:                 "android",
		osversion:          "4.1.1",
		brand:              "mi",
		iosMajorVersion:    0,
		isWechat:           true,
		isWeibo:            false,
		chromeMajorVersion: 0,
		isFireFox:          false,
		uastring:           "Mozilla/5.0 (Linux; U; Android 4.1.1; zh-cn; MI 2A Build/JRO03L) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30 MicroMessenger/5.4.0.51_r798589.480 NetType/WIFI",
	},
	{
		os:                 "ios",
		osversion:          "8.4.1",
		iosMajorVersion:    8,
		isWechat:           false,
		isWeibo:            false,
		chromeMajorVersion: 45,
		isFireFox:          false,
		uastring:           "Mozilla/5.0 (iPhone; CPU iPhone OS 8_4_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) CriOS/45.0.2454.89 Mobile/12H321 Safari/600.1.4",
	},
	{
		os:                 "ios",
		osversion:          "8.4.1",
		iosMajorVersion:    8,
		isWechat:           true,
		isWeibo:            false,
		chromeMajorVersion: 0,
		isFireFox:          false,
		uastring:           "Mozilla/5.0 (iPhone; CPU iPhone OS 8_4_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Mobile/12H321 MicroMessenger/6.2.6 NetType/WIFI Language/zh_CN",
	},
	{
		os:                 "ios",
		osversion:          "8.4.1",
		iosMajorVersion:    8,
		isWechat:           false,
		isWeibo:            true,
		chromeMajorVersion: 0,
		isFireFox:          false,
		uastring:           "Mozilla/5.0 (iPhone; CPU iPhone OS 8_4_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Mobile/12H321 Weibo (iPhone7,2__weibo__5.4.0__iphone__os8.4.1)",
	},
	{
		os:                 "android",
		osversion:          "4.4.4",
		iosMajorVersion:    0,
		isWechat:           false,
		isWeibo:            false,
		isQQ:               true,
		chromeMajorVersion: 0,
		isFireFox:          false,
		uastring:           "Mozilla/5.0 (Linux; U; Android 4.4.4; zh-cn; MI 4LTE Build/KTU84P) AppleWebKit/533.1 (KHTML, like Gecko)Version/4.0 MQQBrowser/5.4 TBS/025477 Mobile Safari/533.1 V1_AND_SQ_5.9.5_288_YYB_D QQ/5.9.5.2575 NetType/WIFI WebP/0.3.0 Pixel/1080",
	},
	{
		os:                 "ios",
		osversion:          "9.1",
		iosMajorVersion:    9,
		isWechat:           false,
		isWeibo:            false,
		isQQ:               true,
		chromeMajorVersion: 0,
		isFireFox:          false,
		uastring:           "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Mobile/13B143 QQ/5.9.5.451 Pixel/750 NetType/WIFI Mem/101",
	},
	{
		os:                 "android",
		osversion:          "4.3",
		iosMajorVersion:    0,
		isWechat:           false,
		isWeibo:            false,
		isQQ:               false,
		chromeMajorVersion: 0,
		isFireFox:          false,
		uastring:           "HUAWEI G521-L076_TD/S100 Linux/3.4.39 Android/4.3 Release/08.15.2013 Browser/AppleWebkit534.30 Mobile Safari/534.30",
	},
	{
		os:                 "android",
		osversion:          "4.4.4",
		iosMajorVersion:    0,
		isWechat:           false,
		isWeibo:            false,
		isQQ:               false,
		chromeMajorVersion: 0,
		isFireFox:          false,
		uastring:           "Xiaomi_2014216_TD-LTE/V1 Linux/3.4.0 Android/4.4.4 Release/20.10.2014 Browser/AppleWebKit537.36 Mobile Safari/537.36 System/Android 4.4.4 XiaoMi/MiuiBrowser/2.0.1",
	},
	{
		os:                 "android",
		osversion:          "4.4.4",
		iosMajorVersion:    0,
		isWechat:           false,
		isWeibo:            false,
		isQQ:               false,
		chromeMajorVersion: 0,
		isFireFox:          false,
		uastring:           "Dalvik/1.6.0 (Linux; U; Android 4.4.4; MI 4LTE MIUI/5.11.19)",
	},
	{
		os:                 "android",
		osversion:          "6.0",
		iosMajorVersion:    0,
		isWechat:           true,
		isWeibo:            false,
		isQQ:               false,
		chromeMajorVersion: 45,
		isFireFox:          false,
		uastring:           "Mozilla/5.0 (Linux; Android 6.0; Nexus 6 Build/MRA58N; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/45.0.2454.95 Mobile Safari/537.36 MicroMessenger/6.3.7.51_rbb7fa12.660 NetType/WIFI Language/zh_CN",
	},
}

func TestAllPatch(t *testing.T) {
	for _, expected := range uas {
		uao := New(expected.uastring)
		os, osv := uao.OsNameVersion()
		isWechat := uao.IsWechat()
		isWeibo := uao.IsWeibo()
		isQQ := uao.IsQQ()
		chromeV := uao.ChromeMajorVersion()
		isFirefox := uao.IsFirefox()
		brand := uao.Brand()
		if os != expected.os {
			t.Errorf("parse os failed, expect: %v, got: %v for ua: %v\n", expected.os, os, expected.uastring)
		}
		if osv != expected.osversion {
			t.Errorf("parse os version failed, expect: %v, got: %v for ua: %v\n", expected.osversion, osv, expected.uastring)
		}
		if isWechat != expected.isWechat {
			t.Errorf("parse isWechat failed, expect: %v, got: %v for ua: %v\n", expected.isWechat, isWechat, expected.uastring)
		}
		if isWeibo != expected.isWeibo {
			t.Errorf("parse isWeibo failed, expect: %v, got: %v for ua: %v\n", expected.isWeibo, isWeibo, expected.uastring)
		}
		if isQQ != expected.isQQ {
			t.Errorf("parse isQQ failed, expect: %v, got: %v for ua: %v\n", expected.isQQ, isQQ, expected.uastring)
		}
		if chromeV != expected.chromeMajorVersion {
			t.Errorf("parse chromeMajorVersion failed, expect: %v, got: %v for ua: %v\n", expected.chromeMajorVersion, chromeV, expected.uastring)
		}
		if isFirefox != expected.isFireFox {
			t.Errorf("parse isFireFox failed, expect: %v, got: %v for ua: %v\n", expected.isFireFox, isFirefox, expected.uastring)
		}
		if expected.brand != "" && brand != expected.brand {
			t.Errorf("parse brand failed, expect: %v, got: %v for ua: %v\n", expected.brand, brand, expected.uastring)
		}
	}
}
