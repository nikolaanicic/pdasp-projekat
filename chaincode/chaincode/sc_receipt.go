package chaincode

import (
	"chaincode/models"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (sc *SmartContract) ReadReceipt(ctx contractapi.TransactionContextInterface, id string) (*models.Receipt, error) {
	return readModel[models.Receipt](ctx, models.ToReceiptID(id))
}

func (sc *SmartContract) DeleteReceipt(ctx contractapi.TransactionContextInterface, id string) error {
	return deleteModel(ctx, models.ToReceiptID(id))
}

func (sc *SmartContract) CreateReceipt(ctx contractapi.TransactionContextInterface, receipt models.Receipt) error {
	receipt.ID = models.ToReceiptID(receipt.ID)
	return createModel(ctx, receipt)
}

func (sc *SmartContract) GetAllReceips(ctx contractapi.TransactionContextInterface) ([]*models.Receipt, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(models.BuildQueryIdStartsWith(models.RECEIPT_TYPE))
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*models.Receipt
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset models.Receipt
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
