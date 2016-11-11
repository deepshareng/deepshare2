package uainfo

import (
	"strconv"
	"strings"
)

// UAInfo holds user information parsed from http request.
type UAInfo struct {
	Ua                    string
	Ip                    string
	Os                    string
	OsVersion             string
	Brand                 string
	Browser               string
	IsWechat              bool
	IsWeibo               bool
	IsQQ                  bool
	IsTwitter             bool
	IsFacebook            bool
	IsFirefox             bool
	IsQQBrowser           bool
	ChromeMajor           int
	CannotDeeplink        bool
	CannotGoMarket        bool
	CannotGetWindowsEvent bool
	ForceUseScheme        bool
}

func (ui *UAInfo) IsAndroid() bool {
	return ui.Os == "android"
}

func (ui *UAInfo) IsIos() bool {
	return ui.Os == "ios"
}

func (ui *UAInfo) IosMajorVersion() int {
	if !ui.IsIos() {
		return 0
	}
	return ui.majorVersion()
}

func (ui *UAInfo) majorVersion() int {
	osv := ui.OsVersion
	if osv == "" {
		return 0
	}

	idx := strings.Index(osv, ".")
	majorStr := osv
	if idx != -1 {
		majorStr = osv[:idx]
	}
	n, err := strconv.Atoi(majorStr)
	if err == nil && n > 0 {
		return n
	}
	return 0
}

type UATransformer interface {
	Transform() string
	Os() string
}

func NewUAFingerPrinter(ua *UAInfo) *UAFingerPrinter {
	return &UAFingerPrinter{ua}
}

type UAFingerPrinter struct {
	ua *UAInfo
}

func (u *UAFingerPrinter) Transform() string {
	if u.ua.Ua == "" {
		return ""
	}
	return ("UA:" + u.ua.Ip + "_" + u.ua.Os + "_" + u.ua.OsVersion + "_" + u.ua.Brand)
}

func (u *UAFingerPrinter) Os() string {
	return u.ua.Os
}
