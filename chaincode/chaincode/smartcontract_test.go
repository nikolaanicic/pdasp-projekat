package chaincode

import (
	"chaincode/chaincode/mocks"
	"chaincode/models"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/stretchr/testify/require"
)

func TestInitLedgerProducts(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	assetTransfer := SmartContract{}
	err := assetTransfer.InitLedger(transactionContext)
	require.NoError(t, err)

	chaincodeStub.PutStateReturns(fmt.Errorf("failed inserting key"))
	err = assetTransfer.InitLedger(transactionContext)
	require.Error(t, err)
}

func TestCreateModel(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	assetTransfer := SmartContract{}
	id := uuid.NewString()
	err := assetTransfer.CreateProduct(transactionContext, models.Product{ID: id, Name: "p1", ExpirationDate: time.Now().Format(time.RFC3339), Price: 1, Quantity: 1})
	require.NoError(t, err)

	chaincodeStub.GetStateReturns([]byte{}, nil)
	err = assetTransfer.CreateProduct(transactionContext, models.Product{ID: id, Name: "p2", ExpirationDate: time.Now().Format(time.RFC3339), Price: 1, Quantity: 1})
	require.Error(t, err)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = assetTransfer.CreateProduct(transactionContext, models.Product{ID: id, Name: "p3", ExpirationDate: time.Now().Format(time.RFC3339), Price: 1, Quantity: 1})
	require.Error(t, err)
}

func TestReadModel(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	expectedAsset := &models.Product{ID: "asset1"}
	bytes, err := json.Marshal(expectedAsset)
	require.NoError(t, err)

	assetTransfer := SmartContract{}
	chaincodeStub.GetStateReturns(bytes, nil)
	readModel, err := assetTransfer.ReadProduct(transactionContext, "asset1")

	require.NoError(t, err)
	require.Equal(t, expectedAsset, readModel)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	_, err = assetTransfer.ReadProduct(transactionContext, "")
	require.Error(t, err)

	chaincodeStub.GetStateReturns(nil, nil)
	_, err = assetTransfer.ReadProduct(transactionContext, "asset1")
	require.Error(t, err)
}

func TestUpdateModel(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	id := "asset1"
	expectedAsset := &models.Product{ID: id}
	bytes, err := json.Marshal(expectedAsset)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	toCreate := &models.Product{ID: id}
	assetTransfer := SmartContract{}

	err = assetTransfer.UpdateProduct(transactionContext, id, toCreate)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, nil)
	err = assetTransfer.UpdateProduct(transactionContext, id, toCreate)
	require.Error(t, err)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = assetTransfer.UpdateProduct(transactionContext, id, toCreate)
	require.Error(t, err)
}

func TestDeleteModel(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	id := "asset1"
	toDelete := &models.Product{ID: id}
	bytes, err := json.Marshal(toDelete)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	chaincodeStub.DelStateReturns(nil)
	assetTransfer := SmartContract{}
	err = assetTransfer.DeleteProduct(transactionContext, "")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, nil)
	err = assetTransfer.DeleteProduct(transactionContext, id)
	require.Error(t, err)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = assetTransfer.DeleteProduct(transactionContext, "")
	require.Error(t, err)
}

