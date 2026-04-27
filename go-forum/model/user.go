package model

import (
	"go-forum/config"
	"time"
)

type User struct {
	Id        uint      `gorm:"column:id;primary_key" json:"id"`
	Username  string    `gorm:"column:username;unique" json:"username"`
	Password  string    `gorm:"column:password" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func InitUserTable() {
	config.DB.AutoMigrate(&User{})
}
