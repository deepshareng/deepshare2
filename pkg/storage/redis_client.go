package storage

import (
	"github.com/MISingularity/deepshare2/pkg/log"
	"gopkg.in/redis.v4"
)

func NewRedisClient(url string, password string) *redis.Client {
	cli := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
	})
	if err := cli.Ping().Err(); err != nil {
		log.Fatal("Failed to connect to redis, err:", err)
	}
	return cli
}
