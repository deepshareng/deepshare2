package retention

import (
	"strconv"
	"strings"
	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg/condition"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"gopkg.in/redis.v4"
)

type RetentionService interface {
	FindMatchUserForRetention(dp *messaging.Event) error

	InsertInstallEventForRetention(dp *messaging.Event) error
}

type RedisRetentionService struct {
	day          int
	name         string
	cli          *redis.Client
	produceTopic []byte
	producer     messaging.Producer
}

func NewRedisRetentionService(day int, name string, cli *redis.Client, produceTopic []byte, producer messaging.Producer) *RedisRetentionService {
	return &RedisRetentionService{day, name, cli, produceTopic, producer}
}

func (r *RedisRetentionService) getEventKey(t time.Time, dp *messaging.Event) string {
	year := strconv.Itoa(t.Year())
	day := strconv.Itoa(t.Day())
	return r.name + "#" + year + t.Month().String() + day + "#" + dp.AppID + "#" + dp.UniqueID
}

func (r *RedisRetentionService) FindMatchUserForRetention(dp *messaging.Event) error {
	tNow := condition.ConvertTimeToDay(time.Unix(dp.TimeStamp, 0))
	tBegin := tNow.Add(-time.Hour * 24 * time.Duration(r.day))
	key := r.getEventKey(tBegin, dp)

	if !r.cli.Exists(key).Val() {
		return nil
	}
	val := r.cli.Get(key).Val()
	log.Infof("[Retention Service][%s] For request=%v, find matching install event=%v", r.name, dp, key)
	err := r.cli.Del(key).Err()
	if err != nil {
		return err
	}
	r.produceRetentionBackToNSQ(dp.AppID, dp.UniqueID, r.name, strings.Split(val, "#"), tNow.Unix(), dp.UAInfo)
	return nil
}

func (r *RedisRetentionService) InsertInstallEventForRetention(dp *messaging.Event) error {
	if condition.IsEmptyChannel(dp.Channels) {
		dp.Channels = []string{}
	}
	t := condition.ConvertTimeToDay(time.Unix(dp.TimeStamp, 0))
	if dp.AppID == "" || dp.UniqueID == "" {
		return nil
	}
	key := r.getEventKey(t, dp)
	log.Infof("[Retention Service][%s] Insert an install event, key=%s", r.name, key)
	if r.cli.Exists(key).Val() {
		return nil
	}
	r.produceRetentionBackToNSQ(dp.AppID, dp.UniqueID, r.name+"_install", dp.Channels, t.Add(time.Hour*24*time.Duration(r.day)).Unix(), dp.UAInfo)
	return r.cli.Set(key, strings.Join(dp.Channels, "#"), time.Hour*24*time.Duration(r.day+1)).Err()
}

func (r *RedisRetentionService) produceRetentionBackToNSQ(appid, uniqueid, eventname string, channels []string, t int64, ua uainfo.UAInfo) {
	if r.producer == nil {
		return
	}
	e := &messaging.Event{
		AppID:     appid,
		EventType: api.RetentionPrefix + eventname,
		UniqueID:  uniqueid,
		Channels:  channels,
		SenderID:  "",
		CookieID:  "",
		Count:     1,
		UAInfo:    ua,
		KVs:       map[string]interface{}{},
		TimeStamp: t,
	}
	r.producer.Produce(r.produceTopic, e)
}
