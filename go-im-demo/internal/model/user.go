package model

import (
	"go-im-demo/config"
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"` // json:"-" 防止密码字段序列化时泄露
	CreatedAt time.Time `json:"created_at"`
}

func InitUserTable() {
	config.DB.AutoMigrate(&User{})
}
