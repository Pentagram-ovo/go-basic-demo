package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Pentagram-ovo/go-basic-demo/gorm-mysql-demo/config"
	"github.com/Pentagram-ovo/go-basic-demo/gorm-mysql-demo/dao"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// 启动异步同步协程：每 10 秒同步一次 Redis 点赞数到 MySQL
func StartSyncTask() {
	go func() {
		ticker := time.NewTicker(10 * time.Second) // 10秒一次
		defer ticker.Stop()

		for {
			<-ticker.C
			dao.SyncLikeCountToDB(context.Background())
		}
	}()
}

func main() {
	config.InitMysql()
	config.InitRedis()
	config.Rdb.FlushAll(context.Background())

	StartSyncTask()
	fmt.Println("开启异步同步，10s更新一次！")

	// 测试缓存
	post, _ := dao.GetPostByIDWithCache(context.Background(), 3)
	fmt.Println("缓存帖子：", post)

	// 测试点赞
	dao.LikePost(context.Background(), 3, 4)
	dao.LikePost(context.Background(), 1, 3)

	// 测试点赞数
	fmt.Println("点赞数：", dao.GetLikeCount(context.Background(), 1))

	// 测试排行
	list, _ := dao.GetTopNPosts(context.Background(), 10)
	fmt.Println("热门帖子：", list)

	select {}
}
