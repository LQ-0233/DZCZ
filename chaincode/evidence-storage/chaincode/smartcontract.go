package chaincode

import (
	"bytes"
	_ "crypto/elliptic"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//go:embed vk
var vkBytes []byte

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
	vk groth16.VerifyingKey
}

const (
	evidenceKeyPrefix   = "evidence"
	authorizedKeyPrefix = "authorized"
	viewRecordKeyPrefix = "viewRecord"
)

const (
	EvidenceType   = "evidence"
	AuthorizedType = "authorized"
	ViewRecordType = "viewRecord"
)

type Evidence struct {
	CaseNumber      string `json:"caseNumber"`
	CaseInfo        string `json:"caseInfo"`
	Manager1        string `json:"manager1"`
	Manager2        string `json:"manager2"`
	EvidenceType    string `json:"evidenceType"`
	IpfsLink        string `json:"ipfsLink"`
	FileHash        string `json:"fileHash"`
	FileName        string `json:"fileName"`
	EvidenceCreator string `json:"evidenceCreator"`
	CreateTime      string `json:"createTime"`
	ViewCount       int    `json:"viewCount"`
	Type            string `json:"type"`
}

type Authorized struct {
	Id              string `json:"id"`
	CaseNumber      string `json:"caseNumber"`
	EvidenceCreator string `json:"evidenceCreator"`
	AuthTime        string `json:"authTime"`
	AuthorizedUser  string `json:"authorizedUser"`
	Status          string `json:"status"`
	ViewCount       int    `json:"viewCount"`
	Type            string `json:"type"`
}

const (
	AUTHORIZED = "authorized"
	CANCELED   = "canceled"
)

