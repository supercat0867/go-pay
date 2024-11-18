package repository

import (
	"go-pay/internal/model"
	"gorm.io/gorm"
)

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{
		db: db,
	}
}

func (r *OrderRepo) Create(pay *model.Order) error {
	return r.db.Create(&pay).Error
}

func (r *OrderRepo) FindByMchIDAndOrderID(mchId, orderId string) (*model.Order, error) {
	var order model.Order
	if err := r.db.Where("MchID = ?", mchId).Where("TradeNo = ?", orderId).Preload("Merchant").
		First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepo) Update(pay *model.Order) error {
	return r.db.Save(pay).Error
}

// FindPendingPay 查询支付中的订单
func (r *OrderRepo) FindPendingPay() []model.Order {
	var orders []model.Order
	r.db.Where("PayState = ?", model.PayStatePending).Preload("Merchant").Find(&orders)
	return orders
}

// FindPendingRefund 查询退款中的订单
func (r *OrderRepo) FindPendingRefund() []model.Order {
	var orders []model.Order
	r.db.Where("PayState = ?", model.PayStateRefund).Preload("Merchant").Find(&orders)
	return orders
}
