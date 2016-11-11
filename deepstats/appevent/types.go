package appevent

type AppEvent struct {
	AppID string `bson:"app_id,omitempty"`
	Event string `bson:"event,omitempty"`
}

type AppEvents struct {
	AppID  string
	Events []string
}
