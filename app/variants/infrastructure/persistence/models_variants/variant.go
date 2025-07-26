package models_variants

import "github.com/shopspring/decimal"

type VariantModel struct {
	ID        uint             `gorm:"primaryKey"`
	ProductID uint             `gorm:"not null"`
	Name      string           `gorm:"not null"`
	SKU       string           `gorm:"uniqueIndex;not null"`
	Price     *decimal.Decimal `gorm:"type:decimal(10,2);null"`
}

func (VariantModel) TableName() string {
	return "product_variants"
}
