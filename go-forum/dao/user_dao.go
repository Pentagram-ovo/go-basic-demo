package dao

import (
	"go-forum/config"
	"go-forum/model"
)

// CreateUser 创建用户
func CreateUser(user *model.User) error {
	return config.DB.Create(user).Error
}

// GetUserById 根据用户id来查询用户
func GetUserById(id uint) (*model.User, error) {
	var user model.User
	err := config.DB.Where("id=?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 根据用户名称来查询用户
func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := config.DB.Where("username=?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
