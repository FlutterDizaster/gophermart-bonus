package models

//go:generate easyjson -all credentials.go
type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
