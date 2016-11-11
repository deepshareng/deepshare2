package appcookiedevice

const (
	PairCookiesPath   = "paircookies"
	RefreshCookiePath = "refreshcookie"
)

type DeviceInfo struct {
	UniqueId string `json:"unique_id"`
}

type PostPairCookieBody struct {
	Cookie1 string `json:"cookie1"`
	Cookie2 string `json:"cookie2"`
}

type PostRefreshCookieBody struct {
	Cookie    string `json:"cookie"`
	NewCookie string `json:"new_cookie"`
}
