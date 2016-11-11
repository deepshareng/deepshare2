package attribpush

const AttributionTopic = "attribution"

type AttributionPushInfo struct {
	SenderID  string `json:"sender_id"`
	Tag       string `json:"tag"`
	Value     int    `json:"value"`
	Timestamp int64  `json:"timestamp"`
}
