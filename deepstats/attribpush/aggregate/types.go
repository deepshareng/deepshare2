package aggregate

type ValueTimePair struct {
	Value     int   `json:"value" bson:"v"`
	Timestamp int64 `json:"timestamp" bson:"t"`
}
type AttributionInfo struct {
	AppID        string                     `json:"app_id" bson:"app_id"`
	SenderID     string                     `json:"sender_id" bson:"sender_id"`
	TaggedValues map[string][]ValueTimePair `json:"tagged_values" bson:"tagged_values"`
}
