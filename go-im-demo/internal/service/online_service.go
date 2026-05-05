package service

import (
	"context"
	"fmt"
	"go-im-demo/config"
	"time"
)

// 添加用户到在线集合，并设置个体过期
func AddUserOnline(username string) error {
	ctx := context.Background()
	err := config.Rdb.SAdd(ctx, "online_users", username).Err()
	if err != nil {
		return err
	}

	UserKey := fmt.Sprintf("user:%s:online", username)
	err = config.Rdb.Set(ctx, UserKey, 1, time.Second*60).Err()
	if err != nil {
		return err
	}
	return nil
}

// 移除用户
func RemoveUserOnline(username string) error {
	ctx := context.Background()
	err := config.Rdb.SRem(ctx, "online_users", username).Err()
	if err != nil {
		return err
	}
	UserKey := fmt.Sprintf("user:%s:online", username)
	err = config.Rdb.Del(ctx, UserKey).Err()
	if err != nil {
		return err
	}
	return nil
}

// 刷新用户在线状态（心跳时续期）
func RefreshUserOnline(username string) error {
	ctx := context.Background()
	UserKey := fmt.Sprintf("user:%s:online", username)
	err := config.Rdb.Expire(ctx, UserKey, 60*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

// 获取在线用户列表
func GetOnlineUsers() ([]string, error) {
	ctx := context.Background()
	users, err := config.Rdb.SMembers(ctx, "online_users").Result()
	if err != nil {
		return nil, err
	}
	return users, nil
}
