package models_test

import (
	"testing"

	"api/models"
)

func TestNewOrderBy(t *testing.T) {
	ttc := []struct {
		name     string
		orderBy  string
		expected models.OrderBy
	}{
		{"must be asc", "asc", models.OrderByAsc},
		{"must be desc", "desc", models.OrderByDesc},
		{"must be asc", "not valid", models.OrderByAsc},
	}

	for _, tt := range ttc {
		t.Run(tt.name, func(t *testing.T) {
			newOrderBy := models.NewOrderBy(tt.orderBy)
			if newOrderBy != tt.expected {
				t.Errorf("NewOrderBy returned %v, expected %v", newOrderBy, tt.expected)
			}
		})
	}

}
