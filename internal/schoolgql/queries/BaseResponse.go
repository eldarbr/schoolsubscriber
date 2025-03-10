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

	sb := strings.Builder{}

	for i := range br.Errors {
		if sb.Len() != 0 {
			sb.WriteString("; ")
		}

		sb.WriteString(br.Errors[i].Message)
	}

	res := sb.String()

	return &res
}
