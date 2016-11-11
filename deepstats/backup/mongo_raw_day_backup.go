package backup

import (
	"fmt"
	"time"

	"regexp"

	pb "github.com/MISingularity/deepshare2/deepstats/backup/proto"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoRawDayBackup struct {
	session      *mgo.Session
	dbNamePrefix string
}

func (d *MongoRawDayBackup) extractMongoColl(event pb.Event) *mgo.Collection {
	t := time.Unix(event.Timestamp, 0)
	dbName := fmt.Sprintf("%s_%04d%02d", d.dbNamePrefix, t.Year(), t.Month())
	db := d.session.DB(dbName)
	collName := fmt.Sprintf("coll%02d", t.Day())
	return db.C(collName)
}

func (d *MongoRawDayBackup) Insert(event pb.Event) error {
	col := d.extractMongoColl(event)
	return col.Insert(event)
}

func (d *MongoRawDayBackup) RetriveAllEvents() ([]pb.Event, error) {
	dbs, err := d.session.DatabaseNames()
	if err != nil {
		return []pb.Event{}, err
	}
	result := make([]pb.Event, 0)
	regDB := regexp.MustCompile(fmt.Sprintf("%s_\\d{6}", d.dbNamePrefix))
	regColl := regexp.MustCompile(`coll\d{2}`)
	for _, dbName := range dbs {
		if !regDB.Match([]byte(dbName)) {
			continue
		}
		db := d.session.DB(dbName)
		colls, err := db.CollectionNames()
		if err != nil {
			return []pb.Event{}, err
		}
		for _, v := range colls {
			if !regColl.Match([]byte(v)) {
				continue
			}
			coll := db.C(v)
			res := make([]pb.Event, 0)
			err := coll.Find(bson.M{}).All(&res)
			if err != nil {
				return []pb.Event{}, err
			}
			result = append(result, res...)
		}
	}

	return result, nil
}

func NewMongoRawDayBackupService(session *mgo.Session, dbNamePrefix string) BackupService {
	return &MongoRawDayBackup{
		session:      session,
		dbNamePrefix: dbNamePrefix,
	}
}
