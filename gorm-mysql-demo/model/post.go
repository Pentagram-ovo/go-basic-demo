package model

import "time"

type Post struct {
	Id        uint      `gorm:"column:id;primary_key"`
	Title     string    `gorm:"column:title"`
	Content   string    `gorm:"column:content"`
	UserId    uint      `gorm:"column:user_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
