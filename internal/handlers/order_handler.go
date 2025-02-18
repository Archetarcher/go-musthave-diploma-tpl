package handlers

import (
	"context"
	"encoding/json"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/domain"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/logger"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/util"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

type OrderHandler struct {
	service OrderService
}

type OrderService interface {
	RegisterAccrual(ctx context.Context, order *domain.OrderAccrualRequest) (*domain.SuccessResponse, *domain.Error)
	RegisterWithdrawal(ctx context.Context, order *domain.OrderWithdrawalRequest) (*domain.SuccessResponse, *domain.Error)
	GetAllAccrual(ctx context.Context) ([]domain.OrderAccrual, *domain.Error)
	GetAllWithdrawal(ctx context.Context) ([]domain.OrderWithdrawal, *domain.Error)
	GetUserBalance(ctx context.Context) (*domain.UserBalanceResponse, *domain.Error)
}

func NewOrderHandler(service OrderService) *OrderHandler {

	return &OrderHandler{service: service}
}

func (h *OrderHandler) RegisterAccrualOrder(writer http.ResponseWriter, request *http.Request) {
	enc := json.NewEncoder(writer)
	writer.Header().Set("Content-Type", "application/json")

	order, err := validateOrderAccrualRequest(request)
	logger.Log.Info("validated order", zap.Any("order", order))

	if err != nil {
		sendResponse(enc, err, err.Code, writer)
		return
	}

	response, sErr := h.service.RegisterAccrual(request.Context(), order)
	if sErr != nil {
		sendResponse(enc, sErr, sErr.Code, writer)
		return
	}

	sendResponse(enc, response, response.Code, writer)

}

func (h *OrderHandler) GetAllAccrual(writer http.ResponseWriter, request *http.Request) {
	enc := json.NewEncoder(writer)
	writer.Header().Set("Content-Type", "application/json")

	orders, err := h.service.GetAllAccrual(request.Context())
	logger.Log.Info("orders handler", zap.Any("order", orders))

	if err != nil {
		sendResponse(enc, err, err.Code, writer)
		return
	}

	sendResponse(enc, orders, http.StatusOK, writer)

}

func (h *OrderHandler) RegisterWithdrawalOrder(writer http.ResponseWriter, request *http.Request) {
	enc := json.NewEncoder(writer)
	writer.Header().Set("Content-Type", "application/json")

	order, err := validateOrderWithdrawalRequest(request)
	if err != nil {
		sendResponse(enc, err, err.Code, writer)
		return
	}

	response, sErr := h.service.RegisterWithdrawal(request.Context(), order)
	if sErr != nil {
		sendResponse(enc, sErr, sErr.Code, writer)
		return
	}

	sendResponse(enc, response, response.Code, writer)

}

func (h *OrderHandler) GetAllWithdrawal(writer http.ResponseWriter, request *http.Request) {
	enc := json.NewEncoder(writer)
	writer.Header().Set("Content-Type", "application/json")

	orders, err := h.service.GetAllWithdrawal(request.Context())
	if err != nil {
		sendResponse(enc, err, err.Code, writer)
		return
	}

	sendResponse(enc, orders, http.StatusOK, writer)

}

func (h *OrderHandler) GetUserBalance(writer http.ResponseWriter, request *http.Request) {
	enc := json.NewEncoder(writer)
	writer.Header().Set("Content-Type", "application/json")

	orders, err := h.service.GetUserBalance(request.Context())
	if err != nil {
		sendResponse(enc, err, err.Code, writer)
		return
	}

	sendResponse(enc, orders, http.StatusOK, writer)

}

func validateOrderWithdrawalRequest(request *http.Request) (*domain.OrderWithdrawalRequest, *RestError) {
	var o domain.OrderWithdrawalRequest

	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, &RestError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Err:     err,
		}
	}

	err = json.Unmarshal(body, &o)
	if err != nil {
		return nil, &RestError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Err:     err,
		}
	}

	v := validator.New()
	err = v.Struct(o)
	if err != nil {
		return nil, &RestError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Err:     err,
		}
	}

	orderID, err := strconv.Atoi(o.OrderID)
	if err != nil {
		return nil, &RestError{
			Code:    http.StatusUnprocessableEntity,
			Message: "invalid order format conversion",
			Err:     err,
		}
	}
	if !util.LuhnValid(orderID) {
		return nil, &RestError{
			Code:    http.StatusUnprocessableEntity,
			Message: "invalid order format",
			Err:     err,
		}
	}

	return &o, nil
}
func validateOrderAccrualRequest(request *http.Request) (*domain.OrderAccrualRequest, *RestError) {
	var o domain.OrderAccrualRequest

	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, &RestError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Err:     err,
		}
	}
	var orderID int

	err = json.Unmarshal(body, &orderID)
	if err != nil {
		return nil, &RestError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Err:     err,
		}
	}
	if !util.LuhnValid(orderID) {
		return nil, &RestError{
			Code:    http.StatusUnprocessableEntity,
			Message: "invalid order format",
			Err:     err,
		}
	}

	o.OrderID = strconv.Itoa(orderID)

	v := validator.New()
	err = v.Struct(o)

	if err != nil {
		return nil, &RestError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Err:     err,
		}
	}

	return &o, nil
}
