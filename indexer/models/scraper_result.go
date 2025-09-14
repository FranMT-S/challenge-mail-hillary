package models

type EmailResult struct {
	Email *Email `json:"email"`
	Error error  `json:"error"`
}

type PageResultStateType string

const (
	PageResultStatePending    PageResultStateType = "pending"
	PageResultStateProcessing PageResultStateType = "processing"
	PageResultStateFinished   PageResultStateType = "finished"
)

type PageResult struct {
	Page  int                 `json:"page"`
	Error string              `json:"error"`
	Total int                 `json:"total"`
	State PageResultStateType `json:"state"` // PageResultStatePending, PageResultStateProcessing, PageResultStateFinished
}
