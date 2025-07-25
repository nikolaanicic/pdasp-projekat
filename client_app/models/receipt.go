package models

import "time"

type Receipt struct {
	ID        string    `json:"id"`
	TraderID  string    `json:"trader"`
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	Date      time.Time `json:"date"`
}

func (r Receipt) GetID() string {
	return r.ID
}
