package models

type Operator string

// Operator constants to search by date
const (
	OperatorEqual              Operator = "="
	OperatorLessThan           Operator = "<"
	OperatorLessThanOrEqual    Operator = "<="
	OperatorGreaterThan        Operator = ">"
	OperatorGreaterThanOrEqual Operator = ">="
)

// Validate checks if the operator is valid
func (o Operator) Validate() bool {
	return o == OperatorEqual || o == OperatorLessThan || o == OperatorLessThanOrEqual || o == OperatorGreaterThan || o == OperatorGreaterThanOrEqual
}
