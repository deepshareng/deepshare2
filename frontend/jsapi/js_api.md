JS Api is the entry for JS, accessed by JS SDK. 

#### POST /v2/jsapi/:appID

Request:

```json
[
  {
    "deeplink_id": "1",
    "inapp_data": {
      "name": "n1"
    },
    "sender_id": "s1",
    "channels": [
      "ch1_x",
      "ch1_y"
    ],
    "download_title": "download_title1",
    "download_btn_text": "download_btn_text1",
    "download_msg": "download_msg1",
    "download_url_ios": "",
    "download_url_android": ""
  },
  {
    "deeplink_id": "2",
    "inapp_data": {
      "name": "n2"
    },
    "sender_id": "s2",
    "channels": [
      "ch2_x",
      "ch2_y"
    ],
    "download_title": "download_title2",
    "download_btn_text": "download_btn_text2",
    "download_msg": "download_msg2",
    "download_url_ios": "",
    "download_url_android": ""
  }
]
```

Response:
```json
{
  "app_id": "38CCA4C77072DDC9",
  "ds_urls": {
    "1": "https://fds.so/d/38CCA4C77072DDC9/4ToyLAHuFi",
    "2": "https://fds.so/d/38CCA4C77072DDC9/4ToyLFOh0s"
  },
  "chrome_major": 44,
  "is_android": true,
  "is_ios": false,
  "ios_major": 0,
  "is_wechat": false,
  "is_weibo": false,
  "is_qq": false,
  "is_firefox": false,
  "is_qq_browser": false,
  "is_uc": false,
  "cannot_deeplink": false,
  "cannot_get_win_event": false,
  "cannot_go_market": false,
  "force_use_scheme": false,
  "app_name": "linux-command",
  "icon_url": "https://nzjddxpun.qnssl.com/38CCA4C77072DDC9/appicon/20160408063820",
  "scheme": "ds38CCA4C77072DDC9",
  "host": "com.misingularity.linuxcommand",
  "bundle_id": "",
  "pkg": "com.misingularity.linuxcommand",
  "url": "http://a.app.qq.com/o/simple.jsp?pkgname=com.misingularity.linuxcommand",
  "is_download_directly": false,
  "is_universal_link": false,
  "is_yyb_enable_ios_below_9": false,
  "is_yyb_enable_ios_above_9": false,
  "is_yyb_enable_android": true,
  "yyb_url": "http://a.app.qq.com/o/simple.jsp?pkgname=com.misingularity.linuxcommand",
  "match_id": "460AKWEXqE",
  "timestamp": 1461294074,
  "ds_tag": "",
  "app_ins_status": 0
}
```

#### POST /v2/jsapi/:appID?clicked=true

Request:

```json
{
  "inapp_data": {
    "name": "n1"
  },
  "sender_id": "s1",
  "channels": [
    "ch1_x",
    "ch1_y"
  ]
}
```

Response:

```json
{
  "ok": true
}
```