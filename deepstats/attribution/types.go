package attribution

import (
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"gopkg.in/mgo.v2/bson"
)

var (
	nsqProduceTopic = messaging.CounterTopicAttributed
)

type Attribution struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"-"`
	AppID    string        `bson:"app_id"`
	UniqueID string        `bson:"unique_id"`
	SenderID string        `bson:"sender_id"`
	Channels []string      `bson:"channels"`
	Tracking string        `bson:"tracking"`
	CreateAt int64         `bson:"create_at"`
	CloseAt  int64         `bson:"close_at"`
}
