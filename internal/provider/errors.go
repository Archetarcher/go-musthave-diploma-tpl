package provider

import (
	"fmt"
	"time"
)

type Error struct {
	Time    time.Time
	Message string
	Err     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("error occured in provider. error: %s, time: %v", e.Message, e.Time)
}
