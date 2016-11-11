##genurl

```
curl -v -X POST -H 'content-type: application/json' --data '{}' https://fds.so/v2/url/54378fd28a6c81e3
```

检查点：

	1) http code 
		200
	2) response body 格式为
		{"url":"https://fds.so/d/54378fd28a6c81e3/5KSBpQmsg0"}

##sharelink

将 "https://fds.so/d/54378fd28a6c81e3/5KSBpQmsg0" 替换为上一步返回的deepshare链接。

```
curl -v -b 'dscookie=CCCCC' -H 'User-Agent:Mozilla/5.0 (iPhone; CPU iPhone OS 9_0 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13A344 Safari/601.1' https://fds.so/d/54378fd28a6c81e3/5KSBpQmsg0
```

检查点：

	1) http code 
		200
	2) Set-Cookie: dscookie=CCCCC
	3) response body 
		<html>开头
		在其中能搜到 "deepshare-redirect.min.js"

##binddevicetocookie

```
curl -v -b 'dscookie=CCCCC' https://fds.so/v2/binddevicetocookie/uuuu
```

检查点：

	1) http code
		200


##inappdata

```
curl -v -X POST -H 'content-type: application/json' --data '{"is_newuser":false,"click_id":"","os":"iOS","os_version":"9.1"}' https://fds.so/v2/inappdata/54378fd28a6c81e3
```

检查点：

	1) http code
		200
	2) response body 格式为
		{"inapp_data":"","channels":[]}

##dsusage
```
curl -v  https://fds.so/v2/dsusages/54378fd28a6c81e3/test_sender
```

检查点：

	1) http code
		200
	2) response body 格式为
		{"new_install":0,"new_open":1}
		
##jsapi
```
curl -v -X POST -H 'content-type: application/json' -H 'User-Agent:Mozilla/5.0 (iPhone; CPU iPhone OS 9_0 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13A344 Safari/601.1' --data '[{"deeplink_id":"1","inapp_data":{"name":"n1"}}]' https://fds.so/v2/jsapi/54378fd28a6c81e3
```

检查点:

	1) http code
		200
	2) response body 格式为
	```json
		{"app_id":"54378fd28a6c81e3","ds_urls":{"1":"https://fds.so/d/54378fd28a6c81e3/6Ahy4UIQmY"},"chrome_major":0,"is_android":false,"is_ios":true,"ios_major":9,"is_wechat":false,"is_weibo":false,"is_qq":false,"is_facebook":false,"is_twitter":false,"is_firefox":false,"is_qq_browser":false,"is_uc":false,"cannot_deeplink":false,"cannot_get_win_event":false,"cannot_go_market":false,"force_use_scheme":false,"app_name":"testapi","icon_url":"https://nzjddxpun.qnssl.com/54378fd28a6c81e3/appicon/20151225051527","scheme":"ds54378fd28a6c81e3","host":"","bundle_id":"com.singulariti.testapi","pkg":"","url":"http://baidu.com","is_download_directly":false,"is_universal_link":true,"is_yyb_enable_ios_below_9":false,"is_yyb_enable_ios_above_9":false,"is_yyb_enable_android":false,"yyb_url":"","match_id":"6Ahy4QKsy4","timestamp":1466754103,"ds_tag":"","app_ins_status":0}
	```

##counter
```
curl -v -X POST -H 'content-type: application/json' --data '{"receiver_info":{"unique_id":"testdevice_curl"},"counters":[{"event":"curl_test","count":5}]}' https://fds.so/v2/counters/54378fd28a6c81e3
```

检查点:

	1) http code
		200

##dsaction
```
curl -v -X POST -H 'content-type: application/json' --data '{"action":"js/curl"}' https://fds.so/v2/dsactions/54378fd28a6c81e3
```

检查点:

	1) http code
		200