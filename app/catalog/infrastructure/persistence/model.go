package persistence

import (
	"github.com/shopspring/decimal"
)

type CategoryModel struct {
	ID   uint   `gorm:"primaryKey"`
	Code string `gorm:"uniqueIndex;not null"`
	Name string `gorm:"not null"`
}

// CategoryModel represents a category in the database.
func (CategoryModel) TableName() string {
	return "categories"
}

type ProductModel struct {
	ID         uint             `gorm:"primaryKey"`
	Code       string           `gorm:"uniqueIndex;not null"`
	Price      *decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	CategoryID uint             `gorm:"not null"`
	Category   CategoryModel    `gorm:"foreignKey:CategoryID"`
}

// ProductModel represents a product in the database.
func (ProductModel) TableName() string {
	return "products"
}
