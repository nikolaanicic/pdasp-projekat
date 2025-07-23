package models

type Model interface {
	Product | User | Trader | Receipt

	GetID() string
}
