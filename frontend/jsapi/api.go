package jsapi

const (
	clickedTag        = "clicked"
	eventSuffixClickd = ":clicked"
)

type JsApiPost []struct {
	DeepLinkID         string                 `json:"deeplink_id"`
	InAppDataReq       map[string]interface{} `json:"inapp_data"`
	InAppData          string                 `json:"-"`
	SenderID           string                 `json:"sender_id"`
	Channels           []string               `json:"channels"`
	DownloadTitle      string                 `json:"download_title"`
	DownloadBtnText    string                 `json:"download_btn_text"`
	DownloadMsg        string                 `json:"download_msg"`
	DownloadUrlIos     string                 `json:"download_url_ios"`
	DownloadUrlAndroid string                 `json:"download_url_android"`
}

type JsApiResp struct {
	AppID string `json:"app_id"`

	DSUrls map[string]string `json:"ds_urls"`

	//From UA info
	ChromeMajor       int  `json:"chrome_major"`
	IsAndroid         bool `json:"is_android"`
	IsIos             bool `json:"is_ios"`
	IosMajor          int  `json:"ios_major"`
	IsWechat          bool `json:"is_wechat"`
	IsWeibo           bool `json:"is_weibo"`
	IsQQ              bool `json:"is_qq"`
	IsFacebook        bool `json:"is_facebook"`
	IsTwitter         bool `json:"is_twitter"`
	IsFirefox         bool `json:"is_firefox"`
	IsQQBrowser       bool `json:"is_qq_browser"`
	IsUC              bool `json:"is_uc"`
	CannotDeeplink    bool `json:"cannot_deeplink"`
	CannotGetWinEvent bool `json:"cannot_get_win_event"`
	CannotGoMarket    bool `json:"cannot_go_market"`
	ForceUseScheme    bool `json:"force_use_scheme"`

	//From App info
	AppName              string `json:"app_name"`
	IconUrl              string `json:"icon_url"`
	Scheme               string `json:"scheme"`
	Host                 string `json:"host"`
	BundleID             string `json:"bundle_id"`
	Pkg                  string `json:"pkg"`
	Url                  string `json:"url"`
	IsDownloadDirectly   bool   `json:"is_download_directly"`
	IsUniversallink      bool   `json:"is_universal_link"`
	IsYYBEnableIosBelow9 bool   `json:"is_yyb_enable_ios_below_9"`
	IsYYBEnableIosAbove9 bool   `json:"is_yyb_enable_ios_above_9"`
	IsYYBEnableAndroid   bool   `json:"is_yyb_enable_android"`
	YYBUrl               string `json:"yyb_url"`

	//From sender info
	MatchId      string `json:"match_id"`
	Timestamp    int64  `json:"timestamp"`
	DsTag        string `json:"ds_tag"`
	AppInsStatus int    `json:"app_ins_status"`
}

type JsApiPostClicked struct {
	InAppDataReq map[string]interface{} `json:"inapp_data"`
	InAppData    string                 `json:"-"`
	SenderID     string                 `json:"sender_id"`
	Channels     []string               `json:"channels"`
}

type JsApiRespClicked struct {
	OK bool `json:"ok"`
}
