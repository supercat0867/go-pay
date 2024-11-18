package model

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model
	MerchantID    uint      `gorm:"index"`    // 商户id
	MchID         string    `gorm:"not null"` // 商户号
	TradeNo       string    `gorm:"not null"` // 系统订单号
	RefundNo      string    // 退款单号
	TransactionID string    // 平台订单号
	RefundID      string    // 退款单号
	PayType       PayType   `gorm:"not null;size:1"` // 交易类型
	OpenID        string    // 微信openid
	Amount        float32   // 金额
	ExpireAt      time.Time `gorm:"not null"`        // 过期时间
	PayState      PayState  `gorm:"not null;size:1"` // 支付状态
	NotifyUrl     string    // 回调地址

	Merchant Merchant `gorm:"foreignKey:MerchantID" `
}

type PayState int

const (
	PayStatePending       PayState = iota + 1 // 待支付
	PayStateSuccess                           // 支付成功
	PayStateFail                              // 支付失败
	PayStateRefund                            // 退款中
	PayStateRefundSuccess                     // 退款成功
	PayStateRefundFail                        // 退款失败
	PayStateClosed                            // 已过期/已关闭
)

type PayType int

const (
	PayTypeWechatJSAPI    PayType = iota + 1 // 微信公众号支付
	PayTypeWechatNative                      // 微信扫码支付
	PayTypeWechatApp                         // 微信APP支付
	PayTypeWechatMicroPay                    // 微信付款码支付
	PayTypeWechatMWEB                        // 微信H5支付
	PayTypeWechatFACEPay                     // 微信刷脸支付
)
