package domain

type AuthRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type OrderAccrualRequest struct {
	OrderId string `json:"order_id" validate:"required"`
}

type OrderWithdrawalRequest struct {
	OrderId string  `json:"order" validate:"required"`
	Sum     float64 `json:"sum" validate:"required"`
}
