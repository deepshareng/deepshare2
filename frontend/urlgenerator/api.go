package urlgenerator

import (
	"encoding/json"
	"errors"
)

type ShareLinkInfoObj struct {
}

type GenURLPostBody struct {
	InAppDataReq       interface{} `json:"inapp_data"`
	InAppData          string      `json:"-"`
	DownloadTitle      string      `json:"download_title"`
	DownloadBtnText    string      `json:"download_btn_text"`
	DownloadMsg        string      `json:"download_msg"`
	DownloadUrlIos     string      `json:"download_url_ios"`
	DownloadUrlAndroid string      `json:"download_url_android"`
	//uninstall url is typically the wap page url, useful for some customer who what user visit there wap site other than installing the app
	UninstallUrl      string   `json:"uninstall_url"`
	RedirectUrl       string   `json:"redirect_url"`
	IsShort           bool     `json:"is_short"`
	IsPermanent       bool     `json:"is_permanent"`
	UseShortID        bool     `json:"use_shortid"`
	SDKInfo           string   `json:"sdk_info"`            //android1.3.1
	SenderID          string   `json:"sender_id"`           //unique_id for sdk, indicate the url is generated from sender
	ForwardedSenderID string   `json:"forwarded_sender_id"` //unique_id for sdk, indicate the url is generated from receiver
	Channels          []string `json:"channels"`
	ForceDownload     string   `json:"force_download"`
}

func StringfyInappData(gupb *GenURLPostBody) (*GenURLPostBody, error) {
	if gupb.InAppDataReq == nil {
		gupb.InAppData = ""
		return gupb, nil
	}
	switch gupb.InAppDataReq.(type) {
	case string:
		gupb.InAppData = gupb.InAppDataReq.(string)
		return gupb, nil
	case map[string]interface{}:
		if b, err := json.Marshal(gupb.InAppDataReq); err != nil {
			return gupb, err
		} else {
			gupb.InAppData = string(b)
			return gupb, err
		}
	}
	return gupb, errors.New("InAppData type error")
}

type GenURLResponseBody struct {
	Url  string `json:"url"`
	Path string `json:"path"`
}
