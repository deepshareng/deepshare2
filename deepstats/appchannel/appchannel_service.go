package appchannel

import (
	"github.com/MISingularity/deepshare2/pkg/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// AppChannelsService serves app query functionalities.
type AppChannelService interface {
	// GetChannels returns all channels of an app
	// e.g. If client request all channels of an appID,
	// the server will retrieve the records and find distinct channel
	// of this app, and it will return the channelID array to client.
	// #GetChannels(appid) => AppChannels: [
	// 		AppID : appid
	//		ChannelIDs : [channelid1, channelid2...]
	// ]
	GetChannels(appID string) (AppChannels, error)

	InsertChannel(appID, channel string) (bool, error)
	DeleteChannel(appID, channel string) error
}

type mongoAppChannelService struct {
	coll *mgo.Collection
}

func NewMongoAppChannelService(c *mgo.Collection) AppChannelService {
	return &mongoAppChannelService{c}
}

func (m *mongoAppChannelService) GetChannels(appid string) (AppChannels, error) {
	log.Infof("[APPCHANNEL SERVICE] Get channel list of app %s", appid)
	var result []AppChannel
	err := m.coll.Find(bson.M{"app_id": appid}).Sort("channel").All(&result)
	appchannels := AppChannels{AppID: appid}
	appchannels.Channels = make([]string, len(result))
	for i, _ := range result {
		appchannels.Channels[i] = result[i].Channel
	}
	log.Debugf("[--Result][APPCHANNEL SERVICE][GetChannels] AppChannels=%v Err=%v", appchannels, err)
	return appchannels, err
}

func (m *mongoAppChannelService) InsertChannel(appid, channel string) (bool, error) {
	log.Infof("[APPCHANNEL SERVICE] Insert channel[%s] to app[%s]", channel, appid)
	res, err := m.coll.Upsert(AppChannel{AppID: appid, Channel: channel}, AppChannel{AppID: appid, Channel: channel})
	if res.Updated == 1 {
		return false, err
	}
	return true, err
}

func (m *mongoAppChannelService) DeleteChannel(appid, channel string) error {
	log.Infof("[APPCHANNEL SERVICE] Delete channel[%s] of app[%s]", channel, appid)
	return m.coll.Remove(AppChannel{AppID: appid, Channel: channel})
}
