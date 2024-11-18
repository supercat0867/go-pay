package service

import (
	"errors"
	"fmt"
	"go-pay/internal/model"
	"go-pay/internal/repository"
	"go-pay/pkg/utils"
	"go-pay/pkg/wechatpay/jsapi"
	pb "go-pay/proto"
	"log"
	"strconv"
)

type Service struct {
	MerchantRepo *repository.MerchantRepo
	OrderRepo    *repository.OrderRepo
}

func NewService(merchantRepo *repository.MerchantRepo, orderRepo *repository.OrderRepo) *Service {
	return &Service{
		MerchantRepo: merchantRepo,
		OrderRepo:    orderRepo,
	}
}

// CreateMerchant 新增商户
func (s *Service) CreateMerchant(req *pb.CreateMerchantRequest) (*pb.CreateMerchantResponse, error) {
	var merchant *model.Merchant
	switch model.PlantForm(req.Platform) {
	case model.WeChatPay:
		// 微信支付平台
		merchant = &model.Merchant{
			Name:      req.Name,
			PlantForm: model.WeChatPay,
			AppID:     req.AppId,
			MchID:     req.MchId,
			Cert:      req.Cert,
			CertNum:   req.CertNum,
			Secret:    req.Secret,
		}
	default:
		return nil, fmt.Errorf("invalid platform")
	}

	if err := s.MerchantRepo.Create(merchant); err != nil {
		return nil, err
	}

	return &pb.CreateMerchantResponse{
		Merchant: &pb.Merchant{
			Id:            uint64(merchant.ID),
			Name:          merchant.Name,
			PlantForm:     uint64(merchant.PlantForm),
			PlantFormName: merchant.PlantForm.String(),
			AppId:         merchant.AppID,
			MchId:         merchant.MchID,
			Cert:          "************",
			CertNum:       "************",
			Secret:        "************",
			CreatedAt:     strconv.FormatInt(merchant.CreatedAt.Unix(), 10),
		},
	}, nil
}

// GetMerchants 查询商户列表
func (s *Service) GetMerchants(req *pb.GetMerchantsRequest) *pb.GetMerchantsResponse {
	items, total := s.MerchantRepo.Find(int(req.Page), int(req.PageSize), map[string]interface{}{
		"Name":      req.Name,
		"PlantForm": req.PlantForm,
		"MchID":     req.MchId,
	})
	var merchants []*pb.Merchant
	for _, item := range items {
		merchants = append(merchants, &pb.Merchant{
			Id:            uint64(item.ID),
			Name:          item.Name,
			PlantForm:     uint64(item.PlantForm),
			PlantFormName: item.PlantForm.String(),
			AppId:         item.AppID,
			MchId:         item.MchID,
			Cert:          "************",
			CertNum:       "************",
			Secret:        "************",
			CreatedAt:     strconv.FormatInt(item.CreatedAt.Unix(), 10),
		})
	}
	return &pb.GetMerchantsResponse{
		Total:     uint64(total),
		Merchants: merchants,
		Page:      req.Page,
		PageSize:  req.PageSize,
	}
}

// GetWechatPrePayInfoJsAPI 发起微信支付-jsapi
func (s *Service) GetWechatPrePayInfoJsAPI(req *pb.WechatPrepayInfoJsAPIRequest) (*pb.WechatPrepayInfoJsAPIResponse, error) {
	// 查询商户
	merchant, err := s.MerchantRepo.FindByMchID(req.MchId)
	if err != nil {
		return nil, errors.New("merchant not found")
	}
	// 检查商户平台是否匹配
	if merchant.PlantForm != model.WeChatPay {
		return nil, errors.New("platform not match")
	}

	// 检查订单是否已存在
	_, err = s.OrderRepo.FindByMchIDAndOrderID(merchant.MchID, req.OutTradeNo)
	if err == nil {
		return nil, errors.New("order already exists")
	}

	// 实例化jsapi client
	client := jsapi.NewClient(merchant.AppID, merchant.MchID, merchant.Secret, merchant.Cert, merchant.CertNum)
	// 获取预支付信息
	expireTime, err := utils.ConvertMillisecondsToTime(req.ExpireTime)
	if err != nil {
		return nil, errors.New("expire time format error")
	}

	prepayInfo, err := client.GetPrepayInfo(req.Description, req.OutTradeNo, req.Openid, req.NotifyUrl,
		expireTime, req.Amount)
	if err != nil {
		return nil, err
	}

	// 创建支付记录
	pay := &model.Order{
		MchID:     merchant.MchID,
		TradeNo:   req.OutTradeNo,
		PayState:  model.PayStatePending,
		PayType:   model.PayTypeWechatJSAPI,
		Amount:    req.Amount,
		OpenID:    req.Openid,
		ExpireAt:  expireTime,
		NotifyUrl: req.NotifyUrl,
	}
	if err = s.OrderRepo.Create(pay); err != nil {
		log.Println(err)
	}

	return &pb.WechatPrepayInfoJsAPIResponse{
		AppId:     prepayInfo.AppId,
		TimeStamp: prepayInfo.Timestamp,
		NonceStr:  prepayInfo.NonceStr,
		Package:   prepayInfo.Package,
		PaySign:   prepayInfo.PaySign,
		SignType:  prepayInfo.SignType,
	}, nil
}

// WechatPayRefund 微信支付退款
func (s *Service) WechatPayRefund(req *pb.WechatPayRefundRequest) (*pb.WechatPayRefundResponse, error) {
	// 查询商户
	merchant, err := s.MerchantRepo.FindByMchID(req.MchId)
	if err != nil {
		return nil, errors.New("merchant not found")
	}
	// 检查商户平台是否匹配
	if merchant.PlantForm != model.WeChatPay {
		return nil, errors.New("platform not match")
	}

	// 查询订单
	order, err := s.OrderRepo.FindByMchIDAndOrderID(merchant.MchID, req.OutTradeNo)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// 实例化jsapi client
	client := jsapi.NewClient(merchant.AppID, merchant.MchID, merchant.Secret, merchant.Cert, merchant.CertNum)
	// 发起退款
	resp, err := client.Refund(order.TradeNo, req.OutRefundNo, req.Reason, req.NotifyUrl, req.Refund, req.Total)
	if err != nil {
		return nil, err
	}

	order.RefundNo = req.OutRefundNo
	order.RefundID = *resp.RefundId
	order.PayState = model.PayStateRefund
	order.NotifyUrl = req.NotifyUrl
	if err = s.OrderRepo.Update(order); err != nil {
		log.Println(err)
	}

	return &pb.WechatPayRefundResponse{
		Channel:       string(*resp.Channel),
		OutRefundNo:   *resp.OutRefundNo,
		OutTradeNo:    *resp.OutTradeNo,
		RefundId:      *resp.RefundId,
		TransactionId: *resp.TransactionId,
	}, nil
}
