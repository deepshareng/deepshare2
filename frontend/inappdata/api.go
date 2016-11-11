package inappdata

const (
	RequestReceiverInfoKey      = "receiver_info"
	RequestTrackingKey          = "tracking"
	MatchRequestClientIPKey     = "client_ip"
	MatchRequestClientUAKey     = "client_ua"
	RequestTrackingValueInstall = "install"
	RequestTrackingValueOpen    = "open"
	CookieRequestDeviceIDKey    = "deviceID"
	noValueFromSDK              = "dl_no_value"
)

type InAppDataPostBody struct {
	IsNewUser bool `json:"is_newuser"`
	//Below fields only possible to be sent when it is a old user
	ClickID    string `json:"click_id"`
	WCookie    string `json:"wcookie"`
	DeeplinkID string `json:"deeplink_id"`
	ShortSeg   string `json:"short_seg"`

	//Below fields should be keep identical as match.MatchRequestReceiverInfo
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
	//Below fields only need to be sent when it is a new user
	HardwareID string `json:"hardware_id"`

	//TODO we need to discuss if below fields should be sent
	UriScheme string `json:"uri_scheme"`
}

type InAppDataResponseBody struct {
	InAppData string   `json:"inapp_data"`
	Channels  []string `json:"channels"`
}
