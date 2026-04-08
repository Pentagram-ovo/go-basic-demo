package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/Pentagram-ovo/go-basic-demo/gorm-mysql-demo/config"
	"github.com/Pentagram-ovo/go-basic-demo/gorm-mysql-demo/model"
	"github.com/redis/go-redis/v9"

	"time"
)

// 创建帖子
func CreatePost(post *model.Post) error {
	return config.DB.Create(post).Error
}

// 根据id查询帖子
func GetPostByID(id uint) (*model.Post, error) {
	var post model.Post
	err := config.DB.Where("id = ?", id).First(&post).Error
	return &post, err
}

// 根据id修改帖子标题
func UpdatePostTitle(id uint, title string) error {
	return config.DB.Model(&model.Post{}).Where("id = ?", id).Update("title", title).Error
}

// 根据id修改帖子内容
func UpdatePostContent(id uint, content string) error {
	return config.DB.Model(&model.Post{}).Where("id = ?", id).Update("content", content).Error
}

// 根据id删除帖子
func DeletePost(id uint) error {
	return config.DB.Delete(&model.Post{}, id).Error
}

// 带缓存的帖子详情
func GetPostByIDWithCache(ctx context.Context, id uint) (*model.Post, error) {
	key := fmt.Sprintf("post:%d", id)
	cacheData, err := config.Rdb.Get(ctx, key).Result()
	//命中
	if err == nil {
		var post model.Post
		err = json.Unmarshal([]byte(cacheData), &post)
		return &post, err
	}

	post, err := GetPostByID(id)
	if err != nil {
		return nil, err
	}
	//序列化后存入 Redis，15分钟过期
	postBytes, _ := json.Marshal(post)
	config.Rdb.Set(ctx, key, postBytes, 15*time.Minute)
	return post, nil
}

// LikePost 点赞：去重+计数
func LikePost(ctx context.Context, postID uint, userID uint) error {
	keyset := fmt.Sprintf("like:users:%d", postID)   // 谁点过赞
	keyCount := fmt.Sprintf("like:count:%d", postID) // 点赞数
	exists, _ := config.Rdb.SIsMember(ctx, keyset, userID).Result()
	if exists {
		return errors.New("请勿重复点赞！")
	}

	config.Rdb.SAdd(ctx, keyset, userID)
	config.Rdb.Incr(ctx, keyCount)
	UpdatePostRank(ctx, postID)
	return nil
}

// GetLikeCount 获取当前点赞数
func GetLikeCount(ctx context.Context, postID uint) int {
	keyset := fmt.Sprintf("like:count:%d", postID)
	countStr, err := config.Rdb.Get(ctx, keyset).Result()
	if err != nil {
		return 0
	}
	count, _ := strconv.Atoi(countStr)
	return count
}

// UpdatePostRank 更新帖子热度（分数=点赞数）
func UpdatePostRank(ctx context.Context, postID uint) {
	key := "post:rank"
	score := float64(GetLikeCount(ctx, postID))
	// ZSet 固定写法
	config.Rdb.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: postID,
	})
}

// GetTopNPosts 获取前 N 热门帖子
func GetTopNPosts(ctx context.Context, top int64) ([]redis.Z, error) {
	key := "post:rank"
	return config.Rdb.ZRevRangeWithScores(ctx, key, 0, top-1).Result()
}

// 把 Redis 中的点赞数 同步到 MySQL
func SyncLikeCountToDB(ctx context.Context) {
	// 1. 获取所有有过点赞的帖子 ID（从 ZSet 排行榜里拿）
	postIDs, err := config.Rdb.ZRange(ctx, "post:rank", 0, -1).Result()
	if err != nil {
		return
	}

	// 2. 遍历每个帖子，同步 like_count
	for _, pidStr := range postIDs {
		postID, _ := strconv.Atoi(pidStr)
		keyCount := fmt.Sprintf("like:count:%d", postID)

		// 从 Redis 拿最新点赞数
		likeCount, err := config.Rdb.Get(ctx, keyCount).Int()
		if err != nil {
			continue
		}

		// 同步到 MySQL
		config.DB.Model(&model.Post{}).
			Where("id = ?", postID).
			Update("like_count", likeCount)
	}

	fmt.Println("✅ 异步同步完成：Redis 点赞数 → MySQL")
}
