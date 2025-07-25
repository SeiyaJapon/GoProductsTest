package product

type Product struct {
	ID       uint
	Code     string
	Price    float64
	Category Category
	Variants []Variant
}