func TestBuyProduct(t *testing.T) {
	sc := SmartContract{}

	// Setup mock context and stub
	stub := new(mocks.ChaincodeStub)
	ctx := new(mocks.TransactionContext)
	ctx.GetStubReturns(stub)

	// Initial data
	user := models.User{
		ID:             "u1",
		Name:           "User",
		LastName:       "Test",
		Email:          "user@test.com",
		AccountBalance: 100,
		ReceiptsID:     []string{},
	}

	product := models.Product{
		ID:             "p1",
		Name:           "Product",
		TraderID:       "t1",
		ExpirationDate: time.Now().Format(time.RFC3339),
		Price:          10,
		Quantity:       2,
	}

	storedUser := models.User{
		ID:             "USER-u1",
		Name:           "User",
		LastName:       "Test",
		Email:          "user@test.com",
		AccountBalance: 100,
		ReceiptsID:     []string{},
	}
	storedProduct := models.Product{
		ID:             "PRODUCT-p1",
		Name:           "Product",
		TraderID:       "t1",
		ExpirationDate: time.Now().Format(time.RFC3339),
		Price:          10,
		Quantity:       2,
	}
	storedTrader := models.Trader{
		ID:             "TRADER-t1",
		PIB:            "123456",
		AccountBalance: 100,
		Receipts:       []string{},
	}

	userBytes, _ := json.Marshal(storedUser)
	traderBytes, _ := json.Marshal(storedTrader)
	productBytes, _ := json.Marshal(storedProduct)

	state := map[string][]byte{
		storedUser.ID:    userBytes,
		storedTrader.ID:  traderBytes,
		storedProduct.ID: productBytes,
	}

	stub.GetStateStub = func(key string) ([]byte, error) {
		return state[key], nil
	}
	stub.PutStateStub = func(key string, value []byte) error {
		state[key] = value
		return nil
	}
	stub.DelStateStub = func(key string) error {
		delete(state, key)
		return nil
	}

	err := sc.BuyProduct(ctx, product.ID, user.ID)
	require.NoError(t, err)

	var updatedUser models.User
	var updatedTrader models.Trader
	var updatedProduct models.Product
	require.NoError(t, json.Unmarshal(state[storedUser.ID], &updatedUser))
	require.NoError(t, json.Unmarshal(state[storedTrader.ID], &updatedTrader))
	require.NoError(t, json.Unmarshal(state[storedProduct.ID], &updatedProduct))

	require.Equal(t, uint(90), updatedUser.AccountBalance)
	require.Equal(t, uint(110), updatedTrader.AccountBalance)
	require.Equal(t, uint(1), updatedProduct.Quantity)
	require.Len(t, updatedUser.ReceiptsID, 1)
	require.Len(t, updatedTrader.Receipts, 1)
}

func TestQueryProducts(t *testing.T) {
	sc := SmartContract{}

	product := models.Product{
		ID:             "p1",
		Name:           "Apple",
		ExpirationDate: time.Now().Format(time.RFC3339),
		Price:          10,
		Quantity:       5,
		TraderID:       "t1",
	}

	productBytes, err := json.Marshal(product)
	require.NoError(t, err)

	stub := new(mocks.ChaincodeStub)
	ctx := new(mocks.TransactionContext)
	iterator := new(mocks.StateQueryIterator)

	ctx.GetStubReturns(stub)

	stub.GetQueryResultReturns(iterator, nil)

	iterator.HasNextReturnsOnCall(0, true)
	iterator.HasNextReturnsOnCall(1, false)

	iterator.NextReturnsOnCall(0, &queryresult.KV{
		Key:   product.ID,
		Value: productBytes,
	}, nil)

	filters := map[string]string{
		"name":      "Apple",
		"trader_id": "t1",
		"price":     "10.0",
	}

	results, err := sc.QueryProducts(ctx, filters)
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Equal(t, product.ID, results[0].ID)
	require.Equal(t, product.Name, results[0].Name)

	queryArg := stub.GetQueryResultArgsForCall(0)
	expectedQuery := `{"selector":{"name":{"$eq":"Apple"},"trader_id":{"$eq":"t1"}, "price":{"$eq":10}}}`
	require.JSONEq(t, expectedQuery, queryArg)
}

