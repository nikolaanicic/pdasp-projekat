package chaincode

import (
	"chaincode/models"
	"encoding/json"
	"fmt"

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
