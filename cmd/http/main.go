package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-pay/internal/handler"
	"go-pay/internal/repository"
	service2 "go-pay/internal/service"
	"go-pay/pkg/database"
	"log"
)

func main() {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	r := gin.Default()

	database.InitDB()
	merchantRepo := repository.NewMerchantRepo(database.DB)
	payRepo := repository.NewPayRepo(database.DB)
	service := service2.NewService(merchantRepo, payRepo)
	controller := &handler.Handler{Service: service}

	// 微信支付通知回调
	r.POST("/wechatpay/notify/:id", controller.DealWechatPayNotify)

	r.Run(":8080")
}
