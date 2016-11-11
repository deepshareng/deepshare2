package testutil

import "gopkg.in/redis.v4"

func MustNewRedisClient(password string) *redis.Client {
	cli := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: password,
	})
	return cli
}
