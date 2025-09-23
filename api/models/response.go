package models

type Status string

const (
	StatusSuccess Status = "success"
	StatusError   Status = "error"
	StatusNoData  Status = "no data"
)

// Response is the response type for all operations
// msg  status of the operation
// data  data of the operation
// error error of the operation
type Response[T any] struct {
	Msg   Status `json:"msg"`
	Data  T      `json:"data"`
	Error string `json:"error,omitempty"`
}

func NewResponse[T any](msg Status, data T, error string) Response[T] {
	return Response[T]{
		Msg:   msg,
		Data:  data,
		Error: error,
	}
}

// MailResponse is the response type for mail operations
// mails  array of emails
// total s total number of emails in the database
type MailResponse struct {
	Mails []Email `json:"mails"`
	Total int64   `json:"total"`
}
