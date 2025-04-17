/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	_ "github.com/XiaoYao-austin/goRecrypt/recrypt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"authentication-center/chaincode"
	"log"
)

func main() {

	authenticationCenterChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating authentication-center chaincode: %v", err)
	}

	if err := authenticationCenterChaincode.Start(); err != nil {
		log.Panicf("Error starting authentication-center chaincode: %v", err)
	}
}