type ViewRecord struct {
	Id              string `json:"id"`
	CaseNumber      string `json:"caseNumber"`
	EvidenceCreator string `json:"evidenceCreator"`
	ViewTime        string `json:"viewTime"`
	ViewUser        string `json:"viewUser"`
	Type            string `json:"type"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	var buf bytes.Buffer
	vk := groth16.NewVerifyingKey(ecc.BN254)
	buf.Write(vkBytes)
	_, err := vk.ReadFrom(&buf)
	if err != nil {
		return fmt.Errorf("failed to read vk: %v", err)
	}

	s.vk = vk
	return nil
}

func (s *SmartContract) Ping(ctx contractapi.TransactionContextInterface) (string, error) {
	return "pong", nil
}

func (s *SmartContract) AddEvidence(ctx contractapi.TransactionContextInterface, username, publicWitnessBytes,
	proofBytes string, caseNumber, caseInfo, manager1, manager2, evidenceType, ipfsLink, fileHash, fileName string) error {
	err := s.Verify(ctx, publicWitnessBytes, proofBytes)
	if err != nil {
		return fmt.Errorf("proof verification failed: %v", err)
	}
	// 检测证据是否存在
	query := fmt.Sprintf(`{"selector":{"type":"%s", "caseNumber":"%s"}}`, EvidenceType, caseNumber)
	evidenceIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return fmt.Errorf("failed to get evidence iterator: %v", err)
	}
	defer evidenceIterator.Close()
	if evidenceIterator.HasNext() {
		return fmt.Errorf("evidence %s already exists", caseNumber)
	}

	// 创建证据
	txTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get tx timestamp: %v", err)
	}
	createTime := txTimestamp.AsTime().Format("2006-01-02 15:04:05")
	evidenceKey, err := ctx.GetStub().CreateCompositeKey(evidenceKeyPrefix, []string{username, caseNumber})
	if err != nil {
		return fmt.Errorf("failed to create evidence key: %v", err)
	}
	evidence := Evidence{
		CaseNumber:      caseNumber,
		CaseInfo:        caseInfo,
		Manager1:        manager1,
		Manager2:        manager2,
		EvidenceType:    evidenceType,
		IpfsLink:        ipfsLink,
		FileHash:        fileHash,
		FileName:        fileName,
		EvidenceCreator: username,
		CreateTime:      createTime,
		ViewCount:       0,
		Type:            EvidenceType,
	}
	evidenceBytes, err := json.Marshal(evidence)
	if err != nil {
		return fmt.Errorf("failed to marshal evidence: %v", err)
	}
	return ctx.GetStub().PutState(evidenceKey, evidenceBytes)
}

func (s *SmartContract) GetEvidence(ctx contractapi.TransactionContextInterface, username, caseNumber string) (*Evidence, error) {
	evidenceKey, err := ctx.GetStub().CreateCompositeKey(evidenceKeyPrefix, []string{username, caseNumber})
	if err != nil {
		return nil, fmt.Errorf("failed to create evidence key: %v", err)
	}
	evidenceBytes, err := ctx.GetStub().GetState(evidenceKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get evidence: %v", err)
	}
	if evidenceBytes == nil {
		return nil, fmt.Errorf("evidence %s not found", caseNumber)
	}
	var evidence Evidence
	err = json.Unmarshal(evidenceBytes, &evidence)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal evidence: %v", err)
	}
	return &evidence, nil
}

func (s *SmartContract) GetUserEvidence(ctx contractapi.TransactionContextInterface, username string) ([]*Evidence, error) {
	evidenceIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(evidenceKeyPrefix, []string{username})
	if err != nil {
		return nil, fmt.Errorf("failed to get evidence iterator: %v", err)
	}
	defer evidenceIterator.Close()
	var evidences []*Evidence
	for evidenceIterator.HasNext() {
		evidenceResponse, err := evidenceIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get evidence response: %v", err)
		}
		var evidence Evidence
		err = json.Unmarshal(evidenceResponse.Value, &evidence)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal evidence: %v", err)
		}
		evidences = append(evidences, &evidence)
	}
	return evidences, nil
}

func (s *SmartContract) AddAuthorized(ctx contractapi.TransactionContextInterface, username, caseNumber, authorizedUser string) error {
	// 检测证据是否存在
	evidenceKey, err := ctx.GetStub().CreateCompositeKey(evidenceKeyPrefix, []string{username, caseNumber})
	if err != nil {
		return fmt.Errorf("failed to create evidence key: %v", err)
	}
	evidenceBytes, err := ctx.GetStub().GetState(evidenceKey)
	if err != nil {
		return fmt.Errorf("failed to get evidence: %v", err)
	}
	if evidenceBytes == nil {
		return fmt.Errorf("evidence %s not found", caseNumber)
	}
	// 检测是否已授权
	query := fmt.Sprintf(`{"selector":{"type":"%s", "caseNumber":"%s", "authorizedUser":"%s", "status":"%s"}}`, AuthorizedType, caseNumber, authorizedUser, AUTHORIZED)
	authorizedUserIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return fmt.Errorf("failed to get authorized user iterator: %v", err)
	}
	defer authorizedUserIterator.Close()
	if authorizedUserIterator.HasNext() {
		return fmt.Errorf("authorized user %s already exists", authorizedUser)
	}
	// 创建授权记录
	txId := ctx.GetStub().GetTxID()
	authorizedKey, err := ctx.GetStub().CreateCompositeKey(authorizedKeyPrefix, []string{username, caseNumber, txId})
	if err != nil {
		return fmt.Errorf("failed to create authorized key: %v", err)
	}
	txTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get tx timestamp: %v", err)
	}
	authTime := txTimestamp.AsTime().Format("2006-01-02 15:04:05")
	authorized := Authorized{
		Id:              txId,
		CaseNumber:      caseNumber,
		EvidenceCreator: username,
		AuthTime:        authTime,
		AuthorizedUser:  authorizedUser,
		Status:          AUTHORIZED,
		ViewCount:       0,
		Type:            AuthorizedType,
	}
	authorizedBytes, err := json.Marshal(authorized)
	if err != nil {
		return fmt.Errorf("failed to marshal authorized: %v", err)
	}
	err = ctx.GetStub().PutState(authorizedKey, authorizedBytes)
	if err != nil {
		return fmt.Errorf("failed to put authorized: %v", err)
	}
	return nil
}

func (s *SmartContract) GetUserAllAuthorized(ctx contractapi.TransactionContextInterface, username string) ([]*Authorized, error) {
	authorizedIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(authorizedKeyPrefix, []string{username})
	if err != nil {
		return nil, fmt.Errorf("failed to get authorized iterator: %v", err)
	}
	defer authorizedIterator.Close()
	var authorizeds []*Authorized
	for authorizedIterator.HasNext() {
		authorizedResponse, err := authorizedIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get authorized response: %v", err)
		}
		var authorized Authorized
		err = json.Unmarshal(authorizedResponse.Value, &authorized)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal authorized: %v", err)
		}
		authorizeds = append(authorizeds, &authorized)
	}
	return authorizeds, nil
}

func (s *SmartContract) CancelAuthorized(ctx contractapi.TransactionContextInterface, username, txId string) error {
	query := fmt.Sprintf(`{"selector":{"type":"%s", "id":"%s", "status":"%s", "evidenceCreator":"%s"}}`, AuthorizedType, txId, AUTHORIZED, username)
	authorizedIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return fmt.Errorf("failed to get authorized iterator: %v", err)
	}
	defer authorizedIterator.Close()
	if !authorizedIterator.HasNext() {
		return fmt.Errorf("authorized %s not found", txId)
	}
	authorizedResponse, err := authorizedIterator.Next()
	if err != nil {
		return fmt.Errorf("failed to get authorized response: %v", err)
	}
	var authorized Authorized
	err = json.Unmarshal(authorizedResponse.Value, &authorized)
	if err != nil {
		return fmt.Errorf("failed to unmarshal authorized: %v", err)
	}
	authorized.Status = CANCELED
	authorizedBytes, err := json.Marshal(authorized)
	if err != nil {
		return fmt.Errorf("failed to marshal authorized: %v", err)
	}
	authorizedKey, err := ctx.GetStub().CreateCompositeKey(authorizedKeyPrefix, []string{authorized.EvidenceCreator, authorized.CaseNumber, authorized.Id})
	if err != nil {
		return fmt.Errorf("failed to create authorized key: %v", err)
	}
	err = ctx.GetStub().PutState(authorizedKey, authorizedBytes)
	if err != nil {
		return fmt.Errorf("failed to put authorized: %v", err)
	}
	return nil
}

func (s *SmartContract) GetUserReceivedAuthorized(ctx contractapi.TransactionContextInterface, username string) ([]*Authorized, error) {
	query := fmt.Sprintf(`{"selector":{"type":"%s", "authorizedUser":"%s", "status":"%s"}}`, AuthorizedType, username, AUTHORIZED)
	authorizedIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get authorized iterator: %v", err)
	}
	defer authorizedIterator.Close()
	var authorizeds []*Authorized
	for authorizedIterator.HasNext() {
		authorizedResponse, err := authorizedIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get authorized response: %v", err)
		}
		var authorized Authorized
		err = json.Unmarshal(authorizedResponse.Value, &authorized)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal authorized: %v", err)
		}
		authorizeds = append(authorizeds, &authorized)
	}
	return authorizeds, nil
}

func (s *SmartContract) View(ctx contractapi.TransactionContextInterface, username, caseNumber string) (*Evidence, error) {
	// 检测是否已授权
	query := fmt.Sprintf(`{"selector":{"type":"%s", "caseNumber":"%s", "authorizedUser":"%s", "status":"%s"}}`, AuthorizedType, caseNumber, username, AUTHORIZED)
	authorizedIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get authorized iterator: %v", err)
	}
	defer authorizedIterator.Close()
	if !authorizedIterator.HasNext() {
		return nil, fmt.Errorf("authorized %s not found", caseNumber)
	}
	authorizedResponse, err := authorizedIterator.Next()
	if err != nil {
		return nil, fmt.Errorf("failed to get authorized response: %v", err)
	}
	var authorized Authorized
	err = json.Unmarshal(authorizedResponse.Value, &authorized)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal authorized: %v", err)
	}
	// 更新授权记录的查看次数
	authorized.ViewCount++
	authorizedBytes, err := json.Marshal(authorized)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal authorized: %v", err)
	}
	authorizedKey, err := ctx.GetStub().CreateCompositeKey(authorizedKeyPrefix, []string{authorized.EvidenceCreator, authorized.CaseNumber, authorized.Id})
	if err != nil {
		return nil, fmt.Errorf("failed to create authorized key: %v", err)
	}
	err = ctx.GetStub().PutState(authorizedKey, authorizedBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to put authorized: %v", err)
	}
	// 更新证据的查看次数
	evidenceKey, err := ctx.GetStub().CreateCompositeKey(evidenceKeyPrefix, []string{authorized.EvidenceCreator, authorized.CaseNumber})
	if err != nil {
		return nil, fmt.Errorf("failed to create evidence key: %v", err)
	}
	evidenceBytes, err := ctx.GetStub().GetState(evidenceKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get evidence: %v", err)
	}
	var evidence Evidence
	err = json.Unmarshal(evidenceBytes, &evidence)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal evidence: %v", err)
	}
	evidence.ViewCount++
	evidenceBytes, err = json.Marshal(evidence)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal evidence: %v", err)
	}
	err = ctx.GetStub().PutState(evidenceKey, evidenceBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to put evidence: %v", err)
	}
	txTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return nil, fmt.Errorf("failed to get tx timestamp: %v", err)
	}
	txId := ctx.GetStub().GetTxID()
	// 创建查看记录
	viewRecordKey, err := ctx.GetStub().CreateCompositeKey(viewRecordKeyPrefix, []string{authorized.EvidenceCreator, caseNumber, txId})
	if err != nil {
		return nil, fmt.Errorf("failed to create view record key: %v", err)
	}
	viewRecord := ViewRecord{
		Id:              txId,
		CaseNumber:      caseNumber,
		EvidenceCreator: authorized.EvidenceCreator,
		ViewTime:        txTimestamp.AsTime().Format("2006-01-02 15:04:05"),
		ViewUser:        username,
		Type:            ViewRecordType,
	}
	viewRecordBytes, err := json.Marshal(viewRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal view record: %v", err)
	}
	err = ctx.GetStub().PutState(viewRecordKey, viewRecordBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to put view record: %v", err)
	}
	return &evidence, nil
}

func (s *SmartContract) GetUserViewRecord(ctx contractapi.TransactionContextInterface, username string) ([]*ViewRecord, error) {
	query := fmt.Sprintf(`{"selector":{"type":"%s", "evidenceCreator":"%s"}}`, ViewRecordType, username)
	viewRecordIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get view record iterator: %v", err)
	}
	defer viewRecordIterator.Close()
	var viewRecords []*ViewRecord
	for viewRecordIterator.HasNext() {
		viewRecordResponse, err := viewRecordIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get view record response: %v", err)
		}
		var viewRecord ViewRecord
		err = json.Unmarshal(viewRecordResponse.Value, &viewRecord)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal view record: %v", err)
		}
		viewRecords = append(viewRecords, &viewRecord)
	}
	return viewRecords, nil
}

func (s *SmartContract) VerifyBase64(ctx contractapi.TransactionContextInterface, publicWitnessBase64, proofBase64 string) error {
	// 零知识证明 验证
	publicWitnessBytes, err := base64.StdEncoding.DecodeString(publicWitnessBase64)
	if err != nil {
		return err
	}
	proofBytes, err := base64.StdEncoding.DecodeString(proofBase64)
	if err != nil {
		return err
	}
	return Verify(s.vk, publicWitnessBytes, proofBytes)
}

func (s *SmartContract) Verify(ctx contractapi.TransactionContextInterface, publicWitnessBytes, proofBytes string) error {
	// 零知识证明 验证
	err := Verify(s.vk, []byte(publicWitnessBytes), []byte(proofBytes))
	if err != nil {
		return err
	}
	return nil
}

func Verify(vk groth16.VerifyingKey, publicWitnessBytes, proofBytes []byte) error {

	publicWitness, err := witness.New(ecc.BN254.ScalarField()) // note that schema is optional for binary encoding
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	_, err = buf.Write(publicWitnessBytes)
	if err != nil {
		return err
	}
	_, err = publicWitness.ReadFrom(&buf)
	if err != nil {
		return err
	}
	buf.Reset()
	_, err = buf.Write(proofBytes)
	if err != nil {
		return err
	}
	newProof := groth16.NewProof(ecc.BN254)
	_, err = newProof.ReadFrom(&buf)
	if err != nil {
		return err
	}
	err = groth16.Verify(newProof, vk, publicWitness)
	if err != nil {
		return err
	}
	return nil
}
