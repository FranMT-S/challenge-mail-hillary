package models

import (
	"api/config"
)

type TypeSearch string

// TypeSearch constants to search by type
const (
	TypeSearchAND TypeSearch = "AND"
	TypeSearchOR  TypeSearch = "OR"
)

// Validate checks if the type search is valid
func (ts TypeSearch) Validate() bool {
	return ts == TypeSearchAND || ts == TypeSearchOR
}

// DateSearch represents a date search with operator
type DateSearch struct {
	Date     *DateParam `json:"date,omitempty"`
	Operator Operator   `json:"operator"`
}

// QuerySearch represents a query search with pagination and date search and order
type QuerySearch struct {
	Query      string     `json:"query"` // Query to search
	TypeSearch TypeSearch `json:"type"`  // AND or OR
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
	Date       *DateParam `json:"date"`
	OrderBy    OrderBy    `json:"orderBy"` // ASC or DESC
	DateSearch DateSearch `json:"dateSearch"`
}

func NewQuerySearch(query string, typeSearch TypeSearch, orderBy OrderBy, page int, limit int, dateSearch DateSearch) *QuerySearch {
	if !typeSearch.Validate() {
		typeSearch = TypeSearchAND
	}

	if !orderBy.Validate() {
		orderBy = OrderByDesc
	}

	return &QuerySearch{
		Query:      query,
		TypeSearch: typeSearch,
		Page:       page,
		Limit:      limit,
		DateSearch: dateSearch,
		OrderBy:    orderBy,
	}
}

func (qs *QuerySearch) Normalize() *QuerySearch {
	if qs.Query == "" {
		qs.Query = ""
	}
	if !qs.TypeSearch.Validate() {
		qs.TypeSearch = TypeSearchAND
	}

	if !qs.OrderBy.Validate() {
		qs.OrderBy = OrderByDesc
	}

	if qs.Page < 1 {
		qs.Page = 1
	}

	if qs.Limit < 1 {
		qs.Limit = 1
	}

	if qs.DateSearch.Date != nil {
		if qs.DateSearch.Date.Valid {
			if !qs.DateSearch.Operator.Validate() {
				qs.DateSearch.Operator = OperatorLessThanOrEqual
			}
		}
	}

	if qs.Limit > config.GetApiConfig().MaxLimitPagination {
		qs.Limit = config.GetApiConfig().MaxLimitPagination
	}

	return qs
}
