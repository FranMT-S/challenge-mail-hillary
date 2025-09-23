package models

import "strings"

type OrderBy string

const (
	OrderByAsc  OrderBy = "asc"
	OrderByDesc OrderBy = "desc"
)

// NewOrderBy creates a new OrderBy instance
// If the orderBy is not valid, it returns OrderByAsc
func NewOrderBy(orderBy string) OrderBy {
	orderBy = strings.ToLower(orderBy)
	if orderBy != "asc" && orderBy != "desc" {
		return OrderByAsc
	}

	return OrderBy(orderBy)
}

func (ob OrderBy) Validate() bool {
	return ob == OrderByAsc || ob == OrderByDesc
}
