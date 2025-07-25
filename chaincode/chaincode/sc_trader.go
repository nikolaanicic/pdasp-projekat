package chaincode

import (
	"chaincode/models"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (sc *SmartContract) ReadTrader(ctx contractapi.TransactionContextInterface, id string) (*models.Trader, error) {
	return readModel[models.Trader](ctx, models.ToTraderID(id))
}

func (sc *SmartContract) UpdateTrader(ctx contractapi.TransactionContextInterface, id string, model *models.Trader) error {
	return updateModel(ctx, models.ToTraderID(id), model)
}
func (sc *SmartContract) DeleteTrader(ctx contractapi.TransactionContextInterface, id string) error {
	return deleteModel(ctx, models.ToTraderID(id))
}

func (sc *SmartContract) CreateTrader(ctx contractapi.TransactionContextInterface, trader models.Trader) error {
	trader.ID = models.ToTraderID(trader.ID)
	trader.Receipts = make([]string, 0)
	trader.Products = make([]string, 0)

	return createModel(ctx, trader)
}

func (sc *SmartContract) GetAllTraders(ctx contractapi.TransactionContextInterface) ([]*models.Trader, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(models.BuildQueryIdStartsWith(models.TRADER_TYPE))
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*models.Trader
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset models.Trader
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
