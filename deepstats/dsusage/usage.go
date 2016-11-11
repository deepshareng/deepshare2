package dsusage

import (
	"strings"

	"fmt"

	"time"

	in "github.com/MISingularity/deepshare2/pkg/instrumentation"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"gopkg.in/redis.v4"
)

const (
	RedisKeyUsageFmt = "usage:%s:%s" //usage:<app_id>:<sender_id>
	RedisHKeyInstall = "install"
	RedisHKeyOpen    = "open"
)

type redisAggregateSenderService struct {
	redisCli *redis.ClusterClient
}

type AggregateSenderService interface {
	Insert(*messaging.Event) error
}

func (as *redisAggregateSenderService) Insert(dp *messaging.Event) error {
	if dp.SenderID == "" || dp.UniqueID == "" || dp.SenderID == dp.UniqueID {
		log.Debug("drop the event:", dp)
		return nil
	}
	var err error
	switch {
	case strings.HasSuffix(dp.EventType, "install"):
		err = as.IncrInstall(dp.AppID, dp.SenderID)
	case strings.HasSuffix(dp.EventType, "open"):
		err = as.IncrOpen(dp.AppID, dp.SenderID)
	}
	if err != nil {
		log.Error("Increase usage error:", err)
	}
	return nil
}

func NewAggregateSenderService(redisAddrs []string, password string, poolSize int) AggregateSenderService {
	cli := storage.NewRedisClusterCli(redisAddrs, password, poolSize)
	return &redisAggregateSenderService{redisCli: cli}
}

func (as *redisAggregateSenderService) IncrInstall(appID string, senderID string) error {
	k := fmt.Sprintf(RedisKeyUsageFmt, appID, senderID)
	start := time.Now()
	if err := as.redisCli.HIncrBy(k, RedisHKeyInstall, 1).Err(); err != nil {
		return err
	}
	in.PrometheusForDSUsage.StorageIncDuration(start)
	return nil
}

func (as *redisAggregateSenderService) IncrOpen(appID string, senderID string) error {
	k := fmt.Sprintf(RedisKeyUsageFmt, appID, senderID)
	start := time.Now()
	if err := as.redisCli.HIncrBy(k, RedisHKeyOpen, 1).Err(); err != nil {
		return err
	}
	in.PrometheusForDSUsage.StorageIncDuration(start)
	return nil
}
