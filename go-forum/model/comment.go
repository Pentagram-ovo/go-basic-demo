package model

import (
	"go-forum/config"
	"time"
)

type Comment struct {
	Id        uint      `gorm:"column:id;primary_key"`
	UserID    uint      `gorm:"column:user_id"`
	PostID    uint      `gorm:"column:post_id"`
	Content   string    `gorm:"column:content"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func InitCommentTable() {
	config.DB.AutoMigrate(&Comment{})
}
