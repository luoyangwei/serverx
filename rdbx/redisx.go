package rdbx

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

const KeyPrefix = "redis"

type redisConfig struct {
	Addr     string // Addr 链接地址
	DB       int    // DB 数据库, 一般默认是0
	Password string // Password 密码

	// Channel 来控制订阅频道, 连接发送消息会发送到频道
	Channel string // Channel 订阅频道
	// ConnectKey 连接 key, websocket connect 创建时会在 redis 里保存信息,
	// 这里的 ConnectKey 相当于前缀保存在 redis 里的 key
	ConnectKey string
}

func Connect() *redis.Client {
	cfg := readRedisConfig()
	addr := cfg.Addr
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalln(err)
	}
	log.Printf("Redis connected to %s \n", addr)
	return rdb
}

func GetRedisConfig() redisConfig {
	return readRedisConfig()
}

func readRedisConfig() redisConfig {
	var cfg redisConfig
	if err := viper.UnmarshalKey("redis", &cfg); err != nil {
		log.Fatalln(err)
	}
	return cfg
}
