package repositories

import (
	"fmt"
	"time"
)

const (
	StatusDBTransactionException = 0
)

func ErrorStatusText(code int) string {
	switch code {
	case StatusDBTransactionException:
		return "Database Transaction Exception"
	default:
		return ""
	}
}

type Error struct {
	Code    time.Time
	Time    time.Time
	Message string
	Err     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("error occured in repository. error: %s, time: %v", e.Message, e.Time)
}
