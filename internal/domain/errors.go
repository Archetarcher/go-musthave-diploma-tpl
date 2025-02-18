package domain

import "fmt"

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("error occured. error: %s", e.Message)
}
