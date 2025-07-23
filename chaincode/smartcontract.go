package chaincode

import (
	"chain/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

func putWorldState[T models.Model](models []T, ctx contractapi.TransactionContextInterface) error {

	for _, model := range models {
		modelJson, err := json.Marshal(model)
		if err != nil {
			return err
		}

		if err := ctx.GetStub().PutState(model.GetID(), modelJson); err != nil {
			return fmt.Errorf("failed to put an asset into the world state: id:%v err:%v", model.GetID(), err)
		}
	}

	return nil
}

func (sc *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	initialState := models.GetInitialChainState()

	if err := putWorldState(initialState.Products, ctx); err != nil {
		return err
	}

	if err := putWorldState(initialState.Traders, ctx); err != nil {
		return err
	}

	if err := putWorldState(initialState.Users, ctx); err != nil {
		return err
	}

	return nil
}

func modelExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	modelJson, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, err
	}

	return modelJson != nil, nil
}

func (sc *SmartContract) ModelExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	return modelExists(ctx, id)
}

func createModel[T models.Model](ctx contractapi.TransactionContextInterface, model T) error {
	exists, err := modelExists(ctx, model.GetID())
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("asset with id:%s already exists", model.GetID())
	}

	modelJson, err := json.Marshal(model)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(model.GetID(), modelJson)
}

func readModel[T models.Model](ctx contractapi.TransactionContextInterface, id string) (*T, error) {
	modelJson, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read the model: %v", err)
	}

	if modelJson == nil {
		return nil, fmt.Errorf("the model with the id: %s doesn't exist", id)
	}

	var model T

	if err := json.Unmarshal(modelJson, &model); err != nil {
		return nil, fmt.Errorf("failed to deserialize the model: %v", err)
	}

	return &model, nil
}

func updateModel[T models.Model](ctx contractapi.TransactionContextInterface, id string, model *T) error {
	exists, err := modelExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the model %s does not exist", id)
	}

	modelJson, err := json.Marshal(model)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, modelJson)
}

func deleteModel(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := modelExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the model %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

func (sc *SmartContract) ReadProduct(ctx contractapi.TransactionContextInterface, id string) (*models.Product, error) {
	return readModel[models.Product](ctx, id)
}

func (sc *SmartContract) ReadTrader(ctx contractapi.TransactionContextInterface, id string) (*models.Trader, error) {
	return readModel[models.Trader](ctx, id)
}

func (sc *SmartContract) ReadUser(ctx contractapi.TransactionContextInterface, id string) (*models.User, error) {
	return readModel[models.User](ctx, id)
}

func (sc *SmartContract) ReadReceipt(ctx contractapi.TransactionContextInterface, id string) (*models.Receipt, error) {
	return readModel[models.Receipt](ctx, id)
}

func (sc *SmartContract) UpdateProduct(ctx contractapi.TransactionContextInterface, id string, model *models.Product) error {
	return updateModel(ctx, id, model)
}
func (sc *SmartContract) UpdateUser(ctx contractapi.TransactionContextInterface, id string, model *models.User) error {
	return updateModel(ctx, id, model)
}
func (sc *SmartContract) UpdateTrader(ctx contractapi.TransactionContextInterface, id string, model *models.Trader) error {
	return updateModel(ctx, id, model)
}

func (sc *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, user models.User) error {
	return createModel(ctx, user)
}
func (sc *SmartContract) CreateProduct(ctx contractapi.TransactionContextInterface, product models.Product) error {
	return createModel(ctx, product)
}
func (sc *SmartContract) CreateTrader(ctx contractapi.TransactionContextInterface, trader models.Trader) error {
	return createModel(ctx, trader)
}
func (sc *SmartContract) CreateReceipt(ctx contractapi.TransactionContextInterface, receipt models.Receipt) error {
	return createModel(ctx, receipt)
}

func (sc *SmartContract) DeleteModel(ctx contractapi.TransactionContextInterface, id string) error {
	return deleteModel(ctx, id)
}

func (sc *SmartContract) BuyProduct(ctx contractapi.TransactionContextInterface, productId string, userId string) error {
	user, err := sc.ReadUser(ctx, userId)
	if err != nil {
		return err
	}

	product, err := sc.ReadProduct(ctx, productId)
	if err != nil {
		return err
	}

	trader, err := sc.ReadTrader(ctx, product.TraderID)
	if err != nil {
		return err
	}

	if product.Price > user.AccountBalance {
		return fmt.Errorf("user doesn't have enough funds to buy the product")
	}

	product.Quantity -= 1
	user.AccountBalance -= product.Price
	trader.AccountBalance += product.Price

	receipt := models.Receipt{
		ID:        uuid.NewString(),
		TraderID:  product.TraderID,
		UserID:    userId,
		ProductID: productId,
		Date:      time.Now().UTC(),
	}

	if product.Quantity == 0 {
		if err := sc.DeleteModel(ctx, product.ID); err != nil {
			return fmt.Errorf("failed to remove the product: %v", err)
		}
	}

	user.ReceiptsID = append(user.ReceiptsID, receipt.ID)
	trader.Receipts = append(trader.Receipts, receipt.ID)

	if err := sc.CreateReceipt(ctx, receipt); err != nil {
		return err
	}

	if err := sc.UpdateProduct(ctx, productId, product); err != nil {
		return err
	}

	if err := sc.UpdateTrader(ctx, trader.ID, trader); err != nil {
		return err
	}

	if err := sc.UpdateUser(ctx, user.ID, user); err != nil {
		return err
	}

	return nil
}

func (sc *SmartContract) GetAllProducts(ctx contractapi.TransactionContextInterface) ([]*models.Product, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*models.Product
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset models.Product
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func (sc *SmartContract) QueryProducts(ctx contractapi.TransactionContextInterface, filters map[string]string) ([]*models.Product, error) {
	selector := make(map[string]map[string]string)

	for key, value := range filters {
		selector[key] = map[string]string{"$eq": value}
	}

	query := map[string]interface{}{
		"selector": selector,
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	resultsIterator, err := ctx.GetStub().GetQueryResult(string(queryBytes))
	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()

	var products []*models.Product

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var product models.Product
		if err := json.Unmarshal(response.Value, &product); err != nil {
			return nil, err
		}

		products = append(products, &product)
	}

	return products, nil
}
