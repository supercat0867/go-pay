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
	var pay model.Order
	if err := r.db.Where("MchID = ?", mchId).Where("TradeNo = ?", orderId).First(&pay).Error; err != nil {
		return nil, err
	}
	return &pay, nil
}

func (r *OrderRepo) Update(pay *model.Order) error {
	return r.db.Save(pay).Error
}
