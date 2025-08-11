package catalog

import (
	"context"
	"strconv"

	"gorm.io/gorm"
)

type ProductRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{db: db}
}

func (r *ProductRepositoryImpl) FindAll(ctx context.Context, offset int, limit int, category string, priceLt *float64) ([]Product, error) {
	var productModels []ProductModel
	query := r.db.WithContext(ctx).Preload("Category")

	if category != "" {
		query = query.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.code = ?", category)
	}
	if priceLt != nil {
		query = query.Where("price < ?", *priceLt)
	}

	if offset >= 0 {
		query = query.Offset(offset)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&productModels).Error; err != nil {
		return nil, err
	}

	products := make([]Product, len(productModels))
	for i, p := range productModels {
		products[i] = mapToDomainProduct(p)
	}

	return products, nil
}

func (r *ProductRepositoryImpl) FindByID(ctx context.Context, id uint) (*Product, error) {
	var p ProductModel
	if err := r.db.WithContext(ctx).Preload("Category").First(&p, id).Error; err != nil {
		return nil, err
	}

	domainProduct := mapToDomainProduct(p)
	return &domainProduct, nil
}

func mapToDomainProduct(p ProductModel) Product {
	var productPrice float64
	if p.Price != nil {
		productPrice = p.Price.InexactFloat64()
	} else {
		productPrice = 0
	}

	var category Category
	if p.Category.ID != 0 {
		category = Category{
			ID:   strconv.Itoa(int(p.Category.ID)),
			Code: p.Category.Code,
			Name: p.Category.Name,
		}
	} else {
		category = Category{}
	}

	return Product{
		ID:       p.ID,
		Code:     p.Code,
		Price:    productPrice,
		Category: category,
	}
}
