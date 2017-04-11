package sdk

import (
	"fmt"
)

type FPTAIError struct {
	Code    int    `json:"error"`
	Message string `json:"message"`
}

func (e FPTAIError) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Message)
}
