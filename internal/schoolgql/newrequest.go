package schoolgql

import (
	"fmt"

	"github.com/eldarbr/schoolsubscriber/internal/schoolgql/queries"
)

type Request struct {
	OperationName queries.TOperationName `json:"operationName"`
	Query         queries.TQuery         `json:"query"`
	Variables     any                    `json:"variables"`
}

func NewRequest(operationName queries.TOperationName) (*Request, error) {
	query, err := queries.GetQueryByOperationName(operationName)
	if err != nil {
		return nil, fmt.Errorf("get query by operation name: %w", err)
	}

	return &Request{
		OperationName: operationName,
		Query:         *query,
		Variables:     struct{}{},
	}, nil
}
