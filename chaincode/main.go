package main

import (
	"chaincode/chaincode"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	traderChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating bank chaincode: %v", err)
	}

	if err := traderChaincode.Start(); err != nil {
		log.Panicf("Error starting bank chaincode: %v", err)
	}
}
