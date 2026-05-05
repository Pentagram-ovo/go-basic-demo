package model

import (
	"go-im-demo/config"
	"time"
)

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FromID    uint      `gorm:"index;not null" json:"from_id"`
	ToID      uint      `gorm:"index;not null" json:"to_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

func InitMessageTable() {
	config.DB.AutoMigrate(&Message{})
}
