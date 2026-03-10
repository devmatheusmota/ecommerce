package domain

type Address struct {
	ID                string
	UserID            string
	Street            string
	Number            string
	Complement        string
	Neighborhood      string
	City              string
	State             string
	ZipCode           string
	Type              string // "billing" or "shipping"
	IsDefaultBilling  bool
	IsDefaultShipping bool
}
