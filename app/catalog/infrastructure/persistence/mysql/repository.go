package mysql

import (
	"context"
	"github.com/mytheresa/go-hiring-challenge/app/catalog/domain"
	"github.com/mytheresa/go-hiring-challenge/app/catalog/infrastructure/persistence"
	"strconv"

	"gorm.io/gorm"
)

type ProductRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{db: db}
}

// FindAll retrieves products based on the provided criteria.
// It returns a slice of Product and an error if any.
func (r *ProductRepositoryImpl) FindAll(ctx context.Context, offset int, limit int, category string, priceLt *float64) ([]domain.Product, error) {
	var productModels []persistence.ProductModel
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

	products := make([]domain.Product, len(productModels))
	for i, p := range productModels {
		products[i] = mapToDomainProduct(p)
	}

	return products, nil
}

// FindByID retrieves a product by its ID.
// It returns a pointer to Product and an error if the product is not found or any other error occurs.
func (r *ProductRepositoryImpl) FindByID(ctx context.Context, id uint) (*domain.Product, error) {
	var p persistence.ProductModel
	if err := r.db.WithContext(ctx).Preload("Category").First(&p, id).Error; err != nil {
		return nil, err
	}

	domainProduct := mapToDomainProduct(p)
	return &domainProduct, nil
}

func mapToDomainProduct(p persistence.ProductModel) domain.Product {
	var productPrice float64
	if p.Price != nil {
		productPrice = p.Price.InexactFloat64()
	} else {
		productPrice = 0
	}

	var category domain.Category
	if p.Category.ID != 0 {
		category = domain.Category{
			ID:   strconv.Itoa(int(p.Category.ID)),
			Code: p.Category.Code,
			Name: p.Category.Name,
		}
	} else {
		category = domain.Category{}
	}

	return domain.Product{
		ID:       p.ID,
		Code:     p.Code,
		Price:    productPrice,
		Category: category,
	}
}
