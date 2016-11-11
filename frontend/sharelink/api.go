package sharelink

const (
	jsTemplateDir      = "/js_template/"
	jsFileServerDir    = jsTemplateDir + "/fileserver/"
	shortSegForInvalid = "0"
	dsTAG              = "ds_tag"
)

type appInsStatus int

const (
	AppNotInstalled   appInsStatus = iota
	AppInstalled      appInsStatus = iota
	AppInstallUnClear appInsStatus = iota
)

//Raw URL fields
const (
	DownloadTitle      = "download_title"
	DownloadBtnText    = "download_btn_text"
	DownloadMsg        = "download_msg"
	DownloadUrlIos     = "download_url_ios"
	DownloadUrlAndroid = "download_url_android"
	UninstallUrl       = "uninstall_url"
	RedirectUrl        = "redirect_url"
	InAppData          = "inapp_data"
	SenderID           = "sender_id"
	Channels           = "channels"
	SDKInfo            = "sdk_info"
	ForceDownload      = "force_download"
)

type AppCookieDeviceResponseBody struct {
	UniqueID string `json:"unique_id"`
}

type DownloadInfo struct {
	DownloadTitle      string
	DownloadBtnText    string
	DownloadMsg        string
	DownloadUrlIos     string
	DownloadUrlAndroid string
	UninstallUrl       string
	ForceDownload      string
}
