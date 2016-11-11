package backup

import (
	"fmt"

	pb "github.com/MISingularity/deepshare2/deepstats/backup/proto"
	"github.com/golang/protobuf/proto"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoCompressDayBackup struct {
	col *mgo.Collection
}

func keynameOfMongoCompressDayBackup(event pb.Event) string {
	return convertTime(event.Timestamp)
}

func (l *MongoCompressDayBackup) Insert(event pb.Event) error {
	key := keynameOfMongoCompressDayBackup(event)
	result, err := proto.Marshal(&event)
	if err != nil {
		return err
	}
	_, err = l.col.Upsert(bson.M{"key": key}, bson.M{"$push": bson.M{"val": string(result)}})
	if err != nil {
		return err
	}
	return nil
}

func (l *MongoCompressDayBackup) RetriveAllEvents() ([]pb.Event, error) {
	results := []MongoBackupCompressStorageFormat{}
	err := l.col.Find(bson.M{}).All(&results)
	if err != nil {
		return []pb.Event{}, fmt.Errorf("Mongo collection find failed! Err Msg=%v", err)
	}
	events := []pb.Event{}
	for _, r := range results {
		for _, e := range r.Val {
			event := pb.Event{}
			err = proto.Unmarshal([]byte(e), &event)
			if err != nil {
				return []pb.Event{}, err
			}
			events = append(events, event)
		}
	}
	return events, nil
}

func NewMongoCompressDayBackupService(col *mgo.Collection) BackupService {
	return &MongoCompressDayBackup{col}
}
