package fabric

import (
	"encoding/json"
	"fabric-smart-evidence-storage/model"
	"fmt"
)

// AddEvidence adds a new evidence record with zero-knowledge proof
func AddEvidence(username, publicWitnessBytes, proofBytes string, caseNumber, caseInfo, manager1, manager2, evidenceType,
	ipfsLink, fileHash, fileName string) error {
	_, err := evidenceStorageContract.SubmitTransaction("AddEvidence",
		username, publicWitnessBytes, proofBytes,
		caseNumber, caseInfo, manager1, manager2,
		evidenceType, ipfsLink, fileHash, fileName)
	if err != nil {

		return fmt.Errorf("failed to add evidence: %w", err)
	}
	return nil
}

// GetEvidence retrieves evidence by username and case number
func GetEvidence(username, caseNumber string) (*model.Evidence, error) {
	evidenceBytes, err := evidenceStorageContract.EvaluateTransaction("GetEvidence", username, caseNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get evidence: %w", err)
	}

	var evidence model.Evidence
	err = json.Unmarshal(evidenceBytes, &evidence)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal evidence: %w", err)
	}
	return &evidence, nil
}

// GetUserEvidence retrieves all evidence records for a user
func GetUserEvidence(username string) ([]*model.Evidence, error) {
	evidencesBytes, err := evidenceStorageContract.EvaluateTransaction("GetUserEvidence", username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user evidence: %w", err)
	}
	var evidences []*model.Evidence
	if evidencesBytes == nil {
		return evidences, nil
	}
	err = json.Unmarshal(evidencesBytes, &evidences)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal evidences: %w", err)
	}
	return evidences, nil
}

// AddAuthorized adds a new authorization record
func AddAuthorized(username, caseNumber, authorizedUser string) error {
	_, err := evidenceStorageContract.SubmitTransaction("AddAuthorized", username, caseNumber, authorizedUser)
	if err != nil {
		return fmt.Errorf("failed to add authorized: %w", err)
	}
	return nil
}

// GetUserAllAuthorized retrieves all authorization records created by a user
func GetUserAllAuthorized(username string) ([]*model.Authorized, error) {
	authorizedBytes, err := evidenceStorageContract.EvaluateTransaction("GetUserAllAuthorized", username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user authorized: %w", err)
	}

	var authorized []*model.Authorized
	if authorizedBytes == nil {
		return authorized, nil
	}
	err = json.Unmarshal(authorizedBytes, &authorized)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal authorized: %w", err)
	}
	return authorized, nil
}

// CancelAuthorized cancels an authorization
func CancelAuthorized(username, txId string) error {
	_, err := evidenceStorageContract.SubmitTransaction("CancelAuthorized", username, txId)
	if err != nil {
		return fmt.Errorf("failed to cancel authorized: %w", err)
	}
	return nil
}

// GetUserReceivedAuthorized retrieves all authorization records received by a user
func GetUserReceivedAuthorized(username string) ([]*model.Authorized, error) {
	authorizedBytes, err := evidenceStorageContract.EvaluateTransaction("GetUserReceivedAuthorized", username)
	if err != nil {
		return nil, fmt.Errorf("failed to get received authorized: %w", err)
	}

	var authorized []*model.Authorized
	if authorizedBytes == nil {
		return authorized, nil
	}
	err = json.Unmarshal(authorizedBytes, &authorized)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal authorized: %w", err)
	}
	return authorized, nil
}

// View records a view action for an evidence
func View(username, caseNumber string) (*model.Evidence, error) {
	evidenceBytes, err := evidenceStorageContract.SubmitTransaction("View", username, caseNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to record view: %w", err)
	}
	var evidence model.Evidence
	err = json.Unmarshal(evidenceBytes, &evidence)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal evidence: %w", err)
	}
	return &evidence, nil
}

// GetUserViewRecord retrieves all view records for evidence created by a user
func GetUserViewRecord(username string) ([]*model.ViewRecord, error) {
	viewRecordsBytes, err := evidenceStorageContract.EvaluateTransaction("GetUserViewRecord", username)
	if err != nil {
		return nil, fmt.Errorf("failed to get view records: %w", err)
	}

	var viewRecords []*model.ViewRecord
	if viewRecordsBytes == nil {
		return viewRecords, nil
	}
	err = json.Unmarshal(viewRecordsBytes, &viewRecords)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal view records: %w", err)
	}
	return viewRecords, nil
}
