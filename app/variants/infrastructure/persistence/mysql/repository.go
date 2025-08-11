package mysql

import (
	"context"
	domaincatalog "github.com/mytheresa/go-hiring-challenge/app/catalog/domain"
	"github.com/mytheresa/go-hiring-challenge/app/variants/domain"
	"github.com/mytheresa/go-hiring-challenge/app/variants/infrastructure/persistence"
	"strconv"

	"gorm.io/gorm"
)

type VariantORM struct {
	db          *gorm.DB
	productRepo domaincatalog.Repository
}

func NewVariantRepository(db *gorm.DB, productRepo domaincatalog.Repository) *VariantORM {
	return &VariantORM{db: db, productRepo: productRepo}
}

// FindByID retrieves a product by its ID along with its variants.
// It returns a ProductWithVariants and an error if the product is not found or any other error occurs.
func (r *VariantORM) FindByID(ctx context.Context, id uint) (*domain.ProductWithVariants, error) {
	product, err := r.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	variants, err := r.FindVariantsByProductID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.ProductWithVariants{
		Product:  product,
		Variants: variants,
	}, nil
}

// FindVariantsByProductID retrieves all variants for a given product ID.
// It returns a slice of Variant and an error if any.
func (r *VariantORM) FindVariantsByProductID(ctx context.Context, productID uint) ([]domain.Variant, error) {
	var variantModels []persistence.VariantModel
	if err := r.db.WithContext(ctx).Where("product_id = ?", productID).Find(&variantModels).Error; err != nil {
		return nil, err
	}

	return mapVariantModelsToDomain(variantModels), nil
}

func mapVariantModelsToDomain(variantModels []persistence.VariantModel) []domain.Variant {
	v := make([]domain.Variant, len(variantModels))
	for i, vm := range variantModels {
		var price float64
		if vm.Price != nil {
			price = vm.Price.InexactFloat64()
		}
		v[i] = domain.Variant{
			ID:    strconv.Itoa(int(vm.ID)),
			Name:  vm.Name,
			SKU:   vm.SKU,
			Price: price,
		}
	}
	return v
}
