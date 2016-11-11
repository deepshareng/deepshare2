package storage

import (
	"testing"

	"gopkg.in/redis.v4"
)

func initTestingRedis(t *testing.T) *redisKV {
	cli := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	if err := cli.Ping().Err(); err != nil {
		t.Fatal("Failed to connect to redis, err:", err)
	}
	return &redisKV{cli: cli}
}

func TestRedisSimpleKVSet(t *testing.T) {
	rskv := initTestingRedis(t)
	testSimpleKVSet(rskv, t)
}

func TestRedisSimpleKVDelete(t *testing.T) {
	rskv := initTestingRedis(t)
	testSimpleKVDelete(rskv, t)
}

func TestRedisSimpleKVGetNil(t *testing.T) {
	rskv := initTestingRedis(t)
	testSimpleKVDelete(rskv, t)
	b, err := rskv.Get([]byte("foo"))
	if err != nil {
		t.Error("get with a key not exist, err should be nil")
	}
	if b != nil {
		t.Error("get with a key not exist, value should be")
	}
}

func TestRedisSimpleKVHSet(t *testing.T) {
	rskv := initTestingRedis(t)
	testSimpleKVHSet(rskv, t)
}

func TestRedisSimpleKVHDel(t *testing.T) {
	rskv := initTestingRedis(t)
	testSimpleKVHDel(rskv, t)
}

func TestRedisSimpleKVHIncrBy(t *testing.T) {
	rskv := initTestingRedis(t)
	testSimpleKVHIncrBy(rskv, t)
}

func TestRedisSimpleKVSetEx(t *testing.T) {
	rskv := initTestingRedis(t)
	testSimpleKVSetEx(rskv, t)
}

func TestRedisSimpleKVSAdd(t *testing.T) {
	rskv := initTestingRedis(t)
	testSimpleKVSAdd(rskv, t)
}

func TestRedisSimpleKVSMembers(t *testing.T) {
	rskv := initTestingRedis(t)
	testSimpleKVSMembers(rskv, t)
}

func TestRedisSimpleKVSRem(t *testing.T) {
	rskv := initTestingRedis(t)
	testSimpleKVSRem(rskv, t)
}

func TestRedisSimpleKVSCard(t *testing.T) {
	rskv := initTestingRedis(t)
	testSimpleKVSCard(rskv, t)
}
