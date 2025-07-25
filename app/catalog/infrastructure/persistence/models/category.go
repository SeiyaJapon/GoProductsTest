package models

type CategoryModel struct {
	ID   uint   `gorm:"primaryKey"`
	Code string `gorm:"uniqueIndex;not null"`
	Name string `gorm:"not null"`
}

func (CategoryModel) TableName() string {
	return "categories"
}
