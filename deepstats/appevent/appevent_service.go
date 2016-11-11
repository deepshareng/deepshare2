package appevent

import (
	"github.com/MISingularity/deepshare2/pkg/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// AppEventService serves app events query functionalities.
type AppEventService interface {
	// GetEvents returns all events exist in channls of an app
	// e.g. If client request all event of an appID,
	// the server will retrieve the records for distinct events exist
	// the channels of this app, and it will return the events collection
	// to user.
	// #GetEvents(appid) => AppEvents: [
	//      AppID : appid
	//      Events : [event1, event2, event3...]
	// ]
	GetEvents(appID string) (AppEvents, error)

	InsertEvent(appID, event string) (bool, error)
}

type mongoAppEventService struct {
	coll *mgo.Collection
}

func NewMongoAppEventService(c *mgo.Collection) AppEventService {
	return &mongoAppEventService{c}
}

func (m *mongoAppEventService) GetEvents(appid string) (AppEvents, error) {
	log.Infof("[APPEVENT SERVICE] Get event list of app %s", appid)
	var result []AppEvent
	err := m.coll.Find(bson.M{"app_id": appid}).All(&result)
	appevents := AppEvents{AppID: appid}
	appevents.Events = make([]string, len(result))
	for i, _ := range result {
		appevents.Events[i] = result[i].Event
	}
	log.Debugf("[--Result][APPEVENT SERVICE][GetEvents] AppEvents=%v Err=%v", appevents, err)
	return appevents, err
}

func (m *mongoAppEventService) InsertEvent(appid, event string) (bool, error) {
	log.Infof("[APPEVENT SERVICE] Insert event [%s] to app[%s]", event, appid)
	res, err := m.coll.Upsert(AppEvent{AppID: appid, Event: event}, AppEvent{AppID: appid, Event: event})
	if res.Updated == 1 {
		return false, err
	}
	return true, err
}
