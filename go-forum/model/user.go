package model

import (
	"go-forum/config"
	"time"
)

type User struct {
	Id        uint      `gorm:"column:id;primary_key"`
	Username  string    `gorm:"column:username;unique"` //唯一
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func InitUserTable() {
	config.DB.AutoMigrate(&User{})
}
