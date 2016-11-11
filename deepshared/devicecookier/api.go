package devicecookier

const (
	UniqueIDPrefix        = "uid_"
	UniqueIDPrefixWCookie = "uid:w_"
	HardwareIDPrefix      = "hid_"
	CookieIDPrefix        = "cid:uid_"
	WCookieIDPrefix       = "wcid:uid_"
)

const (
	GetDeviceByCookieQuery = "getdevicebycookie"
)

const (
	ApiCookiesPrefix = "cookies"
	ApiDevicesPrefix = "devices"
)

type CookieInfo struct {
	CookieID string `json:"cookie_id"`
}

type DeviceInfo struct {
	DeviceID string `json:"device_id"`
}
