/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-go/chaincode"
)

func main() {
	productChaincode, err := contractapi.NewChaincode(&chaincode.SupplyChain{}) 
	if err != nil {
		log.Panicf("Error creating product status chaincode: %v", err)
	}

	if err := productChaincode.Start(); err != nil {
		log.Panicf("Error starting product status chaincode: %v", err)
	}
}
