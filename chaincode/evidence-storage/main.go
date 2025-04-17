/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	_ "github.com/XiaoYao-austin/goRecrypt/recrypt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"evidence-storage/chaincode"
	"log"
)

func main() {

	evidenceStorageChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating evidence-storage chaincode: %v", err)
	}

	if err := evidenceStorageChaincode.Start(); err != nil {
		log.Panicf("Error starting evidence-storage chaincode: %v", err)
	}
}
