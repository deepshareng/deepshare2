package backup

import (
	"fmt"

	pb "github.com/MISingularity/deepshare2/deepstats/backup/proto"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoRawBackup struct {
	col *mgo.Collection
}

func (l *MongoRawBackup) Insert(event pb.Event) error {
	err := l.col.Insert(event)
	if err != nil {
		return err
	}
	return nil
}

func (l *MongoRawBackup) RetriveAllEvents() ([]pb.Event, error) {
	results := []pb.Event{}
	err := l.col.Find(bson.M{}).All(&results)
	if err != nil {
		return []pb.Event{}, fmt.Errorf("Mongo collection find failed! Err Msg=%v", err)
	}

	return results, nil
}

func NewMongoRawBackupService(col *mgo.Collection) BackupService {
	return &MongoRawBackup{col}
}
