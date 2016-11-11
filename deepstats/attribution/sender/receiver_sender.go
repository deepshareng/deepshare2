package sender

import (
	"github.com/MISingularity/deepshare2/pkg/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ReceiverSender struct {
	ReceiverID string `bson:"r"`
	SenderID   string `bson:"s"`
}
type ReceiverSenderMapping struct {
	mgocoll *mgo.Collection
}

func NewReceiverSenderMapping(coll *mgo.Collection) *ReceiverSenderMapping {
	return &ReceiverSenderMapping{coll}
}

func (rsm *ReceiverSenderMapping) OnInstall(receiverID string, senderID string) {
	if _, err := rsm.mgocoll.Upsert(bson.M{"r": receiverID}, bson.M{"$set": bson.M{"s": senderID}}); err != nil {
		log.Error(err)
	}
}

func (rsm *ReceiverSenderMapping) OnOpen(receiverID string, senderID string) {
	if _, err := rsm.mgocoll.Upsert(bson.M{"r": receiverID}, bson.M{"$set": bson.M{"s": senderID}}); err != nil {
		log.Error(err)
	}
}

func (rsm *ReceiverSenderMapping) OnClose(receiverID string) {
	if err := rsm.mgocoll.Remove(bson.M{"r": receiverID}); err != nil {
		log.Error(err)
	}
}

func (rsm *ReceiverSenderMapping) GetSenderID(receiverID string) string {
	result := ReceiverSender{}
	err := rsm.mgocoll.Find(bson.M{"r": receiverID}).One(&result)
	if err != nil {
		log.Error(err)
	}
	return result.SenderID
}
