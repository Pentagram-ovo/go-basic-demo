package grpc

import (
	pb "go-im-demo/proto/chat"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ChatClient pb.ChatServiceClient

func InitChatClient(addr string) error {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	ChatClient = pb.NewChatServiceClient(conn)
	return nil
}
