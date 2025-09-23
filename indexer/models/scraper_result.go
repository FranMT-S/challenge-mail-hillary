package models

type EmailResult struct {
	Email *Email `json:"email"`
	Error error  `json:"error"`
}

// PageResultStateType represents the state of a page
type PageResultStateType string

const (
	PageResultStatePending    PageResultStateType = "pending"
	PageResultStateProcessing PageResultStateType = "processing"
	PageResultStateFinished   PageResultStateType = "finished"
)

// PageResult represents the result of a page
// Page: Page number
// Error: Error message
// Total: Total number of emails
// State: State of the page
type PageResult struct {
	Page  int                 `json:"page"`
	Error string              `json:"error"`
	Total int                 `json:"total"`
	State PageResultStateType `json:"state"` // PageResultStatePending, PageResultStateProcessing, PageResultStateFinished
}
