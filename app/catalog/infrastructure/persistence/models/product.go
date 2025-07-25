package models

import "github.com/shopspring/decimal"

type ProductModel struct {
	ID         uint             `gorm:"primaryKey"`
	Code       string           `gorm:"uniqueIndex;not null"`
	Price      *decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	Variants   []VariantModel   `gorm:"foreignKey:ProductID"`
	CategoryID uint             `gorm:"not null"`
	Category   CategoryModel    `gorm:"foreignKey:CategoryID"`
}

func (ProductModel) TableName() string {
	return "products"
}
