package pgsql

import (
	"fmt"
	"time"
)

const (
	StatusDBConnectionException = 5
	StatusDBMigrationException  = 6
)

func ErrorStatusText(code int) string {
	switch code {
	case StatusDBConnectionException:
		return "Database Connection Exception"
	case StatusDBMigrationException:
		return "Database Migration Exception"
	default:
		return ""
	}
}

type Error struct {
	Time    time.Time
	Message string
	Err     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("error occured while database execution. error: %s, time: %v", e.Message, e.Time)
}
