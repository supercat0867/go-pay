package database

import (
	"go-pay/internal/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.Merchant{}, &model.Order{}); err != nil {
		return err
	}
	return nil
}
