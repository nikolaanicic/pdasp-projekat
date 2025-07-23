package chaincode

import (
	"chaincode/chaincode/mocks"
	"chaincode/models"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
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
	err = assetTransfer.DeleteModel(transactionContext, "")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, nil)
	err = assetTransfer.DeleteModel(transactionContext, id)
	require.Error(t, err)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = assetTransfer.DeleteModel(transactionContext, "")
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
	trader := models.Trader{
		ID:             "t1",
		PIB:            "123456",
		AccountBalance: 100,
		Receipts:       []string{},
	}
	product := models.Product{
		ID:             "p1",
		Name:           "Product",
		TraderID:       "t1",
		ExpirationDate: time.Now().Format(time.RFC3339),
		Price:          10,
		Quantity:       2,
	}

	userBytes, _ := json.Marshal(user)
	traderBytes, _ := json.Marshal(trader)
	productBytes, _ := json.Marshal(product)

	state := map[string][]byte{
		user.ID:    userBytes,
		trader.ID:  traderBytes,
		product.ID: productBytes,
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
	require.NoError(t, json.Unmarshal(state[user.ID], &updatedUser))
	require.NoError(t, json.Unmarshal(state[trader.ID], &updatedTrader))
	require.NoError(t, json.Unmarshal(state[product.ID], &updatedProduct))

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
	}

	results, err := sc.QueryProducts(ctx, filters)
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Equal(t, product.ID, results[0].ID)
	require.Equal(t, product.Name, results[0].Name)

	queryArg := stub.GetQueryResultArgsForCall(0)
	expectedQuery := `{"selector":{"name":{"$eq":"Apple"},"trader_id":{"$eq":"t1"}}}`
	require.JSONEq(t, expectedQuery, queryArg)
}
