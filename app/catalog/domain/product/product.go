package product

import "github.com/mytheresa/go-hiring-challenge/app/variants/domain/product_variants"

type Product struct {
	ID       uint
	Code     string
	Price    float64
	Category Category
	Variants []product_variants.Variant
}
