package models

// Pagination represents pagination parameters
type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}
