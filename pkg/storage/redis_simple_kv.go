package storage

import (
	"time"

	"github.com/MISingularity/deepshare2/pkg/log"
	"gopkg.in/redis.v4"
)

type redisKV struct {
	cli *redis.Client
}

func NewRedisSimpleKV(url string, password string, poolSize int) SimpleKV {
	cli := redis.NewClient(&redis.Options{
		PoolSize: poolSize,
		Addr:     url,
		Password: password,
	})
	if err := cli.Ping().Err(); err != nil {
		log.Fatal("Failed to connect to redis, err:", err)
	}
	return &redisKV{cli: cli}
}
func NewRedisSentinelSimpleKV(urls []string, masterName string, password string, poolSize int) SimpleKV {
	opt := &redis.FailoverOptions{
		PoolSize:      poolSize,
		MasterName:    masterName,
		SentinelAddrs: urls,
	}
	cli := redis.NewFailoverClient(opt)
	if err := cli.Ping().Err(); err != nil {
		log.Fatal("Connect to redis sentinel master failed, err:", err)
	} else {
		log.Debug("[init]Redis FailoverClient config:", opt)
		if err := cli.Set("foo", "bar", 0).Err(); err != nil {
			log.Fatal("redis set failed:", err)
		} else {
			log.Debug("[init]redis set foo = bar")
		}
		if err := cli.Del("foo").Err(); err != nil {
			log.Fatal("redis del failed:", err)
		} else {
			log.Debug("[init]redis del foo")
		}
	}
	return &redisKV{cli: cli}
}

func (redisKV *redisKV) Get(k []byte) ([]byte, error) {
	b, err := redisKV.cli.Get(string(k)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	log.Debugf("RedisSimpleKV; Get: k = %s; v = %s", string(k), string(b))
	return b, err
}

func (redisKV *redisKV) Delete(k []byte) error {
	return redisKV.cli.Del(string(k)).Err()
}

func (redisKV *redisKV) Set(k []byte, v []byte) error {
	return redisKV.cli.Set(string(k), v, 0).Err()
}

func (redisKV *redisKV) SetEx(k []byte, v []byte, expiration time.Duration) error {
	log.Debugf("RedisSimpleKV; SetEx: k = %s, v = %s, expiration = %v\n", string(k), string(v), expiration)
	return redisKV.cli.Set(string(k), v, expiration).Err()
}

func (redisKV *redisKV) HSet(k []byte, hk string, v []byte) error {
	return redisKV.cli.HSet(string(k), hk, string(v)).Err()
}

func (redisKV *redisKV) HGet(k []byte, hk string) ([]byte, error) {
	b, err := redisKV.cli.HGet(string(k), hk).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return b, err
}

func (redisKV *redisKV) HGetAll(k []byte) (map[string]string, error) {
	v, err := redisKV.cli.HGetAll(string(k)).Result()
	if err == redis.Nil {
		return nil, nil
	}
	return v, err
}

func (redisKV *redisKV) HDel(k []byte, hk string) error {
	return redisKV.cli.HDel(string(k), hk).Err()
}

func (redisKV *redisKV) HIncrBy(k []byte, hk string, n int) error {
	return redisKV.cli.HIncrBy(string(k), hk, int64(n)).Err()
}

func (redisKV *redisKV) Exists(k []byte) bool {
	return redisKV.cli.Exists(string(k)).Val()
}

func (redisKV *redisKV) SAdd(k []byte, v string) error {
	return redisKV.cli.SAdd(string(k), v).Err()
}

func (redisKV *redisKV) SRem(k []byte, v string) error {
	return redisKV.cli.SRem(string(k), v).Err()
}

func (redisKV *redisKV) SCard(k []byte) (int64, error) {
	return redisKV.cli.SCard(string(k)).Result()
}

func (redisKV *redisKV) SMembers(k []byte) ([]string, error) {
	return redisKV.cli.SMembers(string(k)).Result()
}
