package repository

import (
	"context"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/app/catalog/domain/product"
	"github.com/mytheresa/go-hiring-challenge/app/catalog/infrastructure/persistence/models"
	"gorm.io/gorm"
)

type ProductRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{db: db}
}

func (r *ProductRepositoryImpl) FindAll(ctx context.Context, offset int, limit int, category string, priceLt *float64) ([]product.Product, error) {
	var productModels []models.ProductModel
	query := r.db.WithContext(ctx).Preload("Variants").Preload("Category")

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

	products := make([]product.Product, len(productModels))
	for i, p := range productModels {
		products[i] = mapToDomainProduct(p)
	}

	return products, nil
}

func (r *ProductRepositoryImpl) FindByID(ctx context.Context, id uint) (*product.Product, error) {
	var p models.ProductModel
	if err := r.db.WithContext(ctx).Preload("Variants").Preload("Category").First(&p, id).Error; err != nil {
		return nil, err
	}

	domainProduct := mapToDomainProduct(p)
	return &domainProduct, nil
}

func mapToDomainProduct(p models.ProductModel) product.Product {
	var variants []product.Variant
	for _, v := range p.Variants {
		var price float64
		price = 0
		if v.Price != nil {
			price = v.Price.InexactFloat64()
		}

		variants = append(variants, product.Variant{
			ID:    strconv.Itoa(int(v.ID)),
			Name:  v.Name,
			SKU:   v.SKU,
			Price: price,
		})
	}

	var productPrice float64
	if p.Price != nil {
		productPrice = p.Price.InexactFloat64()
	} else {
		productPrice = 0
	}

	var category product.Category
	if p.Category.ID != 0 {
		category = product.Category{
			ID:   strconv.Itoa(int(p.Category.ID)),
			Code: p.Category.Code,
			Name: p.Category.Name,
		}
	} else {
		category = product.Category{}
	}

	return product.Product{
		ID:       p.ID,
		Code:     p.Code,
		Price:    productPrice,
		Category: category,
		Variants: variants,
	}
}
