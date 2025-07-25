package chaincode

import (
	"chaincode/models"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (sc *SmartContract) GetAllProducts(ctx contractapi.TransactionContextInterface) ([]*models.Product, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(models.BuildQueryIdStartsWith(models.PRODUCT_TYPE))
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

func (sc *SmartContract) QueryProducts(ctx contractapi.TransactionContextInterface, filters map[string]string) ([]*models.Product, error) {
	selector := make(map[string]interface{})

	for key, value := range filters {
		// Ako je filter za cenu, koristi numeričko poređenje
		if key == "price" {
			priceVal, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid price value: %v", err)
			}
			selector["price"] = map[string]interface{}{"$eq": priceVal}
		} else {
			// Za ostale koristi direktno poređenje
			selector[key] = map[string]interface{}{"$eq": value}
		}
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
		ID:        fmt.Sprintf("%s-%s-%s-%d", user.ID, product.TraderID, product.ID, len(user.ReceiptsID)),
		TraderID:  product.TraderID,
		UserID:    userId,
		ProductID: productId,
		Date:      time.Now().UTC().Format("02-01-2006"),
	}

	if product.Quantity == 0 {
		if err := sc.DeleteProduct(ctx, productId); err != nil {
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

	if err := sc.UpdateTrader(ctx, product.TraderID, trader); err != nil {
		return err
	}

	if err := sc.UpdateUser(ctx, userId, user); err != nil {
		return err
	}

	return nil
}

func (sc *SmartContract) CreateProduct(ctx contractapi.TransactionContextInterface, product models.Product) error {
	product.ID = models.ToProductID(product.ID)
	return createModel(ctx, product)
}

func (sc *SmartContract) UpdateProduct(ctx contractapi.TransactionContextInterface, id string, model *models.Product) error {
	return updateModel(ctx, models.ToProductID(id), model)
}

func (sc *SmartContract) DeleteProduct(ctx contractapi.TransactionContextInterface, id string) error {
	return deleteModel(ctx, models.ToProductID(id))
}

func (sc *SmartContract) ReadProduct(ctx contractapi.TransactionContextInterface, id string) (*models.Product, error) {
	return readModel[models.Product](ctx, models.ToProductID(id))
}
