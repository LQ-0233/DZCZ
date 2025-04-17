package model


type User struct {
	Username     string `json:"username"`
	Pwd          string `json:"pwd"`
	Nickname     string `json:"nickname"`
	Status       string `json:"status"`
	Role         string `json:"role"`
	RegisterTime string `json:"registerTime"`
}

const (
	enableStatus  = "enable"
	disableStatus = "disable"
)

const (
	AdminRole = "1"
	InputRole = "2"
	QueryRole = "3"
)