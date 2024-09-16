package rds

import "github.com/redis/go-redis/v9"

var RDB *redis.Client

func CreateClient() (*redis.Client, error) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "172.16.44.113:6379",
		Password: "",
		DB:       0,
	})

	return RDB, nil
}
