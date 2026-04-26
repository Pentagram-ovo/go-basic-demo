package service

import (
	"errors"
	"go-forum/dao"
	"go-forum/model"
	"go-forum/utils"
)

// SetUser 注册新用户（用户名不重复）
func SetUser(username, password string) error {
	var user model.User
	exist, _ := dao.GetUserByUsername(username)
	if exist != nil {
		return errors.New("用户已存在")
	}
	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	user.Username = username
	user.Password = hash
	return dao.CreateUser(&user)
}

// Login 用户登录，返回token
func Login(username, password string) (string, error) {
	user, err := dao.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("用户不存在！")
	}
	//检查密码是否正确
	if utils.CheckPassword(password, user.Password) {
		return utils.GenerateToken(user.Id)
	}
	return "", errors.New("密码错误！")
}
