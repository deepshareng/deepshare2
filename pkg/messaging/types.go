package messaging

import (
	"encoding/json"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
)

var (
	CounterTopic           = []byte("counter_raw")
	CounterTopicAttributed = []byte("counter")
	MatchTopic             = []byte("match")
	DSActionTopic          = []byte("dsaction")
	ShareLinkTopic         = []byte("sharelink")
	GenUrlTopic            = []byte("genurl")
	InAppDataTopic         = []byte("inappdata")
)

type Event struct {
	AppID        string                 `json:"appid"`
	EventType    string                 `json:"event"`
	Channels     []string               `json:"channels,omitempty"`
	SenderID     string                 `json:"sender_id,omitempty"`
	CookieID     string                 `json:"cookie_id,omitempty"`
	UniqueID     string                 `json:"unique_id,omitempty"`
	Count        int                    `json:"count"`
	UAInfo       uainfo.UAInfo          `json:"ua_info"`
	ReceiverInfo api.MatchReceiverInfo  `json:"receiver_info"`
	KVs          map[string]interface{} `json:"kvs"`
	TimeStamp    int64                  `json:"timestamp"`
}

func (e *Event) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Event) Unmarshal(b []byte) error {
	return json.Unmarshal(b, e)
}
