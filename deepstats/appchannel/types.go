package appchannel

type AppChannel struct {
	AppID   string `bson:"app_id,omitempty"`
	Channel string `bson:"channel,omitempty"`
}

type AppChannels struct {
	AppID    string
	Channels []string
}
