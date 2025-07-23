package models

import (
	"time"

	"github.com/google/uuid"
)

type InitialChainState struct {
	Products []Product
	Traders  []Trader
	Users    []User
}

func getIds[T Model](models []T) []string {
	ids := make([]string, 0)

	for _, model := range models {
		ids = append(ids, model.GetID())
	}

	return ids
}

func GetInitialChainState() InitialChainState {

	marketProducts := []Product{
		{ID: uuid.NewString(), Name: "Tomato", ExpirationDate: time.Now().Add(10 * 24 * time.Hour).UTC().Format(time.RFC3339), Price: 2, Quantity: 10},
		{ID: uuid.NewString(), Name: "Bread", ExpirationDate: time.Now().Add(2 * 24 * time.Hour).UTC().Format(time.RFC3339), Price: 3, Quantity: 10},
		{ID: uuid.NewString(), Name: "Cucumber", ExpirationDate: time.Now().Add(10 * 24 * time.Hour).UTC().Format(time.RFC3339), Price: 2, Quantity: 10},
		{ID: uuid.NewString(), Name: "Milk", ExpirationDate: time.Now().Add(30 * 24 * time.Hour).UTC().Format(time.RFC3339), Price: 3, Quantity: 10},
	}

	autoParts := []Product{
		{ID: uuid.NewString(), Name: "Steering Wheel", Price: 5, Quantity: 10},
		{ID: uuid.NewString(), Name: "Tire", Price: 8, Quantity: 10},
		{ID: uuid.NewString(), Name: "Gearbox", Price: 20, Quantity: 10},
	}

	motoParts := []Product{
		{ID: uuid.NewString(), Name: "Side Mirrors", Price: 6, Quantity: 10},
		{ID: uuid.NewString(), Name: "Pillion Seat Cover", Price: 4, Quantity: 10},
		{ID: uuid.NewString(), Name: "Braking Pads", Price: 5, Quantity: 10},
	}

	allProducts := append(marketProducts, autoParts...)
	allProducts = append(allProducts, motoParts...)

	traders := []Trader{
		{ID: uuid.NewString(), TraderType: Market, PIB: uuid.NewString(), Products: getIds(marketProducts)},
		{ID: uuid.NewString(), TraderType: AutoParts, PIB: uuid.NewString(), Products: getIds(autoParts)},
		{ID: uuid.NewString(), TraderType: MotorcycleParts, PIB: uuid.NewString(), Products: getIds(motoParts)},
	}

	users := []User{
		{ID: uuid.NewString(), Name: "Jon", LastName: "Jones", Email: "duck@jonjones.com", AccountBalance: 0},
		{ID: uuid.NewString(), Name: "Ilia", LastName: "Topuria", Email: "copycat@connor.com", AccountBalance: 0},
		{ID: uuid.NewString(), Name: "Oleksandr", LastName: "Usyk", Email: "heavy.goat@box.com", AccountBalance: 1000},
	}

	return InitialChainState{Products: allProducts, Traders: traders, Users: users}
}
