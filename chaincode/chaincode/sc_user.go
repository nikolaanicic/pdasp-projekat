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

func (s *SmartContract) QueryUsers(ctx contractapi.TransactionContextInterface, queryString string) ([]*models.User, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var users []*models.User
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var person models.User
		err = json.Unmarshal(queryResponse.Value, &person)
		if err != nil {
			return nil, err
		}
		users = append(users, &person)
	}

	return users, nil
}

func (s *SmartContract) SearchUsersByName(ctx contractapi.TransactionContextInterface, nameQuery string) ([]*models.User, error) {
	queryString := models.BuildQueryFieldContains(models.USER_TYPE, "name", nameQuery)
	return s.QueryUsers(ctx, queryString)
}

func (s *SmartContract) SearchUsersByLastName(ctx contractapi.TransactionContextInterface, lastNameQuery string) ([]*models.User, error) {
	queryString := models.BuildQueryFieldContains(models.USER_TYPE, "last_name", lastNameQuery)
	return s.QueryUsers(ctx, queryString)
}

func (s *SmartContract) SearchUsersByLastNameAndEmail(ctx contractapi.TransactionContextInterface, lastname string, email string) ([]*models.User, error) {
	surnameSelector := models.BuildContainsSelector("last_name", lastname)
	emailSelector := models.BuildContainsSelector("email", email)
	selectors := fmt.Sprintf("%s, %s", surnameSelector, emailSelector)
	queryString := models.BuildQueryForEntityType(models.USER_TYPE, selectors)
	return s.QueryUsers(ctx, queryString)
}

func (sc *SmartContract) GetUsersGTEBalance(ctx contractapi.TransactionContextInterface, balance uint) ([]*models.User, error) {
	queryString := `
	{
	"selector": {
		"$and": [
		{
			"account_balance": {
			"$gte": ` + fmt.Sprintf("%d", balance) + `
			}
		},
		{
			"id": {
			"$regex": "^(USER)"
			}
		}
		]
	}
	}`
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)

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
