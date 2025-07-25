package chaincode

import (
	"chaincode/models"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (sc *SmartContract) ReadUser(ctx contractapi.TransactionContextInterface, id string) (*models.User, error) {
	return readModel[models.User](ctx, models.ToUserID(id))
}

func (sc *SmartContract) UpdateUser(ctx contractapi.TransactionContextInterface, id string, model *models.User) error {
	return updateModel(ctx, models.ToUserID(id), model)
}

func (sc *SmartContract) DeleteUser(ctx contractapi.TransactionContextInterface, id string) error {
	return deleteModel(ctx, models.ToUserID(id))
}

func (sc *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, user models.User) error {
	user.ID = models.ToUserID(user.ID)
	user.ReceiptsID = make([]string, 0)

	return createModel(ctx, user)
}

func (sc *SmartContract) GetAllUsers(ctx contractapi.TransactionContextInterface) ([]*models.User, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(models.BuildQueryIdStartsWith(models.USER_TYPE))
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*models.User
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset models.User
		if queryResponse == nil {
			return nil, fmt.Errorf("no next")
		}
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}
