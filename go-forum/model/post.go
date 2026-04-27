package model

import (
	"go-forum/config"
	"time"
)

type Post struct {
	Id        uint      `gorm:"column:id;primary_key" json:"id"`
	Title     string    `gorm:"column:title;not null" json:"title"`
	Content   string    `gorm:"column:content" json:"content"`
	UserID    uint      `gorm:"column:user_id" json:"user_id"`
	LikeCount int       `gorm:"column:like_count;default:0" json:"like_count"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func InitPostTable() {
	config.DB.AutoMigrate(&Post{})
}
