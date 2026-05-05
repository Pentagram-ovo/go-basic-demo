package ws

import (
	"go-im-demo/internal/service"
	"log"
	"sync"
)

// Hub 管理所有在线连接
type Hub struct {
	clients sync.Map // key: userID(string), value: *Client
}

// NewHub 创建一个新的 Hub
func NewHub() *Hub {
	return &Hub{}
}

// Register 将用户连接注册到 Hub
func (h *Hub) Register(username string, client *Client) {
	h.clients.Store(username, client)
	err := service.AddUserOnline(username)
	if err != nil {
		log.Println(err)
	}
}

// Unregister 从 Hub 中移除用户连接
func (h *Hub) Unregister(username string) {
	h.clients.Delete(username)
	err := service.RemoveUserOnline(username)
	if err != nil {
		log.Println(err)
	}
}

// GetClient 获取指定用户的连接
func (h *Hub) GetClient(username string) (*Client, bool) {
	value, ok := h.clients.Load(username)
	if !ok {
		return nil, false
	}
	return value.(*Client), true
}

// SendToUser 向指定用户发送消息，如果用户不在线则返回 false
func (h *Hub) SendToUser(username string, message []byte) bool {
	client, ok := h.GetClient(username)
	if !ok {
		return false
	}
	client.Send <- message
	return true
}

// OnlineUsers 返回当前在线用户 ID 列表（调试用）
func (h *Hub) OnlineUsers() []string {
	var users []string
	users, err := service.GetOnlineUsers()
	if err != nil {
		log.Println(err)
		return nil
	}
	return users
}
