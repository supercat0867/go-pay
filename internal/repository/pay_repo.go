package repository

import (
	"go-pay/internal/model"
	"gorm.io/gorm"
)

type PayRepo struct {
	db *gorm.DB
}

func NewPayRepo(db *gorm.DB) *PayRepo {
	return &PayRepo{
		db: db,
	}
}

func (r *PayRepo) Create(pay *model.Pay) error {
	return r.db.Create(&pay).Error
}

func (r *PayRepo) FindByMchIDAndOrderID(mchId, orderId string) (*model.Pay, error) {
	var pay model.Pay
	if err := r.db.Where("MchID = ?", mchId).Where("TradeNo = ?", orderId).First(&pay).Error; err != nil {
		return nil, err
	}
	return &pay, nil
}

func (r *PayRepo) Update(pay *model.Pay) error {
	return r.db.Save(pay).Error
}
