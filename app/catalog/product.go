package catalog

type Product struct {
	ID       uint
	Code     string
	Price    float64
	Category Category
}
