package config

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func InitRedis() {
	addr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	if os.Getenv("REDIS_PORT") == "" {
		addr = "192.168.18.132:6379"
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr, // 固定
		Password: "",   // 无密码为空
		DB:       0,    // 默认库
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
