package models

type Receipt struct {
	ID        string `json:"id"`
	TraderID  string `json:"trader"`
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	Date      string `json:"date"`
}

func (r Receipt) GetID() string {
	return r.ID
}
