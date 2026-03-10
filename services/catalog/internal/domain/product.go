package domain

type Product struct {
	ID          string
	SellerID    string
	CategoryID  string
	Title       string
	Description string
	Price       string
	Images      []string
	CreatedAt   string
	UpdatedAt   string
}