func TestGetAllProducts(t *testing.T) {
	sc := SmartContract{}
	stub := new(mocks.ChaincodeStub)
	ctx := new(mocks.TransactionContext)
	iterator := new(mocks.StateQueryIterator)

	products := []models.Product{
		{ID: "PRODUCT-p1", Name: "p1", Price: 10, TraderID: "t1"},
		{ID: "PRODUCT-p2", Name: "p2", Price: 10, TraderID: "t2"},
		{ID: "PRODUCT-p3", Name: "p3", Price: 10, TraderID: "t3"},
	}

	for i, obj := range products {
		bts, err := json.Marshal(obj)
		require.NoError(t, err)

		iterator.HasNextReturnsOnCall(i, true)
		iterator.NextReturnsOnCall(i, &queryresult.KV{
			Key:   obj.ID,
			Value: bts,
		}, nil)
	}
	iterator.HasNextReturnsOnCall(len(products), false)

	stub.GetQueryResultStub = func(query string) (shim.StateQueryIteratorInterface, error) {
		return iterator, nil
	}

	ctx.GetStubReturns(stub)

	// Run function under test
	results, err := sc.GetAllProducts(ctx)

	require.NoError(t, err)
	require.Len(t, results, 3)

	// Ensure each returned item is a valid PRODUCT
	for i := 0; i < 3; i++ {
		require.True(t, strings.HasPrefix(results[i].ID, "PRODUCT"))
		require.Equal(t, products[i].ID, results[i].ID)
		require.Equal(t, products[i].Name, results[i].Name)
	}
}

func TestGetAllProducts_FiltersOnlyProducts(t *testing.T) {
	sc := SmartContract{}
	stub := new(mocks.ChaincodeStub)
	ctx := new(mocks.TransactionContext)

	// Prepare test data: products and a non-product entry
	products := []models.Product{
		{ID: "PRODUCT-p1", Name: "p1", Price: 10, TraderID: "t1"},
		{ID: "PRODUCT-p2", Name: "p2", Price: 20, TraderID: "t2"},
	}
	nonProduct := models.Product{ID: "TRADER-t1", Name: "not a product", Price: 0, TraderID: "t0"}

	state := map[string][]byte{}
	for _, p := range products {
		b, err := json.Marshal(p)
		require.NoError(t, err)
		state[p.ID] = b
	}
	b, err := json.Marshal(nonProduct)
	require.NoError(t, err)
	state[nonProduct.ID] = b

	// Setup the iterator to only return the product entries
	iterator := new(mocks.StateQueryIterator)
	callCount := 0
	keys := []string{"PRODUCT-p1", "PRODUCT-p2"}

	iterator.HasNextStub = func() bool {
		return callCount < len(keys)
	}
	iterator.NextStub = func() (*queryresult.KV, error) {
		key := keys[callCount]
		val := state[key]
		callCount++
		return &queryresult.KV{Key: key, Value: val}, nil
	}
	iterator.CloseReturns(nil)

	stub.GetQueryResultStub = func(query string) (shim.StateQueryIteratorInterface, error) {
		return iterator, nil
	}

	ctx.GetStubReturns(stub)

	results, err := sc.GetAllProducts(ctx)
	require.NoError(t, err)
	require.Len(t, results, len(products))

	for i, product := range products {
		require.Equal(t, product.ID, results[i].ID)
		require.Equal(t, product.Name, results[i].Name)
		require.Equal(t, product.Price, results[i].Price)
		require.Equal(t, product.TraderID, results[i].TraderID)
	}
}

func TestGetAllUsers_FiltersOnlyUsers(t *testing.T) {
	sc := SmartContract{}
	stub := new(mocks.ChaincodeStub)
	ctx := new(mocks.TransactionContext)

	products := []models.User{
		{ID: "USER-p1", Name: "p1", AccountBalance: 10, LastName: "t1", Email: "t2@gmail.com"},
		{ID: "USER-p2", Name: "p2", AccountBalance: 20, LastName: "t2", Email: "t2@gmail.com"},
	}
	nonProduct := models.User{ID: "TRADER-t1", Name: "not a user", AccountBalance: 0, Email: "t0"}

	state := map[string][]byte{}
	for _, p := range products {
		b, err := json.Marshal(p)
		require.NoError(t, err)
		state[p.ID] = b
	}
	b, err := json.Marshal(nonProduct)
	require.NoError(t, err)
	state[nonProduct.ID] = b

	iterator := new(mocks.StateQueryIterator)
	callCount := 0
	keys := []string{"USER-p1", "USER-p2"}

	iterator.HasNextStub = func() bool {
		return callCount < len(keys)
	}
	iterator.NextStub = func() (*queryresult.KV, error) {
		key := keys[callCount]
		val := state[key]
		callCount++
		return &queryresult.KV{Key: key, Value: val}, nil
	}
	iterator.CloseReturns(nil)

	stub.GetQueryResultStub = func(query string) (shim.StateQueryIteratorInterface, error) {
		return iterator, nil
	}

	ctx.GetStubReturns(stub)

	results, err := sc.GetAllUsers(ctx)
	require.NoError(t, err)
	require.Len(t, results, len(products))

	for i, product := range products {
		require.Equal(t, product.ID, results[i].ID)
		require.Equal(t, product.Name, results[i].Name)
		require.Equal(t, product.AccountBalance, results[i].AccountBalance)
	}
}

