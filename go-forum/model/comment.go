package model

import (
	"go-forum/config"
	"time"
)

type Comment struct {
	Id        uint      `gorm:"column:id;primary_key" json:"id"`
	UserID    uint      `gorm:"column:user_id" json:"user_id"`
	PostID    uint      `gorm:"column:post_id" json:"post_id"`
	Content   string    `gorm:"column:content" json:"content"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func InitCommentTable() {
	config.DB.AutoMigrate(&Comment{})
}
