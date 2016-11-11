package job

import "time"

type LastTimeServiceTime struct {
	ServiceName string    `bson:"service-name"`
	LastTime    time.Time `bson:"lasttime"`
}

type LastTimeChannelTime struct {
	ServiceName string    `bson:"service-name"`
	Channel     string    `bson:"channel"`
	LastTime    time.Time `bson:"lasttime"`
}