func TestQueryUsers(t *testing.T) {
	sc := SmartContract{}

	user := models.User{
		ID:             "USER-1",
		Name:           "John",
		LastName:       "Doe",
		Email:          "john@example.com",
		ReceiptsID:     []string{"R1", "R2"},
		AccountBalance: 100,
	}

	userBytes, err := json.Marshal(user)
	require.NoError(t, err)

	stub := new(mocks.ChaincodeStub)
	ctx := new(mocks.TransactionContext)
	iterator := new(mocks.StateQueryIterator)

	ctx.GetStubReturns(stub)
	stub.GetQueryResultReturns(iterator, nil)

	iterator.HasNextReturnsOnCall(0, true)
	iterator.HasNextReturnsOnCall(1, false)

	iterator.NextReturnsOnCall(0, &queryresult.KV{
		Key:   user.ID,
		Value: userBytes,
	}, nil)

	results, err := sc.QueryUsers(ctx, models.BuildQueryIdStartsWith(models.USER_TYPE))
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Equal(t, user.ID, results[0].ID)
	require.Equal(t, user.Email, results[0].Email)
}

func TestGetUsersGTEBalance(t *testing.T) {
	sc := SmartContract{}

	users := []models.User{
		{
			ID:             "USER-001",
			Name:           "Alice",
			LastName:       "Smith",
			Email:          "alice@example.com",
			ReceiptsID:     []string{"R1"},
			AccountBalance: 150,
		},
		{
			ID:             "USER-002",
			Name:           "Bob",
			LastName:       "Jones",
			Email:          "bob@example.com",
			ReceiptsID:     []string{"R2"},
			AccountBalance: 90,
		},
		{
			ID:             "TRADER-001",
			Name:           "Eve",
			LastName:       "Black",
			Email:          "eve@example.com",
			ReceiptsID:     nil,
			AccountBalance: 999,
		},
	}

	mockState := map[string][]byte{}
	for _, u := range users {
		data, err := json.Marshal(u)
		require.NoError(t, err)
		mockState[u.ID] = data
	}

	expectedKeys := []string{}
	for _, u := range users {
		if strings.HasPrefix(u.ID, "USER") && u.AccountBalance >= 100 {
			expectedKeys = append(expectedKeys, u.ID)
		}
	}

	call := 0
	iterator := new(mocks.StateQueryIterator)
	iterator.HasNextStub = func() bool {
		return call < len(expectedKeys)
	}
	iterator.NextStub = func() (*queryresult.KV, error) {
		key := expectedKeys[call]
		val := mockState[key]
		call++
		return &queryresult.KV{Key: key, Value: val}, nil
	}
	iterator.CloseReturns(nil)

	stub := new(mocks.ChaincodeStub)
	ctx := new(mocks.TransactionContext)

	stub.GetQueryResultStub = func(query string) (shim.StateQueryIteratorInterface, error) {
		return iterator, nil
	}
	ctx.GetStubReturns(stub)

	results, err := sc.GetUsersGTEBalance(ctx, 100)
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Equal(t, "USER-001", results[0].ID)
	require.Equal(t, uint(150), results[0].AccountBalance)

}
