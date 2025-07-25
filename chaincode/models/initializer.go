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
		{ID: ToProductID("t1"), Name: "Tomato", ExpirationDate: time.Now().Add(10 * 24 * time.Hour).UTC().Format("02-01-2006"), Price: 2, Quantity: 10, TraderID: "tt1"},
		{ID: ToProductID("b1"), Name: "Bread", ExpirationDate: time.Now().Add(2 * 24 * time.Hour).UTC().Format("02-01-2006"), Price: 3, Quantity: 10, TraderID: "tt1"},
		{ID: ToProductID("c1"), Name: "Cucumber", ExpirationDate: time.Now().Add(10 * 24 * time.Hour).UTC().Format("02-01-2006"), Price: 2, Quantity: 10, TraderID: "tt1"},
		{ID: ToProductID("m1"), Name: "Milk", ExpirationDate: time.Now().Add(30 * 24 * time.Hour).UTC().Format("02-01-2006"), Price: 3, Quantity: 10, TraderID: "tt1"},
	}

	autoParts := []Product{
		{ID: ToProductID("sw1"), Name: "Steering Wheel", Price: 5, Quantity: 10, TraderID: "tt2"},
		{ID: ToProductID("ti1"), Name: "Tire", Price: 8, Quantity: 10, TraderID: "tt2"},
		{ID: ToProductID("ge1"), Name: "Gearbox", Price: 20, Quantity: 10, TraderID: "tt2"},
	}

	motoParts := []Product{
		{ID: ToProductID("si1"), Name: "Side Mirrors", Price: 6, Quantity: 10, TraderID: "tt3"},
		{ID: ToProductID("pi1"), Name: "Pillion Seat Cover", Price: 4, Quantity: 10, TraderID: "tt3"},
		{ID: ToProductID("br1"), Name: "Braking Pads", Price: 5, Quantity: 10, TraderID: "tt3"},
	}

	allProducts := append(marketProducts, autoParts...)
	allProducts = append(allProducts, motoParts...)

	traders := []Trader{
		{ID: ToTraderID("tt1"), TraderType: Market, PIB: "pib1", Products: getIds(marketProducts), Receipts: make([]string, 0)},
		{ID: ToTraderID("tt2"), TraderType: AutoParts, PIB: "pib2", Products: getIds(autoParts), Receipts: make([]string, 0)},
		{ID: ToTraderID("tt3"), TraderType: MotorcycleParts, PIB: "pib3", Products: getIds(motoParts), Receipts: make([]string, 0)},
	}

	users := []User{
		{ID: ToUserID("jj1"), Name: "Jon", LastName: "Jones", Email: "duck@jonjones.com", AccountBalance: 0, ReceiptsID: make([]string, 0)},
		{ID: ToUserID("it1"), Name: "Ilia", LastName: "Topuria", Email: "copycat@connor.com", AccountBalance: 0, ReceiptsID: make([]string, 0)},
		{ID: ToUserID("ou1"), Name: "Oleksandr", LastName: "Usyk", Email: "heavy.goat@box.com", AccountBalance: 1000, ReceiptsID: make([]string, 0)},
	}

	return InitialChainState{Products: allProducts, Traders: traders, Users: users}
}
