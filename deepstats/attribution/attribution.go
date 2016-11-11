package attribution

import (
	"fmt"

	"encoding/json"

	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/mongo"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	redisKeyPrefix  = "attr:"
	trackingInstall = "install"
	trackingOpen    = "open"
)

// AttributionRetriever retrieves attributions from install/open/close events, provide method to get attribution given appID and uniqueID
type AttributionRetriever interface {
	OnInstall(appID string, uniqueID string, senderID string, channels []string, timestamp int64) error
	OnOpen(appID string, uniqueID string, senderID string, channels []string, timestamp int64) error
	OnClose(appID string, uniqueID string, timestamp int64) error
	GetAttribution(appID string, uniqueID string) *Attribution
}

type attributionRetriever struct {
	mgoSession        *mgo.Session
	mgoDBName         string
	mgoCollNamePrefix string

	db storage.SimpleKV
}

func NewAttributionRetriever(mgoSession *mgo.Session, mgoDBName, mgoCollNamePrefix string, db storage.SimpleKV) AttributionRetriever {
	return &attributionRetriever{
		mgoSession:        mgoSession,
		mgoDBName:         mgoDBName,
		mgoCollNamePrefix: mgoCollNamePrefix,
		db:                db,
	}
}

// OnInstall: insert an Attribution to mongodb
// keep the latest install Attribution in redis for query
// clear the latest open Attribution from redis: after a new install, the previous open should be useless
func (ar *attributionRetriever) OnInstall(appID string, uniqueID string, senderID string, channels []string, timestamp int64) error {
	log.Infof("---OnInstall; appID: %s, uniqueID: %s, senderID: %s, channels: %v, timestamp: %d\n", appID, uniqueID, senderID, channels, timestamp)
	if senderID == "" && channels == nil {
		log.Debug("nothing to bind, ignore!")
		return nil
	}
	attrib := Attribution{
		AppID:    appID,
		UniqueID: uniqueID,
		SenderID: senderID,
		Channels: channels,
		CreateAt: timestamp,
		Tracking: trackingInstall,
	}

	if err := mongo.MongoDB(ar.mgoSession, ar.mgoDBName, ar.mgoCollNamePrefix, appID,
		func(mgocoll *mgo.Collection) error {
			return mgocoll.Insert(attrib)
		}); err != nil {
		log.Error("OnInstall; Failed to insert to mongo, err:", err)
		return err
	}

	return ar.updateDB(&Attribution{
		AppID:    appID,
		UniqueID: uniqueID,
		SenderID: senderID,
		Channels: channels,
		CreateAt: timestamp,
		Tracking: trackingInstall,
	})
}

// OnOpen: insert an Attribution to mongodb
// keep the latest open Attribution in redis for query
func (ar *attributionRetriever) OnOpen(appID string, uniqueID string, senderID string, channels []string, timestamp int64) error {
	log.Infof("---OnOpen; appID: %s, uniqueID: %s, senderID: %s, channels: %v, timestamp: %d\n", appID, uniqueID, senderID, channels, timestamp)
	if senderID == "" && channels == nil {
		log.Debug("nothing to bind, ignore!")
		return nil
	}
	attr := Attribution{
		AppID:    appID,
		UniqueID: uniqueID,
		SenderID: senderID,
		Channels: channels,
		CreateAt: timestamp,
		Tracking: trackingOpen,
	}

	if err := mongo.MongoDB(ar.mgoSession, ar.mgoDBName, ar.mgoCollNamePrefix, appID,
		func(mgocoll *mgo.Collection) error {
			return mgocoll.Insert(attr)
		}); err != nil {
		log.Error("OnOpen; Failed to insert to mongo, err:", err)
		return err
	}

	return ar.updateDB(&Attribution{
		AppID:    appID,
		UniqueID: uniqueID,
		SenderID: senderID,
		Channels: channels,
		CreateAt: timestamp,
		Tracking: trackingOpen,
	})

}

// OnClose: update CloseAt field of the latest open or install Attribution
func (ar *attributionRetriever) OnClose(appID string, uniqueID string, timestamp int64) error {
	log.Infof("---OnClose; appID: %s, uniqueID: %s, timestamp: %d\n", appID, uniqueID, timestamp)
	attr := Attribution{}
	if err := mongo.MongoDB(ar.mgoSession, ar.mgoDBName, ar.mgoCollNamePrefix, appID,
		func(mgocoll *mgo.Collection) error {
			return mgocoll.Find(bson.M{"app_id": appID, "unique_id": uniqueID, "close_at": 0}).Sort("-create_at").One(&attr)
		}); err != nil {
		log.Error("OnClose; Failed to find the latest unclosed session, err:", err)
		return err
	}

	if err := mongo.MongoDB(ar.mgoSession, ar.mgoDBName, ar.mgoCollNamePrefix, appID,
		func(mgocoll *mgo.Collection) error {
			return mgocoll.Update(bson.M{"_id": attr.ID}, bson.M{"$set": bson.M{"close_at": timestamp}})
		}); err != nil {
		log.Error("OnClose; Failed to update close_at field, err:", err)
		return err
	}

	attr.CloseAt = timestamp
	return ar.updateDB(&attr)
}

// GetAttribution: return the latest attribution
func (ar *attributionRetriever) GetAttribution(appID string, uniqueID string) *Attribution {
	k := formKey(appID, uniqueID)
	vInstall, err := ar.db.HGet(k, trackingInstall)
	if err != nil {
		log.Panic(err)
	}
	vOpen, err := ar.db.HGet(k, trackingOpen)
	if err != nil {
		log.Panic(err)
	}

	//open should be newer that install, since we deleted open from redis when receive install
	if vOpen != nil {
		attrOpen := Attribution{}
		if err := json.Unmarshal(vOpen, &attrOpen); err != nil {
			log.Panic(err)
		}
		return &attrOpen
	}

	if vInstall != nil {
		attrInstall := Attribution{}
		if err := json.Unmarshal(vInstall, &attrInstall); err != nil {
			log.Panic(err)
		}
		return &attrInstall
	}

	return nil
}

func formKey(appID, uniqueID string) []byte {
	return []byte(fmt.Sprintf("%s:%s:%s", redisKeyPrefix, appID, uniqueID))
}

func (ar attributionRetriever) updateDB(attr *Attribution) error {
	k := formKey(attr.AppID, attr.UniqueID)
	if attr.Tracking != trackingOpen && attr.Tracking != trackingInstall {
		panic("tracking should be open/install")
	}

	if b, err := json.Marshal(attr); err != nil {
		log.Error("updateDB; Failed to marshal to json")
		return err
	} else {
		if err := ar.db.HSet(k, attr.Tracking, b); err != nil {
			return err
		}
		if attr.Tracking == trackingInstall {
			//TODO 这里假设event是按时间顺序来的,如果假设不成立,需要进行时间戳的对比,不能简单删除
			if err := ar.db.HDel(k, trackingOpen); err != nil {
				return err
			}
		}
	}
	return nil
}
