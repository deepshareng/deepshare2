package api

// Producer Definition
const (
	versionPrefix            = "/v2"
	MatchPrefix              = versionPrefix + "/matches/"
	CounterPrefix            = versionPrefix + "/counters/"
	DSActionsPrefix          = versionPrefix + "/dsactions/"
	AppCookieDevicePrefix    = versionPrefix + "/appcookiedevice/"
	DeviceCookiePrefix       = versionPrefix + "/devicecookie/"
	AppInfoPrefix            = versionPrefix + "/appinfo/"
	GenerateUrlPrefix        = versionPrefix + "/url/"
	GetInAppDataPrefix       = versionPrefix + "/inappdata/"
	BindDeviceToCookiePrefix = versionPrefix + "/binddevicetocookie/"
	ShareLinkPrefix          = "/d/"
	JSServerPrefix           = "/jsserver/"
	ChannelPrefix            = versionPrefix + "/channels/"
	AppChannelPrefix         = versionPrefix + "/appchannels/"
	AppEventPrefix           = versionPrefix + "/appevents/"
	CookieName               = "dscookie"
	WechatCookieQueryKey     = "wcookie"
	CookiePlaceHolder        = "cookie_placeholder"
	DSUsagesPrefix           = versionPrefix + "/dsusages/"
	RetentionPrefix          = versionPrefix + "/retention/"
	RetentionAmountPrefix    = versionPrefix + "/retention-amount/"
	JSApiPrefix              = versionPrefix + "/jsapi/"
	TokenPrefix              = versionPrefix + "/tokens/"
	DeviceStatPrefix         = versionPrefix + "/device-stat/"
)

type MatchResponse struct {
	InappData string   `json:"inapp_data"`
	SenderID  string   `json:"sender_id"`
	Channels  []string `json:"channels"`
}

type SenderInfo struct {
	SenderID string   `json:"sender_id"`
	Channels []string `json:"channels"`
}

type ReceiverInfo struct {
	UniqueID string `json:"unique_id"`
}

type MatchReceiverInfo struct {
	//Below fields should be keep identical as inappdata.InAppDataPostBody
	UniqueID       string `json:"unique_id"`
	AppVersionName string `json:"app_version_name"`
	//For android version code
	AppVersionCode int `json:"app_version_code"`
	//For IOS build
	AppVersionBuild  string `json:"app_version_build"`
	SDKInfo          string `json:"sdk_info"`
	ISWifiConnected  bool   `json:"is_wifi_connected"`
	CarrierName      string `json:"carrier_name"`
	BlueToothEnable  bool   `json:"blueTooth_enable"`
	Model            string `json:"model"`
	Brand            string `json:"brand"`
	IsEmulator       bool   `json:"is_emulator"`
	OS               string `json:"os"`
	OSVersion        string `json:"os_version"`
	HasNfc           bool   `json:"has_nfc"`
	HasTelephony     bool   `json:"has_telephone"`
	BlueToothVersion string `json:"bluetooth_version"`
	ScreenDPI        int    `json:"screen_dpi"`
	ScreenWidth      int    `json:"screen_width"`
	ScreenHeight     int    `json:"screen_height"`
	//TODO we need to discuss if below fields should be sent
	UriScheme string `json:"uri_scheme"`
	//Below fields only need to be sent when it is a new user
	HardwareID string `json:"hardware_id"`
}
