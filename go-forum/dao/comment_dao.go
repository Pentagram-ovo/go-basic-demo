package dao

import (
	"go-forum/config"
	"go-forum/model"
)

// CreateComment 创建评论
func CreateComment(comment *model.Comment) error {
	return config.DB.Create(comment).Error
}

// GetCommentById 根据评论id来查询评论
func GetCommentById(id uint) (*model.Comment, error) {
	var comment model.Comment
	err := config.DB.Where("id=?", id).First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetCommentListByPostID 根据帖子id来查看评论列表
func GetCommentListByPostID(postid uint, page, size int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64
	offset := (page - 1) * size
	config.DB.Model(&model.Comment{}).Where("post_id = ?", postid).Count(&total)
	err := config.DB.Order("created_at DESC").Offset(offset).Limit(size).Where("post_id=?", postid).Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}

// GetCommentByUserId 根据用户id来查看其发表的所有评论
func GetCommentByUserId(userid uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := config.DB.Where("user_id=?", userid).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}
