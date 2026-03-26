package dao

import (
	"fmt"

	"github.com/Pentagram-ovo/go-basic-demo/gorm-mysql-demo/config"
	"github.com/Pentagram-ovo/go-basic-demo/gorm-mysql-demo/model"
)

// 创建用户
func CreateUser(user *model.User) error {
	return config.DB.Create(user).Error
}

// 根据id查询用户
func GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := config.DB.Where("id = ?", id).First(&user).Error
	return &user, err
}

// 根据id更新用户
func UpdateUsername(id uint, name string) error {
	return config.DB.Model(&model.User{}).Where("id = ?", id).Update("name", name).Error
}

// 根据id删除用户
func DeleteUser(id uint) error {
	return config.DB.Delete(&model.User{}, id).Error
}

// 查询用户以及他的全部帖子
func GetUserWithPosts(uid uint) (*model.User, error) {
	var user model.User
	err := config.DB.Preload("Post").First(&user).Error
	return &user, err
}

// 查询帖子以及贴主信息
func GetPostWithUser(pid uint) (*model.Post, error) {
	var post model.Post
	err := config.DB.Preload("User").First(&post).Error
	return &post, err
}

// 用户发布帖子加生成日志
func CreatePostByTx(post *model.Post) error {
	// 开启事务——————要么全部成功，要么全部失败，不只干一半
	//确保安全
	tx := config.DB.Begin()
	// 执行帖子插入
	if err := tx.Create(post).Error; err != nil {
		tx.Rollback() // 失败回滚
		fmt.Println("创建帖子失败！")
		return err
	}
	// 提交事务
	fmt.Println("创建帖子成功！")
	return tx.Commit().Error
}
