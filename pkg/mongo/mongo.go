package mongo

import (
	"github.com/MISingularity/deepshare2/pkg/log"
	"gopkg.in/mgo.v2"
)

func MongoDB(session *mgo.Session, dbName, collNamePrefix, appID string, f func(c *mgo.Collection) error) error {
	if session == nil {
		return nil
	}
	s := session.Clone()
	defer s.Close()
	db := s.DB(dbName)
	collName := collNamePrefix
	if appID != "" {
		collName = collNamePrefix + "_" + appID
	}
	c := db.C(collName)
	if err := f(c); err != nil {
		log.Error("Access mongo error:", err, "db:", dbName, "collNamePrefix:", collNamePrefix, "appID:", appID)
		return err
	}

	return nil
}
