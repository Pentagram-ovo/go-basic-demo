package service

import (
	"context"
	"go-forum/dao"
)

var ctx = context.Background()

// GetLikeStatus 获取当前用户对帖子的点赞状态
func GetLikeStatus(userID, postID uint) (bool, error) {
	return dao.IsLiked(postID, userID)
}

// GetLikeCount 获取当前帖子的点赞数
func GetLikeCount(postID uint) (int64, error) {
	return dao.CountLikes(postID)
}

// ToggleLikePost 点赞/取消点赞 二合一
func ToggleLikePost(userID, postID uint) error {
	_, err := dao.ToggleLike(postID, userID)
	return err
}
