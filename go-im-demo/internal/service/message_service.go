package service

import (
	"go-im-demo/internal/model"
	"log"

	"gorm.io/gorm"
)

type MessageService struct {
	Db       *gorm.DB
	saveChan chan model.Message // 缓冲通道，异步保存用
}

func NewMessageService(db *gorm.DB) *MessageService {
	ms := &MessageService{
		Db:       db,
		saveChan: make(chan model.Message, 100), // 缓冲大小100
	}
	go ms.saveWorker() // 启动后台落库协程
	return ms
}

// SaveMessage 供外部调用，把消息丢进通道，非阻塞
func (s *MessageService) SaveMessage(msg model.Message) {
	// 将 msg 放进 s.saveChan，如果通道满可以 log 一下但不阻塞
	select {
	case s.saveChan <- msg:
	default:
		log.Println("[消息落库] 通道满，丢弃一条消息")
	}
}

// saveWorker 在后台 goroutine 中运行，从通道读取消息并写入数据库
func (s *MessageService) saveWorker() {
	for msg := range s.saveChan {
		if err := s.Db.Create(&msg).Error; err != nil {
			log.Println(err)
		}
	}
}

// GetHistory 查询两个用户之间的历史消息（双向），分页按时间倒序
func (s *MessageService) GetHistory(userID1, userID2 uint, page, size int) ([]model.Message, int64, error) {
	query := s.Db.Model(&model.Message{}).
		Where("(from_id = ? AND to_id = ?) OR (from_id = ? AND to_id = ?)",
			userID1, userID2, userID2, userID1)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var messages []model.Message
	offset := (page - 1) * size
	err := query.Order("created_at desc").
		Offset(offset).
		Limit(size).
		Find(&messages).Error
	return messages, total, err
}

// GetUnreadMessages 获取未读消息（发给指定用户且未读）
func (s *MessageService) GetUnreadMessages(userID uint) ([]model.Message, error) {
	var msgs []model.Message
	err := s.Db.Where("to_id = ? AND is_read = ?", userID, false).
		Order("created_at asc").
		Find(&msgs).Error
	return msgs, err
}

// MarkAsRead 标记单条消息为已读
func (s *MessageService) MarkAsRead(msgID uint) error {
	return s.Db.Model(&model.Message{}).Where("id = ?", msgID).
		Update("is_read", true).Error
}
