package models

import (
	"github.com/mytheresa/go-hiring-challenge/app/variants/infrastructure/persistence/models_variants"
	"github.com/shopspring/decimal"
)

type ProductModel struct {
	ID         uint                           `gorm:"primaryKey"`
	Code       string                         `gorm:"uniqueIndex;not null"`
	Price      *decimal.Decimal               `gorm:"type:decimal(10,2);not null"`
	Variants   []models_variants.VariantModel `gorm:"foreignKey:ProductID"`
	CategoryID uint                           `gorm:"not null"`
	Category   CategoryModel                  `gorm:"foreignKey:CategoryID"`
}

func (ProductModel) TableName() string {
	return "products"
}
