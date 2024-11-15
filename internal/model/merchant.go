package model

import "gorm.io/gorm"

// Merchant 商户模型
type Merchant struct {
	gorm.Model
	Name      string    `gorm:"size:100"` // 商户名称
	PlantForm PlantForm `gorm:"size:1"`   // 支付平台
	AppID     string    `gorm:"size:100"` // 应用ID
	MchID     string    `gorm:"size:100"` // 商户号
	Cert      string    // 证书
	CertNum   string    // 证书序列号
	Secret    string    // 密钥
}

type PlantForm int

const (
	WeChatPay PlantForm = iota + 1 // 微信支付
	AliPay                         // 支付宝
)

func (f PlantForm) String() string {
	switch f {
	case WeChatPay:
		return "微信支付"
	case AliPay:
		return "支付宝"
	default:
		return "未知平台"
	}
}
