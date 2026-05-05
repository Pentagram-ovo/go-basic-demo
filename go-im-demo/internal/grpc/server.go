package grpc

import (
	"context"
	"go-im-demo/internal/model"
	"go-im-demo/internal/service"
	pb "go-im-demo/proto/chat"
)

type ChatServer struct {
	pb.UnimplementedChatServiceServer
	UserService    *service.UserService
	MessageService *service.MessageService
}

func NewChatServer(userSvc *service.UserService, msgSvc *service.MessageService) *ChatServer {
	return &ChatServer{UserService: userSvc, MessageService: msgSvc}
}

// SendMessage 实现 RPC：查询目标用户ID，保存消息（已读标记由调用方决定）
func (s *ChatServer) SendMessage(ctx context.Context, req *pb.SendMessageReq) (*pb.SendMessageResp, error) {
	msg := model.Message{
		FromID:  uint(req.FromId),
		ToID:    uint(req.ToId),
		Content: req.Content,
		IsRead:  req.IsRead,
	}
	s.MessageService.SaveMessage(msg)
	return &pb.SendMessageResp{Success: true}, nil
}

// GetHistory 实现 RPC，直接调用 MessageService.GetHistory
func (s *ChatServer) GetHistory(ctx context.Context, req *pb.GetHistoryReq) (*pb.GetHistoryResp, error) {
	msgs, total, err := s.MessageService.GetHistory(uint(req.UserId), uint(req.PeerId), int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}
	var pbMsgs []*pb.Message
	for _, m := range msgs {
		pbMsgs = append(pbMsgs, &pb.Message{
			Id:        uint64(m.ID),
			FromId:    uint64(m.FromID),
			ToId:      uint64(m.ToID),
			Content:   m.Content,
			CreatedAt: m.CreatedAt.Unix(),
		})
	}
	return &pb.GetHistoryResp{Messages: pbMsgs, Total: total}, nil
}
