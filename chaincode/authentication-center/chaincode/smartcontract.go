package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

const (
	adminKeyPrefix = "admin"
	userKeyPrefix  = "user"
)

const (
	enableStatus  = "enable"
	disableStatus = "disable"
)

const (
	AdminRole = "1"
	InputRole = "2"
	QueryRole = "3"
)

type User struct {
	Username     string `json:"username"`
	Pwd          string `json:"pwd"`
	Nickname     string `json:"nickname"`
	Status       string `json:"status"`
	Role         string `json:"role"`
	RegisterTime string `json:"registerTime"`
}

func (s *SmartContract) Ping(ctx contractapi.TransactionContextInterface) (string, error) {
	return "pong", nil
}

func (s *SmartContract) InitAdmin(ctx contractapi.TransactionContextInterface, username, pwd, nickname string) error {
	// 检测admin是否存在
	adminIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(adminKeyPrefix, []string{})
	if err != nil {
		return err
	}
	defer adminIterator.Close()
	if adminIterator.HasNext() {
		return fmt.Errorf("admin already exists")
	}
	adminKey, err := ctx.GetStub().CreateCompositeKey(adminKeyPrefix, []string{username})
	if err != nil {
		return err
	}
	txTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return err
	}
	registerTime := txTimestamp.AsTime().Format("2006-01-02 15:04:05")
	admin := User{
		Username:     username,
		Pwd:          pwd,
		Nickname:     nickname,
		Status:       enableStatus,
		Role:         AdminRole,
		RegisterTime: registerTime,
	}
	adminBytes, err := json.Marshal(admin)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(adminKey, adminBytes)
	if err != nil {
		return err
	}
	userKey, err := ctx.GetStub().CreateCompositeKey(userKeyPrefix, []string{username})
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(userKey, adminBytes)
	if err != nil {
		return err
	}
	return nil
}

func (s *SmartContract) Register(ctx contractapi.TransactionContextInterface, username, pwd, nickname, role string) error {
	if role != InputRole && role != QueryRole {
		return fmt.Errorf("invalid role")
	}
	userKey, err := ctx.GetStub().CreateCompositeKey(userKeyPrefix, []string{username})
	if err != nil {
		return err
	}
	if userBytes, err := ctx.GetStub().GetState(userKey); err != nil {
		return err
	} else if userBytes != nil {
		return fmt.Errorf("user %s already exists", username)
	}
	txTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return err
	}
	registerTime := txTimestamp.AsTime().Format("2006-01-02 15:04:05")
	user := User{
		Username:     username,
		Pwd:          pwd,
		Nickname:     nickname,
		Status:       disableStatus,
		Role:         role,
		RegisterTime: registerTime,
	}
	userBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(userKey, userBytes)
}

func (s *SmartContract) GetUser(ctx contractapi.TransactionContextInterface, username string) (*User, error) {
	userKey, err := ctx.GetStub().CreateCompositeKey(userKeyPrefix, []string{username})
	if err != nil {
		return nil, err
	}
	userBytes, err := ctx.GetStub().GetState(userKey)
	if err != nil {
		return nil, err
	}
	if userBytes == nil {
		return nil, fmt.Errorf("user %s not found", username)
	}
	var user User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *SmartContract) GetAllUsers(ctx contractapi.TransactionContextInterface) ([]*User, error) {
	// Query for non-admin users
	queryString := fmt.Sprintf("{\"selector\":{\"role\":{\"$ne\":\"%s\"}}}", AdminRole)
	userIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer userIterator.Close()
	var users []*User
	for userIterator.HasNext() {
		userResponse, err := userIterator.Next()
		if err != nil {
			return nil, err
		}
		var user User
		err = json.Unmarshal(userResponse.Value, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (s *SmartContract) UpdatePwd(ctx contractapi.TransactionContextInterface, username, newPwd string) error {
	userKey, err := ctx.GetStub().CreateCompositeKey(userKeyPrefix, []string{username})
	if err != nil {
		return err
	}
	userBytes, err := ctx.GetStub().GetState(userKey)
	if err != nil {
		return err
	}
	if userBytes == nil {
		return fmt.Errorf("user %s not found", username)
	}
	var user User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return err
	}
	user.Pwd = newPwd
	userBytes, err = json.Marshal(user)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(userKey, userBytes)
}

func (s *SmartContract) UpdateRoleAndStatus(ctx contractapi.TransactionContextInterface, username, role, status string) error {
	if status != enableStatus && status != disableStatus {
		return fmt.Errorf("invalid status")
	}
	userKey, err := ctx.GetStub().CreateCompositeKey(userKeyPrefix, []string{username})
	if err != nil {
		return err
	}
	userBytes, err := ctx.GetStub().GetState(userKey)
	if err != nil {
		return err
	}
	if userBytes == nil {
		return fmt.Errorf("user %s not found", username)
	}
	var user User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return err
	}
	user.Role = role
	user.Status = status
	userBytes, err = json.Marshal(user)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(userKey, userBytes)
}

func (s *SmartContract) DeleteUser(ctx contractapi.TransactionContextInterface, username string) error {
	userKey, err := ctx.GetStub().CreateCompositeKey(userKeyPrefix, []string{username})
	if err != nil {
		return err
	}
	userBytes, err := ctx.GetStub().GetState(userKey)
	if err != nil {
		return err
	}
	if userBytes == nil {
		return fmt.Errorf("user %s not found", username)
	}
	return ctx.GetStub().DelState(userKey)
}
