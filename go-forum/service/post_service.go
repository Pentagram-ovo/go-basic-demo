package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-forum/config"
	"go-forum/dao"
	"go-forum/model"
	"strconv"
	"time"
)

// CreatePost 发布新帖子
func CreatePost(userid uint, title, content string) error {
	var post model.Post
	post.Title = title
	post.Content = content
	post.UserID = userid
	return dao.CreatePost(&post)
}

// GetPost 根据帖子id获取帖子信息(未使用redis缓存)
func GetPost(id uint) (*model.Post, error) {
	var post *model.Post
	post, err := dao.GetPostByPostId(id)
	if err != nil {
		return nil, err
	}
	return post, err
}

// GetPostListService 分页查询
func GetPostListService(page int, size int) ([]model.Post, int64, error) {
	if page < 1 || size > 20 {
		return nil, 0, errors.New("请调整参数")
	}
	post, total, err := dao.GetPostList(page, size)
	if err != nil {
		return nil, 0, err
	}
	return post, total, err
}

// GetPostCache 获取缓存并存入(根据帖子id）
func GetPostCache(id uint) (*model.Post, error) {
	key := fmt.Sprintf("post:%d", id)
	cacheData, err := config.Rdb.Get(ctx, key).Result()
	if err == nil { //命中
		var post model.Post
		if err := json.Unmarshal([]byte(cacheData), &post); err == nil {
			return &post, nil
		}
		// 反序列化失败，删除脏缓存，回源查库
		config.Rdb.Del(ctx, key)
	}

	post, err := dao.GetPostByPostId(id)
	if err != nil {
		return nil, err
	}
	//序列化后存入 Redis，15分钟过期
	postBytes, _ := json.Marshal(post)
	config.Rdb.Set(ctx, key, postBytes, 15*time.Minute)
	return post, nil
}

// DelPostCache 删除缓存
func DelPostCache(id uint) error {
	// 从热榜移除
	if _, err := config.Rdb.ZRem(ctx, "hot:post:rank", id).Result(); err != nil {
		return err
	}
	// 删除帖子详情缓存
	if _, err := config.Rdb.Del(ctx, fmt.Sprintf("post:%d", id)).Result(); err != nil {
		return err
	}
	return nil
}

func UpdatePost(postid, userid uint, content, title string) error {
	post, err := dao.GetPostByPostId(postid)
	if err != nil {
		return err
	}
	if post.UserID != userid {
		return errors.New("无权限")
	}

	if title == "" || content == "" {
		return errors.New("参数不能为空")
	}

	if err = dao.UpdatePost(postid, title, content); err != nil {
		return err
	}
	// 更新后清理相关缓存
	if err = DelPostCache(postid); err != nil {
		return err
	}
	return nil
}

func DeletePost(id uint) error {
	// 1. 确认帖子存在
	_, err := dao.GetPostByPostId(id)
	if err != nil {
		return err
	}

	// 2. 删除该帖子的所有评论
	if err := dao.DeleteCommentsByPostID(id); err != nil {
		return err
	}

	// 3. 删除帖子本身
	if err := dao.DeletePost(id); err != nil {
		return err
	}

	// 4. 清理缓存和排行
	if err := DelPostCache(id); err != nil {
		return err
	}

	// 5. 删除 Redis 点赞相关键
	config.Rdb.Del(ctx, fmt.Sprintf("like:users:%d", id))
	config.Rdb.Del(ctx, fmt.Sprintf("like:count:%d", id))
	return nil
}

// GetTopNPosts 获取前 N 热门帖子
func GetTopNPosts(top int64) ([]model.Post, error) {
	key := "hot:post:rank"
	zlist, err := config.Rdb.ZRevRangeWithScores(ctx, key, 0, top-1).Result()
	if err != nil {
		return nil, err
	}
	var ids []uint
	for _, z := range zlist {
		// z.Member 是字符串类型的 "1","2"
		idStr := z.Member.(string)
		id, _ := strconv.ParseUint(idStr, 10, 64)
		ids = append(ids, uint(id))
	}
	return dao.GetPostListByIDs(ids)
}
