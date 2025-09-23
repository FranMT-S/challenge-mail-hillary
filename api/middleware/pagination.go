package middleware

import (
	"context"
	"net/http"
	"strconv"

	"api/models"
)

type contextKey string

const (
	paginationContextKey = contextKey("pagination")
)

// Pagination is a middleware that parses pagination query parameters
func Pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Default values
		page := 1
		limit := DefaultPaginationLimit

		// Parse query parameters
		query := r.URL.Query()
		if pageParam := query.Get("page"); pageParam != "" {
			if _page, err := strconv.Atoi(pageParam); err == nil && _page > 0 {
				page = _page
			}
		}

		if limitParam := query.Get("limit"); limitParam != "" {
			if _limit, err := strconv.Atoi(limitParam); err == nil && _limit > 0 && _limit <= MaxPaginationLimit {
				limit = _limit
			}
		}

		// Create pagination model
		pagination := models.Pagination{
			Page:  page,
			Limit: limit,
		}

		// Add pagination to context
		ctx := context.WithValue(r.Context(), paginationContextKey, pagination)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetPaginationFromContext retrieves the pagination from context
func GetPaginationFromContext(ctx context.Context) models.Pagination {
	if p, ok := ctx.Value(paginationContextKey).(models.Pagination); ok {
		return p
	}

	return models.Pagination{Page: 1, Limit: DefaultPaginationLimit} // Default values
}
