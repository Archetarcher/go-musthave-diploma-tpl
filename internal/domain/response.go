package domain

type SuccessResponse struct {
	Code         int           `json:"code"`
	Message      string        `json:"message"`
	OrderAccrual *OrderAccrual `json:"-"`
}

type AuthResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

type OrderAccrualResponse struct {
	ID         int     `json:"-" `
	Number     uint64  `json:"number" `
	Status     string  `json:"status" `
	Accrual    float64 `json:"accrual" `
	UploadedAt string  `json:"uploaded_at"`
}

type OrderWithdrawalResponse struct {
	ID          int     `json:"-" `
	Order       uint64  `json:"order" `
	Sum         float64 `json:"sum" `
	ProcessedAt string  `json:"processed_at"`
}

type UserBalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type AccrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
