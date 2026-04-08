package main

import (
	"context"
	"fmt"
	"grpc-user-demo/userpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		":50051", //自定义地址
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 创建 gRPC 客户端
	// New Xxx ServiceClient
	client := userpb.NewUserServiceClient(conn)

	// 1. 创建用户
	createResp, err := client.CreateUser(context.Background(), &userpb.CreateUserRequest{
		User: &userpb.User{
			Id:       3,
			Username: "zhangqingli",
			Password: "0331",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("创建用户结果：", createResp)

	// 2. 获取用户
	getResp, err := client.GetUser(context.Background(), &userpb.GetUserRequest{Id: 1})
	if err != nil {
		panic(err)
	}
	fmt.Println("查询到用户：", getResp.User)
}
