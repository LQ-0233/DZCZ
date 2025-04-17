package fabric

import (
	"encoding/json"
	"fmt"
	"fabric-smart-evidence-storage/model"
)



// Register creates a new user account
func Register(username, pwd, nickname, role string) error {
	_, err := authenticationCenterContract.SubmitTransaction("Register", username, pwd, nickname, role)
	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}
	return nil
}

// GetUser retrieves user information by username
func GetUser(username string) (*model.User, error) {
	userBytes, err := authenticationCenterContract.EvaluateTransaction("GetUser", username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	var user model.User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}
	return &user, nil
}

// GetAllUsers retrieves all users from the system
func GetAllUsers() ([]*model.User, error) {
	usersBytes, err := authenticationCenterContract.EvaluateTransaction("GetAllUsers")
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	var users []*model.User
	if usersBytes == nil {
		return users, nil
	}
	err = json.Unmarshal(usersBytes, &users)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal users: %w", err)
	}
	return users, nil
}


// UpdatePwd updates user's password
func UpdatePwd(username, newPwd string) error {
	_, err := authenticationCenterContract.SubmitTransaction("UpdatePwd", username, newPwd)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

// UpdateRoleAndStatus updates user's role and status
func UpdateRoleAndStatus(username, role, status string) error {
	_, err := authenticationCenterContract.SubmitTransaction("UpdateRoleAndStatus", username, role, status)
	if err != nil {
		return fmt.Errorf("failed to update role and status: %w", err)
	}
	return nil
}

