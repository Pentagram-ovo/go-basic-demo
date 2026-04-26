package dao

import (
	"fmt"
	"go-forum/config"
	"go-forum/model"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// ZAddHotPost 将数据载入排行榜
func ZAddHotPost(postID uint) error {
	key := "hot:post:rank"
	score, err := CountLikes(postID)
	if err != nil {
		return err
	}
	_, err = config.Rdb.ZAdd(ctx, key, redis.Z{
		Member: postID,
		Score:  float64(score * 2),
	}).Result()
	if err != nil {
		return err
	}
	return nil
}

// SyncLikeRankFromZSet 更新点赞数量
func SyncLikeRankFromZSet() {
	postIDs, err := config.Rdb.ZRange(ctx, "hot:post:rank", 0, -1).Result()
	if err != nil {
		return
	}
	for _, postIDStr := range postIDs {
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			continue
		}
		key := fmt.Sprintf("like:count:%d", postID)

		//从redis拿到最新点赞数
		likecount, err := config.Rdb.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		config.DB.Model(&model.Post{}).Where("id=?", postID).Update("like_count", likecount)
		if err := ZAddHotPost(uint(postID)); err != nil {
			fmt.Printf("同步热榜失败 postID=%d, err=%v\n", postID, err)
		}
	}
	fmt.Println("✅ 异步同步完成：Redis 点赞数 → MySQL")
}
