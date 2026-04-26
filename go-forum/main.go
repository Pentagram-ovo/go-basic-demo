package main

import (
	"go-forum/config"
	"go-forum/dao"
	"go-forum/middleware"
	"go-forum/model"
	"go-forum/router"
	"time"

	"github.com/gin-gonic/gin"
)

// 启动异步同步协程：每 30 秒同步一次 Redis 点赞数到 MySQL
func StartSyncTask() {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			<-ticker.C
			dao.SyncLikeRankFromZSet()
		}
	}()
}

func main() {
	config.InitRedis()
	config.InitMysql()
	model.InitPostTable()
	model.InitUserTable()
	model.InitCommentTable()
	r := gin.Default()
	r.Use(middleware.Cors())
	router.SetupRouter(r)
	StartSyncTask()
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
