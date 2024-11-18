package handler

import (
	"context"
	"github.com/joho/godotenv"
	"go-pay/internal/repository"
	service2 "go-pay/internal/service"
	"go-pay/pkg/database"
	pb "go-pay/proto"
	"log"
	"testing"
)

// 测试创建商户
func TestHandler_CreateMerchant(t *testing.T) {
	// 加载环境变量
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	database.InitDB()
	merchantRepo := repository.NewMerchantRepo(database.DB)
	orderRepo := repository.NewOrderRepo(database.DB)
	service := service2.NewService(merchantRepo, orderRepo)
	handler := Handler{
		Service: service,
	}

	resp, err := handler.CreateMerchant(context.Background(), &pb.CreateMerchantRequest{
		Name:     "测试商户",
		Platform: 1, // 微信支付
		AppId:    "123456",
		MchId:    "123456",
		Secret:   "123456",
		Cert:     "-----BEGIN PRIVATE KEY-----123456-----END PRIVATE KEY-----",
		CertNum:  "123456",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp)
}

// 测试获取商户列表
func TestHandler_GetMerchants(t *testing.T) {
	// 加载环境变量
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	database.InitDB()
	merchantRepo := repository.NewMerchantRepo(database.DB)
	orderRepo := repository.NewOrderRepo(database.DB)
	service := service2.NewService(merchantRepo, orderRepo)
	handler := Handler{
		Service: service,
	}

	resp, err := handler.GetMerchants(context.Background(), &pb.GetMerchantsRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp)
}

// 测试获取微信预支付信息
func TestHandler_GetWechatPrepayInfoJsAPI(t *testing.T) {
	// 加载环境变量
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	database.InitDB()
	merchantRepo := repository.NewMerchantRepo(database.DB)
	orderRepo := repository.NewOrderRepo(database.DB)
	service := service2.NewService(merchantRepo, orderRepo)
	handler := Handler{
		Service: service,
	}

	resp, err := handler.GetWechatPrepayInfoJsAPI(context.Background(), &pb.WechatPrepayInfoJsAPIRequest{
		MchId:       "123456",
		Amount:      1,
		Description: "测试商品",
		ExpireTime:  "1731674980000",
		Openid:      "ouZ****************SLa",
		OutTradeNo:  "S2024111111111111",
		NotifyUrl:   "http://127.0.0.1:8080",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp)
}

// 测试微信支付退款
func TestHandler_WechatPayRefund(t *testing.T) {
	// 加载环境变量
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	database.InitDB()
	merchantRepo := repository.NewMerchantRepo(database.DB)
	orderRepo := repository.NewOrderRepo(database.DB)
	service := service2.NewService(merchantRepo, orderRepo)
	handler := Handler{
		Service: service,
	}

	resp, err := handler.WechatPayRefund(context.Background(), &pb.WechatPayRefundRequest{
		MchId:       "123456",
		Refund:      1,
		Total:       1,
		OutTradeNo:  "S2024111111111111",
		OutRefundNo: "TK2024111111111111",
		Reason:      "用户主动退款",
		NotifyUrl:   "http://127.0.0.1:8080",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp)
}
