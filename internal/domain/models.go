package domain

import (
	"time"
)

type User struct {
	ID      int     `json:"id" db:"id"`
	Login   string  `json:"login" db:"login"`
	Hash    string  `json:"-" db:"hash"`
	Balance float64 `json:"-" db:"balance"`
}

type OrderAccrual struct {
	ID                  int        `json:"-" db:"id"`
	OrderId             uint64     `json:"number" db:"order_id" `
	UserId              int64      `json:"-" db:"user_id" `
	Status              string     `json:"status" db:"status"`
	Amount              *float64   `json:"accrual" db:"amount"`
	UploadedAt          ParsedTime `json:"uploaded_at" db:"uploaded_at"`
	ProcessingStartedAt *string    `json:"-" db:"processing_started_at"`
	ProcessedAt         *string    `json:"-" db:"processed_at"`
	InvalidatedAt       *string    `json:"-" db:"invalidated_at"`
}

type OrderWithdrawal struct {
	ID        int        `json:"-" db:"id"`
	OrderId   uint64     `json:"order_id" db:"order_id" `
	UserId    int64      `json:"user_id" db:"user_id" `
	Amount    *float64   `json:"amount" db:"amount"`
	CreatedAt ParsedTime `json:"created_at" db:"created_at"`
}

type ParsedTime time.Time

func (t ParsedTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(t).Format(time.RFC3339) + `"`), nil
}

const (
	OrderTypeAccrual    = "accrual"
	OrderTypeWithdrawal = "withdrawal"
)
const (
	OrderStatusNew        = "NEW"
	OrderStatusRegistered = "REGISTERED"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)
