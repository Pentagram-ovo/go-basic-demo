package dao

import (
	"go-forum/config"
	"go-forum/model"
)

// CreatePost 创建帖子
func CreatePost(post *model.Post) error {
	return config.DB.Create(post).Error
}

// GetPostByUserId 根据用户id来查询帖子
func GetPostByUserId(userid uint) ([]model.Post, error) {
	var post []model.Post
	err := config.DB.Where("user_id=?", userid).Find(&post).Error
	if err != nil {
		return nil, err
	}
	return post, nil
}

// GetPostByPostId 根据帖子id来查询帖子
func GetPostByPostId(id uint) (*model.Post, error) {
	var post model.Post
	err := config.DB.Where("id=?", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPostList page: 页码  size: 每页条数
// 返回 帖子列表、总数、错误
func GetPostList(page int, size int) ([]model.Post, int64, error) {
	var post []model.Post
	var total int64
	offset := (page - 1) * size
	config.DB.Model(&model.Post{}).Count(&total)
	err := config.DB.Order("created_at DESC").Offset(offset).Limit(size).Find(&post).Error
	if err != nil {
		return nil, 0, err
	}
	return post, total, nil
}

// UpdatePost 更新帖子内容
func UpdatePost(postid uint, title, content string) error {
	return config.DB.Model(&model.Post{}).Where("id = ?", postid).Updates(model.Post{
		Title:   title,
		Content: content,
	}).Error
}

// DeletePost 删除帖子内容
func DeletePost(id uint) error {
	return config.DB.Where("id=?", id).Delete(&model.Post{}).Error
}

// DeleteCommentsByPostID 根据帖子ID删除所有评论
func DeleteCommentsByPostID(postID uint) error {
	return config.DB.Where("post_id = ?", postID).Delete(&model.Comment{}).Error
}

// GetPostListByIDs 根据给出的帖子id列查询相对的帖子详情
func GetPostListByIDs(ids []uint) ([]model.Post, error) {
	var post []model.Post
	err := config.DB.Where("id IN (?)", ids).Find(&post).Error
	if err != nil {
		return nil, err
	}
	return post, nil
}
