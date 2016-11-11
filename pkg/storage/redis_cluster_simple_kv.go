package storage

import (
	"time"

	"strings"

	"github.com/MISingularity/deepshare2/pkg/log"
	"gopkg.in/redis.v4"
)

type redisClusterKV struct {
	cli *redis.ClusterClient
}

func NewRedisClusterSimpleKV(urls []string, password string, poolSize int) SimpleKV {
	return &redisClusterKV{cli: NewRedisClusterCli(urls, password, poolSize)}
}

func NewRedisClusterCli(urls []string, password string, poolSize int) *redis.ClusterClient {
	opt := &redis.ClusterOptions{
		Addrs:    urls,
		Password: password,
		PoolSize: poolSize,
	}
	cli := redis.NewClusterClient(opt)

	if err := cli.Ping().Err(); err != nil {
		log.Fatal("Connect to redis cluster failed, err:", err)
	} else {
		log.Debug("[init]Redis Cluster config:", opt)
		if err := cli.Set("foo", "bar", 0).Err(); err != nil {
			log.Fatal("redis custer set failed:", err)
		} else {
			log.Debug("[init]redis cluster set foo = bar")
		}
		if err := cli.Del("foo").Err(); err != nil {
			log.Fatal("redis cluster del failed:", err)
		} else {
			log.Debug("[init]redis cluster del foo")
		}
	}
	return cli
}

func (redisKV *redisClusterKV) Get(k []byte) ([]byte, error) {
	b, err := redisKV.cli.Get(string(k)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil && strings.Contains(err.Error(), "CLUSTERDOWN") {
		log.Fatal("RedisClusterKV; Get; cluster is down! err:", err)
	}
	log.Debugf("RedisClusterKV; Get: k = %s; v = %s", string(k), string(b))
	return b, err
}

func (redisKV *redisClusterKV) Delete(k []byte) error {
	err := redisKV.cli.Del(string(k)).Err()
	if err != nil && strings.Contains(err.Error(), "CLUSTERDOWN") {
		log.Fatal("RedisClusterKV; Delete; cluster is down! err:", err)
	}
	return err
}

func (redisKV *redisClusterKV) Set(k []byte, v []byte) error {
	err := redisKV.cli.Set(string(k), v, 0).Err()
	if err != nil && strings.Contains(err.Error(), "CLUSTERDOWN") {
		log.Fatal("RedisClusterKV; Set; cluster is down! err:", err)
	}
	return err
}

func (redisKV *redisClusterKV) SetEx(k []byte, v []byte, expiration time.Duration) error {
	log.Debugf("RedisClusterKV; SetEx: k = %s, v = %s, expiration = %v\n", string(k), string(v), expiration)
	err := redisKV.cli.Set(string(k), v, expiration).Err()
	if err != nil && strings.Contains(err.Error(), "CLUSTERDOWN") {
		log.Fatal("RedisClusterKV; SetEx; cluster is down! err:", err)
	}
	return err
}

func (redisKV *redisClusterKV) HSet(k []byte, hk string, v []byte) error {
	err := redisKV.cli.HSet(string(k), hk, string(v)).Err()
	if err != nil && strings.Contains(err.Error(), "CLUSTERDOWN") {
		log.Fatal("RedisClusterKV; HSet; cluster is down! err:", err)
	}
	return err
}

func (redisKV *redisClusterKV) HGet(k []byte, hk string) ([]byte, error) {
	b, err := redisKV.cli.HGet(string(k), hk).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil && strings.Contains(err.Error(), "CLUSTERDOWN") {
		log.Fatal("RedisClusterKV; HGet; cluster is down! err:", err)
	}
	return b, err
}

func (redisKV *redisClusterKV) HGetAll(k []byte) (map[string]string, error) {
	v, err := redisKV.cli.HGetAll(string(k)).Result()
	if err == redis.Nil {
		return nil, nil
	}
	return v, err
}

func (redisKV *redisClusterKV) HDel(k []byte, hk string) error {
	err := redisKV.cli.HDel(string(k), hk).Err()
	if err != nil && strings.Contains(err.Error(), "CLUSTERDOWN") {
		log.Fatal("RedisClusterKV; HDel; cluster is down! err:", err)
	}
	return err
}

func (redisKV *redisClusterKV) HIncrBy(k []byte, hk string, n int) error {
	err := redisKV.cli.HIncrBy(string(k), hk, int64(n)).Err()
	if err != nil && strings.Contains(err.Error(), "CLUSTERDOWN") {
		log.Fatal("RedisClusterKV; HIncrBy; cluster is down! err:", err)
	}
	return err
}

func (redisKV *redisClusterKV) Exists(k []byte) bool {
	v := redisKV.cli.Exists(string(k))
	err := v.Err()
	if err != nil && strings.Contains(err.Error(), "CLUSTERDOWN") {
		log.Fatal("RedisClusterKV; Exists; cluster is down! err:", err)
	}
	return v.Val()
}

func (redisKV *redisClusterKV) SAdd(k []byte, v string) error {
	err := redisKV.cli.SAdd(string(k), v).Err()
	if err != nil && strings.Contains(err.Error(), "CLUSTERDOWN") {
		log.Fatal("RedisClusterKV; SAdd; cluster is down! err:", err)
	}
	return err
}

func (redisKV *redisClusterKV) SRem(k []byte, v string) error {
	err := redisKV.cli.SRem(string(k), v).Err()
	if err != nil && strings.Contains(err.Error(), "CLUSTERDOWN") {
		log.Fatal("RedisClusterKV; SRem; cluster is down! err:", err)
	}
	return err
}

func (redisKV *redisClusterKV) SCard(k []byte) (int64, error) {
	v, err := redisKV.cli.SCard(string(k)).Result()
	if err != nil && strings.Contains(err.Error(), "CLUSTERDOWN") {
		log.Fatal("RedisClusterKV; SCard; cluster is down! err:", err)
	}
	return v, err
}

func (redisKV *redisClusterKV) SMembers(k []byte) ([]string, error) {
	s, err := redisKV.cli.SMembers(string(k)).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil && strings.Contains(err.Error(), "CLUSTERDOWN") {
		log.Fatal("RedisClusterKV; SCard; cluster is down! err:", err)
	}
	return s, err
}
