package config

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "192.168.18.130:6379", // 固定
		Password: "",                    // 无密码为空
		DB:       0,                     // 默认库
	})

	// 测试连通
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		panic("redis 连接失败：" + err.Error())
	}
	println("redis 连接成功")
}
