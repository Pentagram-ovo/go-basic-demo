package service

import (
	"errors"
	"go-im-demo/config"
	"go-im-demo/internal/middleware"
	"go-im-demo/internal/model"

	"gorm.io/gorm"
)

type UserService struct {
	Db *gorm.DB
}

func NewUserService() *UserService {
	return &UserService{
		Db: config.DB,
	}
}

// Register 用户注册，返回创建后的用户（不含密码字段）
func (s *UserService) Register(username, password string) (*model.User, error) {
	var exist model.User
	result := s.Db.Where("username = ?", username).First(&exist)
	if result.Error != gorm.ErrRecordNotFound {
		return nil, errors.New("用户名已被使用")
	}
	if len(username) < 2 || len(password) < 6 {
		return nil, errors.New("用户名或密码太短！")
	}
	hashPassword, err := middleware.HashPassword(password)
	if err != nil {
		return nil, err
	}
	var user model.User
	user.Username = username
	user.Password = hashPassword
	if err := s.Db.Create(&user).Error; err != nil {
		return nil, errors.New("注册用户失败")
	}
	return &user, nil
}

func (s *UserService) Login(username, password string) (*model.User, error) {
	var user model.User
	result := s.Db.Where("username = ?", username).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, errors.New("用户还未注册！")
	}
	if middleware.CheckPassword(password, user.Password) {
		return &user, nil
	}
	return nil, errors.New("用户名或密码错误！")
}

func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	result := s.Db.Where("username = ?", username).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}

func (s *UserService) GetUserById(id uint) (*model.User, error) {
	var user model.User
	result := s.Db.Where("id = ?", id).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}
