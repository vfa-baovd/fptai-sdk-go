package sdk

import (
	"fmt"
)

type Error struct {
	Code int `json:"error"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Message)
}