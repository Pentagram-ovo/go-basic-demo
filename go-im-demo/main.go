package main

import (
	"go-im-demo/config"
	"go-im-demo/internal/api"
	"go-im-demo/internal/grpc"
	"go-im-demo/internal/model"
	"go-im-demo/internal/router"
	"go-im-demo/internal/service"
	pb "go-im-demo/proto/chat"
	"log"
	"net"

	grpclib "google.golang.org/grpc"
)

func main() {
	// 1. 初始化数据库和 Redis
	config.WaitForMySQL()
	config.InitMysql()
	config.InitRedis()
	model.InitMessageTable()
	model.InitUserTable()
	// 2. 初始化服务
	userService := service.NewUserService()
	messageService := service.NewMessageService(config.DB)

	// 3. 初始化 handler
	userHandler := api.NewUserHandler(userService)
	messageHandler := api.NewMessageHandler(messageService, userService)
	if err := grpc.InitChatClient("localhost:50051"); err != nil {
		log.Fatal("gRPC客户端连接失败:", err)
	}
	// gRPC 服务
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	s := grpclib.NewServer()
	pb.RegisterChatServiceServer(s, grpc.NewChatServer(userService, messageService))
	go func() {
		log.Println("gRPC 服务启动在 :50051")
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
	// 4. 创建路由引擎并启动
	r := router.SetupRouter(userHandler, messageHandler, messageService, userService)

	log.Println("服务启动在 :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("启动失败:", err)
	}
}
