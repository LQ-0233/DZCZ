/*
Copyright 2021 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fabric

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

const (
	mspID                             = "Org1MSP"
	cryptoPath                        = "./organizations/peerOrganizations/org1.example.com"
	certPath                          = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem"
	keyPath                           = cryptoPath + "/users/User1@org1.example.com/msp/keystore/"
	tlsCertPath                       = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	gatewayPeer                       = "peer0.org1.example.com"
	peerEndpoint 					  = "192.168.43.146:7051"
	authenticationCenterChannelName   = "authentication"
	authenticationCenterChaincodeName = "authentication"
	evidenceStorageChannelName        = "evidence"
	evidenceStorageChaincodeName      = "evidence"
)



var gw *client.Gateway
var authenticationCenterContract *client.Contract
var evidenceStorageContract *client.Contract

var now = time.Now()
var assetId = fmt.Sprintf("asset%d", now.Unix()*1e3+int64(now.Nanosecond())/1e6)

func Init() {
	clientConnection := newGrpcConnection()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	network := gw.GetNetwork(authenticationCenterChannelName)
	authenticationCenterContract = network.GetContract(authenticationCenterChaincodeName)

	network = gw.GetNetwork(evidenceStorageChannelName)
	evidenceStorageContract = network.GetContract(evidenceStorageChaincodeName)
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection() *grpc.ClientConn {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity() *identity.X509Identity {
	certificate, err := loadCertificate(certPath)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign() identity.Sign {
	files, err := os.ReadDir(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}
	privateKeyPEM, err := os.ReadFile(path.Join(keyPath, files[0].Name()))

	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

// This type of transaction would typically only be run once by an application the first time it was started after its
// initial deployment. A new version of the chaincode deployed later would likely not need to run an "init" function.
func initLedger(contract *client.Contract) {
	fmt.Printf("\n--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger \n")

	_, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Evaluate a transaction to query ledger state.
func getAllAssets(contract *client.Contract) {
	fmt.Println("\n--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

// Submit a transaction synchronously, blocking until it has been committed to the ledger.
func createAsset(contract *client.Contract) {
	fmt.Printf("\n--> Submit Transaction: CreateAsset, creates new asset with ID, Color, Size, Owner and AppraisedValue arguments \n")

	_, err := contract.SubmitTransaction("CreateAsset", assetId, "yellow", "5", "Tom", "1300")
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Evaluate a transaction by assetID to query ledger state.
func readAssetByID(contract *client.Contract) {
	fmt.Printf("\n--> Evaluate Transaction: ReadAsset, function returns asset attributes\n")

	evaluateResult, err := contract.EvaluateTransaction("ReadAsset", assetId)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

// Submit transaction asynchronously, blocking until the transaction has been sent to the orderer, and allowing
// this thread to process the chaincode response (e.g. update a UI) without waiting for the commit notification
func transferAssetAsync(contract *client.Contract) {
	fmt.Printf("\n--> Async Submit Transaction: TransferAsset, updates existing asset owner")

	submitResult, commit, err := contract.SubmitAsync("TransferAsset", client.WithArguments(assetId, "Mark"))
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction asynchronously: %w", err))
	}

	fmt.Printf("\n*** Successfully submitted transaction to transfer ownership from %s to Mark. \n", string(submitResult))
	fmt.Println("*** Waiting for transaction commit.")

	if commitStatus, err := commit.Status(); err != nil {
		panic(fmt.Errorf("failed to get commit status: %w", err))
	} else if !commitStatus.Successful {
		panic(fmt.Errorf("transaction %s failed to commit with status: %d", commitStatus.TransactionID, int32(commitStatus.Code)))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Submit transaction, passing in the wrong number of arguments ,expected to throw an error containing details of any error responses from the smart contract.
func exampleErrorHandling(contract *client.Contract) {
	fmt.Println("\n--> Submit Transaction: UpdateAsset asset70, asset70 does not exist and should return an error")

	_, err := contract.SubmitTransaction("UpdateAsset", "asset70", "blue", "5", "Tomoko", "300")
	if err == nil {
		panic("******** FAILED to return an error")
	}

	fmt.Println("*** Successfully caught the error:")

	switch err := err.(type) {
	case *client.EndorseError:
		fmt.Printf("Endorse error for transaction %s with gRPC status %v: %s\n", err.TransactionID, status.Code(err), err)
	case *client.SubmitError:
		fmt.Printf("Submit error for transaction %s with gRPC status %v: %s\n", err.TransactionID, status.Code(err), err)
	case *client.CommitStatusError:
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("Timeout waiting for transaction %s commit status: %s", err.TransactionID, err)
		} else {
			fmt.Printf("Error obtaining commit status for transaction %s with gRPC status %v: %s\n", err.TransactionID, status.Code(err), err)
		}
	case *client.CommitError:
		fmt.Printf("Transaction %s failed to commit with status %d: %s\n", err.TransactionID, int32(err.Code), err)
	default:
		panic(fmt.Errorf("unexpected error type %T: %w", err, err))
	}

	// Any error that originates from a peer or orderer node external to the gateway will have its details
	// embedded within the gRPC status error. The following code shows how to extract that.
	statusErr := status.Convert(err)

	details := statusErr.Details()
	if len(details) > 0 {
		fmt.Println("Error Details:")

		for _, detail := range details {
			switch detail := detail.(type) {
			case *gateway.ErrorDetail:
				fmt.Printf("- address: %s, mspId: %s, message: %s\n", detail.Address, detail.MspId, detail.Message)
			}
		}
	}
}

// Format JSON data
func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}
