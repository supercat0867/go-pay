package handler

import (
	"context"
	"go-pay/internal/service"
	pb "go-pay/proto"
)

type Handler struct {
	Service *service.Service
	pb.UnimplementedMerchantServiceServer
	pb.UnimplementedPayServiceServer
}

// CreateMerchant 新增商户
func (h *Handler) CreateMerchant(ctx context.Context, req *pb.CreateMerchantRequest) (*pb.CreateMerchantResponse, error) {
	resp, err := h.Service.CreateMerchant(req)
	return resp, err
}

// GetMerchants 获取商户列表
func (h *Handler) GetMerchants(ctx context.Context, req *pb.GetMerchantsRequest) (*pb.GetMerchantsResponse, error) {
	return h.Service.GetMerchants(req), nil
}

// GetWechatPrepayInfoJsAPI 微信jsapi支付
func (h *Handler) GetWechatPrepayInfoJsAPI(ctx context.Context, req *pb.WechatPrepayInfoJsAPIRequest) (*pb.WechatPrepayInfoJsAPIResponse, error) {
	return h.Service.GetWechatPrePayInfoJsAPI(req)
}

// WechatPayRePayJsAPI 微信jsapi重新支付
func (h *Handler) WechatPayRePayJsAPI(ctx context.Context, req *pb.WechatRePayRequest) (*pb.WechatPrepayInfoJsAPIResponse, error) {
	return h.Service.WechatRePayInfoJsAPI(req)
}

// WechatPayRefund 微信支付退款
func (h *Handler) WechatPayRefund(ctx context.Context, req *pb.WechatPayRefundRequest) (*pb.WechatPayRefundResponse, error) {
	return h.Service.WechatPayRefund(req)
}
