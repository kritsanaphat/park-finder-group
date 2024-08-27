package connector

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Client  *redis.Client
	Context context.Context
}

func NewRedisClient(redisURI, address, password string) *Redis {
	var err error = nil
	var redisConfig *redis.Options = &redis.Options{}
	if redisURI != "" {
		redisConfig, err = redis.ParseURL(redisURI)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		redisConfig = &redis.Options{
			Addr:     address,
			Password: password,
			DB:       0,
		}
	}
	rdb := redis.NewClient(redisConfig)
	return &Redis{
		Client:  rdb,
		Context: context.Background(),
	}
}

func (r Redis) Disconnect() error {
	return r.Client.Close()
}

func (r Redis) Ping() *redis.StatusCmd {
	return r.Client.Ping(r.Context)
}

func (r Redis) Set(key string, value interface{}, duration int64) *redis.StatusCmd {
	return r.Client.Set(r.Context, key, value, time.Duration(duration)*time.Second)
}

func (r Redis) Get(key string) *redis.StringCmd {
	return r.Client.Get(r.Context, key)
}
