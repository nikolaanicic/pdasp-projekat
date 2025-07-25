package models

type User struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	LastName       string   `json:"last_name"`
	Email          string   `json:"email"`
	ReceiptsID     []string `json:"receipts_ids"`
	AccountBalance uint     `json:"account_balance"`
}

func (p User) GetID() string {
	return p.ID
}
