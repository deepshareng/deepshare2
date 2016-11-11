package match

import "github.com/MISingularity/deepshare2/api"

type MatchPayload struct {
	SenderInfo api.SenderInfo `json:"sender_info"`
	InappData  string         `json:"inapp_data"`
	ClientIP   string         `json:"client_ip"`
	ClientUA   string         `json:"client_ua"`
}

type SenderInfoObj struct {
	SenderID string   `json:"sender_id"`
	Channels []string `json:"channels"`
	SDKInfo  string   `json:"sdk_info"`
}

type MatchRequestBody struct {
	InAppData  string        `json:"inapp_data"`
	SenderInfo SenderInfoObj `json:"sender_info"`
	ClientIP   string        `json:"client_ip"`
	ClientUA   string        `json:"client_ua"`
}
