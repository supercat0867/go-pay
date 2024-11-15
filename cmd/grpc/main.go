package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"go-pay/internal/handler"
	"go-pay/internal/repository"
	service2 "go-pay/internal/service"
	"go-pay/pkg/database"
	pb "go-pay/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	listen, _ := net.Listen("tcp", ":9090")
	grpcServer := grpc.NewServer()

	database.InitDB()
	merchantRepo := repository.NewMerchantRepo(database.DB)
	orderRepo := repository.NewOrderRepo(database.DB)
	service := service2.NewService(merchantRepo, orderRepo)

	// 注册服务
	pb.RegisterMerchantServiceServer(grpcServer, &handler.Handler{Service: service})
	pb.RegisterPayServiceServer(grpcServer, &handler.Handler{Service: service})

	// 启动服务
	err = grpcServer.Serve(listen)
	if err != nil {
		fmt.Println("启动服务失败", err)
	}
}
