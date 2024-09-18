package handlers

import "fmt"

type RestError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *RestError) Error() string {
	return fmt.Sprintf("error occured. error: %s", e.Message)
}
