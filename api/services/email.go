package services

import (
	"api/models"
	"api/sanatizer"
	"context"
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GetEmailsResponse struct {
	Emails []models.Email `json:"emails"`
	Total  int64          `json:"total"`
}

// EmailService defines the interface for email-related operations
type EmailService interface {
	// GetEmails retrieves a paginated list of emails
	SearchEmails(ctx context.Context, query models.QuerySearch) (*GetEmailsResponse, error)
}

type emailService struct {
	db *gorm.DB
}

// NewEmailService creates a new instance of EmailService
func NewEmailService(db *gorm.DB) EmailService {
	return &emailService{
		db: db,
	}
}

// SearchEmails implements EmailService interface
// SearchEmails retrieves a paginated list of emails based on a search query using inverted index
// can filter and order by Date
func (s *emailService) SearchEmails(ctx context.Context, query models.QuerySearch) (*GetEmailsResponse, error) {

	query = *query.Normalize()

	if ctx == nil {
		ctx = context.Background()
	}

	var emails []models.Email
	var emailsRank []models.EmailRank
	var total int64

	query.Query = sanatizer.SanitizeForSQL(query.Query, 80)

	if len(query.Query) > 80 {
		return nil, models.NewApiError("query too long", nil)
	}

	queriesSplit := strings.Split(query.Query, " ")
	sanitizedQueries := make([]string, 0)
	for _, query := range queriesSplit {
		query = strings.TrimSpace(query)
		query = regexp.MustCompile(`[^a-zA-Z0-9]`).ReplaceAllString(query, "")
		if query == "" {
			continue
		}

		sanitizedQueries = append(sanitizedQueries, query)
	}

	var querySearch string

	if query.TypeSearch == models.TypeSearchAND {
		querySearch = strings.Join(sanitizedQueries, " & ")
	} else {
		querySearch = strings.Join(sanitizedQueries, " | ")
	}

	offset := (query.Page - 1) * query.Limit
	var tx *gorm.DB

	// create query to find mails
	if len(querySearch) < 2 {
		tx = s.createAllMailsSearchQuery(ctx, query)
	} else {
		tx = s.createQuerySearch(ctx, query, querySearch)
	}

	// count total
	tx.Count(&total)

	// limit and offset
	tx = tx.Limit(query.Limit).Offset(offset)

	err := tx.Scan(&emailsRank).Error
	if err != nil {
		return nil, models.NewApiError("cannot retrieve emails", err)
	}

	for _, email := range emailsRank {
		emails = append(emails, models.Email{
			ID:      email.ID,
			Date:    email.Date,
			Subject: email.Subject,
			From:    email.From,
			To:      email.To,
			Content: email.Content,
		})
	}

	return &GetEmailsResponse{Emails: emails, Total: total}, nil
}

func (s *emailService) createAllMailsSearchQuery(ctx context.Context, query models.QuerySearch) *gorm.DB {
	whereComparison := fmt.Sprintf("e.date::date %s ?", string(query.DateSearch.Operator))
	dateOrderBy := clause.OrderByColumn{Column: clause.Column{Name: "e.date"}, Desc: query.OrderBy == models.OrderByDesc}
	orderBy := clause.OrderByColumn{Column: clause.Column{Name: "e.id"}, Desc: query.OrderBy == models.OrderByDesc}

	tx := s.db.WithContext(ctx).
		Table("emails_hillary.emails e").
		Select(`e.id, e.subject, e."from", e."to", e.content, e.date`)

	if query.DateSearch.Date != nil && query.DateSearch.Date.Valid {
		tx = tx.Where(whereComparison, query.DateSearch.Date.Format("2006-01-02"))
	}

	if query.DateSearch.Date != nil && query.DateSearch.Date.Valid {
		orderBy = dateOrderBy
	}

	tx = tx.Order(orderBy)

	return tx
}

func (s *emailService) createQuerySearch(ctx context.Context, query models.QuerySearch, querySearch string) *gorm.DB {
	tx := s.db.WithContext(ctx).Table("emails_hillary.emails_search es").
		Joins("JOIN emails_hillary.emails e ON e.id = es.id").
		Select(`
			e.id, 
			e.subject, 
			e."from", 
			e."to", 
			e.content, 
			e.date,
			ts_rank(es.search_vector, to_tsquery('english', ?)) AS rank`,
			querySearch)

	// if date exist order by date first before order by rank
	if query.DateSearch.Date != nil && query.DateSearch.Date.Valid {
		whereComparison := fmt.Sprintf("e.date::date %s ?", string(query.DateSearch.Operator))
		orderBy := clause.OrderByColumn{Column: clause.Column{Name: "e.date"}, Desc: query.OrderBy == models.OrderByDesc}
		tx = tx.Where(whereComparison, query.DateSearch.Date.Format("2006-01-02")).
			Order(orderBy)
	}

	tx = tx.Order("rank DESC")
	tx = tx.Where("es.search_vector @@ to_tsquery('english', ?)", querySearch)

	return tx
}
