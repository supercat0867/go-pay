package repository

import (
	"fmt"
	"go-pay/internal/model"
	"gorm.io/gorm"
)

type MerchantRepo struct {
	db *gorm.DB
}

func NewMerchantRepo(db *gorm.DB) *MerchantRepo {
	return &MerchantRepo{
		db: db,
	}
}

func (r *MerchantRepo) Create(m *model.Merchant) error {
	return r.db.Create(&m).Error
}

func (r *MerchantRepo) Find(page, pageSize int, params map[string]interface{}) ([]model.Merchant, int) {
	var merchants []model.Merchant
	var total int64
	sql := r.db.Model(&model.Merchant{})
	for k, v := range params {
		if v != "" {
			if k == "Name" {
				sql.Where("Name LIKE ?", fmt.Sprintf("%%%s%%", v))
			} else if k == "PlantForm" {
				sql.Where("PlantForm = ?", v)
			} else if k == "MchID" {
				sql.Where("MchID LIKE ?", fmt.Sprintf("%%%s%%", v))
			}
		}
	}
	sql.Count(&total).Offset((page - 1) * pageSize).Limit(pageSize).Find(&merchants)
	return merchants, int(total)
}

// FindByMchID 根据商户号查询商户信息
func (r *MerchantRepo) FindByMchID(mchId string) (*model.Merchant, error) {
	var merchant model.Merchant
	if err := r.db.Where("MchID = ?", mchId).First(&merchant).Error; err != nil {
		return nil, err
	}
	return &merchant, nil
}
