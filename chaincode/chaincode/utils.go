package chaincode

import (
	"chaincode/models"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (sc *SmartContract) GetEntityById(ctx contractapi.TransactionContextInterface, entityType string, id string) ([]byte, error) {
	entity, err := ctx.GetStub().GetState(models.FormatKey(entityType, id))
	if err != nil {
		return nil, fmt.Errorf("failed to find the entity: %v", err)
	}

	return entity, nil
}

func (sc *SmartContract) EntityExists(ctx contractapi.TransactionContextInterface, entityType string, id string) (bool, error) {
	itemJSON, err := ctx.GetStub().GetState(models.FormatKey(entityType, id))
	if err != nil {
		return false, fmt.Errorf("failed to find the entity: %v", err)
	}

	return itemJSON != nil, nil
}
