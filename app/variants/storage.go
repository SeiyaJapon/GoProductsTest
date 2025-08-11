package variants

import (
	"context"
	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"strconv"

	"gorm.io/gorm"
)

type VariantORM struct {
	db          *gorm.DB
	productRepo catalog.Repository
}

func NewVariantRepository(db *gorm.DB, productRepo catalog.Repository) *VariantORM {
	return &VariantORM{db: db, productRepo: productRepo}
}

func (r *VariantORM) FindByID(ctx context.Context, id uint) (*ProductWithVariants, error) {
	product, err := r.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	variants, err := r.FindVariantsByProductID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &ProductWithVariants{
		Product:  product,
		Variants: variants,
	}, nil
}

func (r *VariantORM) FindVariantsByProductID(ctx context.Context, productID uint) ([]Variant, error) {
	var variantModels []VariantModel
	if err := r.db.WithContext(ctx).Where("product_id = ?", productID).Find(&variantModels).Error; err != nil {
		return nil, err
	}

	v := make([]Variant, len(variantModels))
	for i, vm := range variantModels {
		var price float64
		if vm.Price != nil {
			price = vm.Price.InexactFloat64()
		}
		v[i] = Variant{
			ID:    strconv.Itoa(int(vm.ID)),
			Name:  vm.Name,
			SKU:   vm.SKU,
			Price: price,
		}
	}

	return v, nil
}
