package dsaction

import "github.com/MISingularity/deepshare2/api"

type DSActionRequest struct {
	Action       string                 `json:"action"`
	Channels     []string               `json:"channels"`
	SenderID     string                 `json:"sender_id"`
	ReceiverInfo api.ReceiverInfo       `json:"receiver_info"`
	KVs          map[string]interface{} `json:"kvs"`
}
