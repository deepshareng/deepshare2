package appinfo

const (
	keyAppIDPrefix    = "app:"
	keyShortIDToAppID = "shortid_to_appid"
	keyAppIDToShortID = "appid_to_shortid"
)

type AppInfo struct {
	AppID     string
	ShortID   string
	AppName   string
	Android   AppAndroidInfo
	Ios       AppIosInfo
	YYBUrl    string
	YYBEnable bool
	IconUrl   string
	Theme     string
	UserConf  UserConfig
}

type AppAndroidInfo struct {
	Scheme             string
	Host               string
	Pkg                string
	DownloadUrl        string
	IsDownloadDirectly bool
	YYBEnable          bool
}

type AppIosInfo struct {
	Scheme              string
	BundleID            string
	TeamID              string
	DownloadUrl         string
	UniversalLinkEnable bool
	YYBEnableBelow9     bool
	YYBEnableAbove9     bool
	ForceDownload       bool
}

type UserConfig struct {
	BgWeChatAndroidTipUrl string
	BgWeChatIosTipUrl     string
}
