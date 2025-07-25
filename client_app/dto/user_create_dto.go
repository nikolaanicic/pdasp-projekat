package dto

type UserCreateDto struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	AccountBalance uint   `json:"account_balance"`
}
