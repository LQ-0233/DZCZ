package model



type Evidence struct {
	CaseNumber      string `json:"caseNumber"`
	CaseInfo        string `json:"caseInfo"`
	Manager1        string `json:"manager1"`
	Manager2        string `json:"manager2"`
	EvidenceType    string `json:"evidenceType"`
	IpfsLink        string `json:"ipfsLink"`
	FileHash        string `json:"fileHash"`
	FileName        string `json:"fileName"`
	FilePath        string `json:"filePath"`
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