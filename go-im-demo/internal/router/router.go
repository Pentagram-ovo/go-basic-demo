package router

import (
	"encoding/json"
	"log"
	"net/http"

	"go-im-demo/internal/api"
	"go-im-demo/internal/middleware"
	"go-im-demo/internal/service"
	"go-im-demo/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	hub      = ws.NewHub()
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func SetupRouter(
	userHandler *api.UserHandler,
	messageHandler *api.MessageHandler,
	messageService *service.MessageService,
	userService *service.UserService,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())

	// 静态文件
	r.StaticFile("/", "./public/index.html")
	r.StaticFile("/login.html", "./public/login.html")
	r.Static("/static", "./public/static")

	// REST 接口
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.GET("/messages", messageHandler.GetHistory)

	// WebSocket
	r.GET("/ws", func(c *gin.Context) {
		handleWebSocket(c, messageService, userService)
	})

	return r
}

func handleWebSocket(c *gin.Context, msgService *service.MessageService, userService *service.UserService) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少token"})
		return
	}
	claims, err := middleware.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效token"})
		return
	}
	if claims.UserName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效用户信息"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v\n", err)
		return
	}

	client := ws.NewClient(claims.UserID, claims.UserName, conn, msgService, userService)
	hub.Register(claims.UserName, client)

	// 拉取离线消息
	unreadMsgs, err := msgService.GetUnreadMessages(claims.UserID)
	if err != nil {
		log.Printf("拉取离线消息失败: %v", err)
	} else {
		for _, msg := range unreadMsgs {
			fromUser, err := userService.GetUserById(msg.FromID)
			if err != nil {
				log.Printf("获取发送者信息失败: %v", err)
				continue // 修复：改为 continue，不要 return
			}
			pushMsg := map[string]interface{}{
				"from":    fromUser.Username,
				"to":      claims.UserName,
				"content": msg.Content,
				"time":    msg.CreatedAt.Unix(),
			}
			data, _ := json.Marshal(pushMsg)
			client.Send <- data

			// 标记已读
			if err := msgService.MarkAsRead(msg.ID); err != nil {
				log.Printf("标记已读失败: %v", err)
			}
		}
	}

	client.Send <- []byte("欢迎 " + claims.UserName + " 加入聊天室！")

	go client.WritePump()
	go client.ReadPump(hub)

	log.Printf("用户 %s 已上线\n", claims.UserName)
}
