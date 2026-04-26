package dao

import (
	"context"
	"errors"
	"fmt"
	"go-forum/config"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// 定义 Lua 脚本（Go 变量）
var likeScript = redis.NewScript(`
if redis.call('SISMEMBER', KEYS[1], ARGV[1]) == 1 then
    return 0
end
redis.call('SADD', KEYS[1], ARGV[1])
local new_count = redis.call('INCR', KEYS[2])
redis.call('ZADD', KEYS[3], new_count * 2, ARGV[2])
return 1
`)

var unlikeScript = redis.NewScript(`
if redis.call('SISMEMBER', KEYS[1], ARGV[1]) == 0 then
    return 0
end
redis.call('SREM', KEYS[1], ARGV[1])
local new_count = redis.call('DECR', KEYS[2])
if new_count < 0 then
    redis.call('SET', KEYS[2], 0)
    new_count = 0
end
redis.call('ZADD', KEYS[3], new_count * 2, ARGV[2])
return 1
`)

// 用户点赞
//存在并发安全问题
//func LikePost(postID, userID uint) error {
//	keyset := fmt.Sprintf("like:users:%d", postID)   //谁点过赞
//	keycount := fmt.Sprintf("like:count:%d", postID) //帖子的总点赞数
//	exists, _ := IsLiked(postID, userID)
//	if exists {
//		return errors.New("请勿重复点赞！")
//	}
//	config.Rdb.SAdd(ctx, keyset, userID) //计入点赞的用户
//	config.Rdb.Incr(ctx, keycount)       //点赞+1
//	return ZAddHotPost(postID)           //更新热度排行
//}

// 用户取消点赞
//存在并发安全问题
//func UnlikePost(postID, userID uint) error {
//	keyset := fmt.Sprintf("like:users:%d", postID)
//	keycount := fmt.Sprintf("like:count:%d", postID)
//	exists, _ := IsLiked(postID, userID)
//	if !exists {
//		return errors.New("用户还未点赞，无法取消！")
//	}
//	config.Rdb.SRem(ctx, keyset, userID) //删除点赞的用户
//	config.Rdb.Decr(ctx, keycount)       //点赞-1
//	return ZAddHotPost(postID)           //更新热度排行
//}

// LikePost 原子点赞
func LikePost(postID, userID uint) error {
	keys := []string{
		fmt.Sprintf("like:users:%d", postID), //KEYS[1]
		fmt.Sprintf("like:count:%d", postID), //KEYS[2]
		"hot:post:rank",                      //KEYS[3]
	}
	res, err := likeScript.Run(ctx, config.Rdb, keys, userID, postID).Int()
	//keys后是ARGV[1],ARGV[2]
	if err != nil {
		return err
	}
	if res == 0 {
		return errors.New("请勿重复点赞！")
	}
	return nil
}

// UnlikePost 原子取消点赞
func UnlikePost(postID, userID uint) error {
	keys := []string{
		fmt.Sprintf("like:users:%d", postID),
		fmt.Sprintf("like:count:%d", postID),
		"hot:post:rank",
	}
	res, err := unlikeScript.Run(ctx, config.Rdb, keys, userID, postID).Int()
	if err != nil {
		return err
	}
	if res == 0 {
		return errors.New("用户还未点赞，无法取消！")
	}
	return nil
}

// IsLiked 用户是否点过赞
func IsLiked(postID, userID uint) (bool, error) {
	keyset := fmt.Sprintf("like:users:%d", postID)
	exists, err := config.Rdb.SIsMember(ctx, keyset, userID).Result()
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

// CountLikes 获取当前帖子点赞总数
func CountLikes(postID uint) (int64, error) {
	key := fmt.Sprintf("like:count:%d", postID)
	val, err := config.Rdb.Get(ctx, key).Int64()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	return val, err
}

var toggleScript = redis.NewScript(`
local is_liked = redis.call('SISMEMBER', KEYS[1], ARGV[1])
if is_liked == 1 then
    redis.call('SREM', KEYS[1], ARGV[1])
    local new_count = redis.call('DECR', KEYS[2])
    if new_count < 0 then
        redis.call('SET', KEYS[2], 0)
        new_count = 0
    end
    redis.call('ZADD', KEYS[3], new_count * 2, ARGV[2])
    return -1
else
    redis.call('SADD', KEYS[1], ARGV[1])
    local new_count = redis.call('INCR', KEYS[2])
    redis.call('ZADD', KEYS[3], new_count * 2, ARGV[2])
    return 1
end
`)

func ToggleLike(postID, userID uint) (int, error) {
	keys := []string{
		fmt.Sprintf("like:users:%d", postID),
		fmt.Sprintf("like:count:%d", postID),
		"hot:post:rank",
	}
	return toggleScript.Run(ctx, config.Rdb, keys, userID, postID).Int()
}
