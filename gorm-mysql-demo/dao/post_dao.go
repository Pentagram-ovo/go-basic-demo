package dao

import (
	"github.com/Pentagram-ovo/go-basic-demo/gorm-mysql-demo/config"
	"github.com/Pentagram-ovo/go-basic-demo/gorm-mysql-demo/model"
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
