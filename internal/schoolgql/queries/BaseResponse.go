package queries

import "strings"

type BaseResponse struct {
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type BasePlusResponse struct {
	BaseResponse
	Data any `json:"data"`
}

func (br *BaseResponse) GetErrorText() *string {
	if br == nil || br.Errors == nil {
		return nil
	}

	builder := strings.Builder{}

	for i := range br.Errors {
		if builder.Len() != 0 {
			builder.WriteString("; ")
		}

		builder.WriteString(br.Errors[i].Message)
	}

	res := builder.String()

	return &res
}
