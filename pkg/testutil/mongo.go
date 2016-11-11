package testutil

import (
	"time"

	"gopkg.in/mgo.v2"
)

func MustNewMongoColl(dbName, collName string) *mgo.Collection {
	s, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	db := s.DB(dbName)
	if err := db.DropDatabase(); err != nil {
		panic(err)
	}
	coll := db.C(collName)
	return coll
}

func MustNewLocalSession() *mgo.Session {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:   []string{"127.0.0.1"},
		Timeout: 20 * time.Second,
	}
	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		panic("Mongo Dial failed")
	}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return session
}
