package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"go-pay/internal/handler"
	"go-pay/internal/model"
	"go-pay/internal/repository"
	service2 "go-pay/internal/service"
	"go-pay/pkg/database"
	"go-pay/pkg/wechatpay/jsapi"
	pb "go-pay/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
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

	// 定时检查订单状态
	go func() {
		for {
			orders := orderRepo.FindPendingPay()
			for _, order := range orders {
				client := jsapi.NewClient(order.Merchant.AppID, order.Merchant.MchID, order.Merchant.Secret,
					order.Merchant.Cert, order.Merchant.CertNum)
				resp, err := client.QueryByOutTradeNo(order.TradeNo)
				if err != nil {
					log.Printf("订单商家单号：%s状态查询失败：%v", order.TradeNo, err)
				}
				if *resp.TradeState == "SUCCESS" {
					order.PayState = model.PayStateSuccess
					order.TransactionID = *resp.TransactionId
					if err = orderRepo.Update(&order); err != nil {
						log.Printf("订单商家单号：%s状态更新失败：%v", order.TradeNo, err)
					}
				}
			}
			time.Sleep(time.Second * 1)
		}
	}()

	// 定时检查订单状态
	go func() {
		for {
			orders := orderRepo.FindPendingRefund()
			for _, order := range orders {
				client := jsapi.NewClient(order.Merchant.AppID, order.Merchant.MchID, order.Merchant.Secret,
					order.Merchant.Cert, order.Merchant.CertNum)
				resp, err := client.QueryRefundByOutRefundNo(order.RefundNo)
				if err != nil {
					log.Printf("订单商家退款单号：%s状态查询失败：%v", order.RefundNo, err)
				}
				if *resp.Status == "SUCCESS" {
					order.PayState = model.PayStateRefundSuccess
					if err = orderRepo.Update(&order); err != nil {
						log.Printf("订单商家退款单号：%s状态更新失败：%v", order.RefundNo, err)
					}
				}
			}
			time.Sleep(time.Second * 1)
		}
	}()

	// 启动服务
	err = grpcServer.Serve(listen)
	if err != nil {
		fmt.Println("启动服务失败", err)
	}
}
