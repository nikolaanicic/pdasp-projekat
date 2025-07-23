package models

type Product struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ExpirationDate string `json:"expiration_date"`
	Price          uint   `json:"price"`
	Quantity       uint   `json:"quantity"`
	TraderID       string `json:"trader_id"`
}

func (p Product) GetID() string {
	return p.ID
}
