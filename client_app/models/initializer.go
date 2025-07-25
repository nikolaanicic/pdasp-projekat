package models

import (
	"time"
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
		{ID: "t1", Name: "Tomato", ExpirationDate: time.Now().Add(10 * 24 * time.Hour).UTC().Format(time.RFC3339), Price: 2, Quantity: 10},
		{ID: "b1", Name: "Bread", ExpirationDate: time.Now().Add(2 * 24 * time.Hour).UTC().Format(time.RFC3339), Price: 3, Quantity: 10},
		{ID: "c1", Name: "Cucumber", ExpirationDate: time.Now().Add(10 * 24 * time.Hour).UTC().Format(time.RFC3339), Price: 2, Quantity: 10},
		{ID: "m1", Name: "Milk", ExpirationDate: time.Now().Add(30 * 24 * time.Hour).UTC().Format(time.RFC3339), Price: 3, Quantity: 10},
	}

	autoParts := []Product{
		{ID: "sw1", Name: "Steering Wheel", Price: 5, Quantity: 10},
		{ID: "ti1", Name: "Tire", Price: 8, Quantity: 10},
		{ID: "ge1", Name: "Gearbox", Price: 20, Quantity: 10},
	}

	motoParts := []Product{
		{ID: "si1", Name: "Side Mirrors", Price: 6, Quantity: 10},
		{ID: "pi1", Name: "Pillion Seat Cover", Price: 4, Quantity: 10},
		{ID: "br1", Name: "Braking Pads", Price: 5, Quantity: 10},
	}

	allProducts := append(marketProducts, autoParts...)
	allProducts = append(allProducts, motoParts...)

	traders := []Trader{
		{ID: "tt1", TraderType: Market, PIB: "pib1", Products: getIds(marketProducts)},
		{ID: "tt2", TraderType: AutoParts, PIB: "pib2", Products: getIds(autoParts)},
		{ID: "tt3", TraderType: MotorcycleParts, PIB: "pib3", Products: getIds(motoParts)},
	}

	users := []User{
		{ID: "jj1", Name: "Jon", LastName: "Jones", Email: "duck@jonjones.com", AccountBalance: 0},
		{ID: "it1", Name: "Ilia", LastName: "Topuria", Email: "copycat@connor.com", AccountBalance: 0},
		{ID: "ou1", Name: "Oleksandr", LastName: "Usyk", Email: "heavy.goat@box.com", AccountBalance: 1000},
	}

	return InitialChainState{Products: allProducts, Traders: traders, Users: users}
}
