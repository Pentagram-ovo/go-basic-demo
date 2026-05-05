package api

import (
	"context"
	"net/http"
	"time"

	"go-im-demo/internal/grpc"
	"go-im-demo/internal/middleware"
	"go-im-demo/internal/model"
	"go-im-demo/internal/service"
	pb "go-im-demo/proto/chat"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	MsgService  *service.MessageService
	UserService *service.UserService
}

func NewMessageHandler(msgService *service.MessageService, userService *service.UserService) *MessageHandler {
	return &MessageHandler{
		MsgService:  msgService,
		UserService: userService,
	}
}

type HistoryRequest struct {
	PeerID   uint   `form:"peer_id"`
	PeerName string `form:"peer_name"`
	Page     int    `form:"page" binding:"min=1"`
	Size     int    `form:"size" binding:"min=1,max=50"`
}

// GetHistory 拉取当前用户与 peer_id 的聊天记录（通过 gRPC）
func (h *MessageHandler) GetHistory(c *gin.Context) {
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

	var req HistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var peerID uint
	if req.PeerName != "" {
		user, err := h.UserService.GetUserByUsername(req.PeerName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "对方用户不存在"})
			return
		}
		peerID = user.ID
	} else if req.PeerID != 0 {
		peerID = req.PeerID
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少peer_id或peer_name"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := grpc.ChatClient.GetHistory(ctx, &pb.GetHistoryReq{
		UserId: uint64(claims.UserID),
		PeerId: uint64(peerID),
		Page:   int32(req.Page),
		Size:   int32(req.Size),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询消息失败"})
		return
	}

	// 转换 gRPC 消息为 model.Message
	var messages []model.Message
	for _, m := range resp.Messages {
		messages = append(messages, model.Message{
			ID:        uint(m.Id),
			FromID:    uint(m.FromId),
			ToID:      uint(m.ToId),
			Content:   m.Content,
			IsRead:    false, // 历史消息 is_read 字段前端未使用，无所谓
			CreatedAt: time.Unix(m.CreatedAt, 0),
		})
	}
	if messages == nil {
		messages = []model.Message{}
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"total":    resp.Total,
		"page":     req.Page,
		"size":     req.Size,
		"peer_id":  peerID,
	})
}
