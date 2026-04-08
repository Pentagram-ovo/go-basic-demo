package main

import (
	"context"
	"fmt"
	"grpc-user-demo/userpb"
	"net"

	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 定义服务结构体
type userServer struct {
	userpb.UnimplementedUserServiceServer
	// 模拟数据库
	db *gorm.DB
}

func newUserServer(db *gorm.DB) *userServer {
	return &userServer{db: db}
}

type User struct {
	ID       uint32 `gorm:"primaryKey"`
	Username string
	Password string
}

// 实现 CreateUser
func (s *userServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	pbUser := req.User
	user := User{
		ID:       pbUser.Id,
		Username: pbUser.Username,
		Password: pbUser.Password,
	}

	// MySQL 插入
	err := s.db.Create(&user).Error
	if err != nil {
		return &userpb.CreateUserResponse{
			Success: false,
			Message: "创建失败：" + err.Error(),
		}, nil
	}

	return &userpb.CreateUserResponse{
		Success: true,
		Message: "创建成功",
	}, nil
}

// 实现 GetUser
func (s *userServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	id := req.Id
	var user User

	// MySQL 查询
	err := s.db.First(&user, id).Error
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	pbUser := &userpb.User{
		Id:       user.ID,
		Username: user.Username,
		Password: user.Password,
	}

	return &userpb.GetUserResponse{User: pbUser}, nil
}

func main() {

	dsn := "root:123456@tcp(127.0.0.1:3306)/grpc-user?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("MySQL 连接失败")
	}

	// 自动建表
	db.AutoMigrate(&User{})
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	// 创建gRPC服务器
	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, newUserServer(db))

	fmt.Println("gRPC 服务运行在 :50051，已连接 MySQL")
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
