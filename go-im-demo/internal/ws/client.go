package ws

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"go-im-demo/internal/grpc"
	"go-im-demo/internal/service"
	pb "go-im-demo/proto/chat"

	"github.com/gorilla/websocket"
)

// Client 代表一个 WebSocket 连接
type Client struct {
	UserID      uint
	UserName    string
	Conn        *websocket.Conn
	Send        chan []byte
	MsgService  *service.MessageService
	UserService *service.UserService
}

type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
	Time    int64  `json:"time"`
}

// NewClient 创建一个新的客户端实例
func NewClient(userID uint, userName string, conn *websocket.Conn, msgService *service.MessageService, userService *service.UserService) *Client {
	return &Client{
		UserID:      userID,
		UserName:    userName,
		Conn:        conn,
		Send:        make(chan []byte, 256),
		MsgService:  msgService,
		UserService: userService,
	}
}

// ReadPump 从 WebSocket 读取消息
func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.Unregister(c.UserName)
		close(c.Send)
		c.Conn.Close()
		log.Printf("用户 %s 已断开连接\n", c.UserName)
	}()

	for {
		_, rawMsg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		// 心跳处理
		var typeMsg struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(rawMsg, &typeMsg); err == nil && typeMsg.Type == "ping" {
			c.Send <- []byte(`{"type":"pong"}`)
			service.RefreshUserOnline(c.UserName)
			continue
		}

		var msg Message
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			log.Println("JSON解析失败:", err)
			continue
		}
		msg.From = c.UserName
		msg.Time = time.Now().Unix()

		// /online 命令
		if msg.Content == "/online" {
			onlineUsers := hub.OnlineUsers()
			jsonUsers, _ := json.Marshal(onlineUsers)
			hub.SendToUser(c.UserName, jsonUsers)
			continue
		}

		// 消息校验
		if msg.To == "" {
			c.Send <- []byte("消息缺少接收人")
			continue
		}
		if msg.To == c.UserName {
			c.Send <- []byte("不能给自己发消息")
			continue
		}

		// 序列化消息用于实时转发
		jsonData, err := json.Marshal(msg)
		if err != nil {
			c.Send <- []byte("系统错误：消息发送失败")
			continue
		}

		// 尝试在线转发
		if success := hub.SendToUser(msg.To, jsonData); success {
			// 在线情况：异步通过 gRPC 落库，标记已读
			if targetClient, ok := hub.GetClient(msg.To); ok {
				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
					defer cancel()
					_, err := grpc.ChatClient.SendMessage(ctx, &pb.SendMessageReq{
						FromId:  uint64(c.UserID),
						ToId:    uint64(targetClient.UserID),
						Content: msg.Content,
						IsRead:  true,
					})
					if err != nil {
						log.Println("gRPC落库失败(在线):", err)
					}
				}()
			}
		} else {
			// 离线情况：查询目标用户ID，然后异步通过 gRPC 落库，标记未读
			targetUser, err := c.UserService.GetUserByUsername(msg.To)
			if err != nil {
				c.Send <- []byte("系统错误：目标用户不存在")
				continue
			}
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				_, err := grpc.ChatClient.SendMessage(ctx, &pb.SendMessageReq{
					FromId:  uint64(c.UserID),
					ToId:    uint64(targetUser.ID),
					Content: msg.Content,
					IsRead:  false,
				})
				if err != nil {
					log.Println("gRPC落库失败(离线):", err)
				}
			}()
			c.Send <- []byte("系统提示：对方不在线，消息已离线保存")
		}
	}
}

// WritePump 从 Send 通道取出消息并写入 WebSocket
func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()

	for message := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("向用户 %s 发送消息失败: %v\n", c.UserName, err)
			break
		}
	}
}
