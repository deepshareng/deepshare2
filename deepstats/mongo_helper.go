package deepstats

import (
	"github.com/MISingularity/deepshare2/pkg/log"
	"gopkg.in/mgo.v2"
)

const (
	MongoTestingAddr = "127.0.0.1:27017"
)

func MustCreateMongoSession(mongoURL string) *mgo.Session {
	session, err := mgo.Dial(mongoURL)
	if err != nil {
		log.Fatal("mongo server connection failed:", err)
	}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return session
}
