package services

import (
	"context"
	"fmt"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/domain"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/handlers"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/util"
	"net/http"
	"strconv"
)

type OrderService struct {
	accrualRepo    OrderAccrualRepository
	withdrawalRepo OrderWithdrawalRepository
	userRepo       UserRepository
}
type OrderAccrualRepository interface {
	Create(ctx context.Context, order domain.OrderAccrual) (*domain.OrderAccrual, error)
	Update(ctx context.Context, order domain.OrderAccrual) (*domain.OrderAccrual, error)
	GetAllByUser(ctx context.Context, user int) ([]domain.OrderAccrual, error)
	GetOrderByUser(ctx context.Context, user int, order uint64) (*domain.OrderAccrual, error)
	GetById(ctx context.Context, order uint64) (*domain.OrderAccrual, error)
	GetOrdersByStatus(ctx context.Context, status []string) ([]domain.OrderAccrual, error)
}

type OrderWithdrawalRepository interface {
	Create(ctx context.Context, order domain.OrderWithdrawal) (*domain.OrderWithdrawal, error)
	GetAllByUser(ctx context.Context, user int) ([]domain.OrderWithdrawal, error)
	GetOrderByUser(ctx context.Context, user int, order uint64) (*domain.OrderWithdrawal, error)
	GetAllSumByUser(ctx context.Context, user int) (float64, error)
}

func NewOrderService(accrualRepo OrderAccrualRepository, withdrawalRepo OrderWithdrawalRepository, userRepo UserRepository) *OrderService {
	return &OrderService{accrualRepo: accrualRepo, withdrawalRepo: withdrawalRepo, userRepo: userRepo}
}

func (s *OrderService) RegisterAccrual(ctx context.Context, request *domain.OrderAccrualRequest) (*domain.SuccessResponse, *handlers.RestError) {
	userId, err := util.GetIdFromToken(ctx)
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}
	orderId, err := strconv.ParseUint(request.OrderId, 10, 64)
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusUnprocessableEntity,
			Message: err.Error(),
			Err:     err,
		}
	}

	order, err := s.accrualRepo.GetById(ctx, orderId)
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}
	fmt.Println("order")
	fmt.Println(order)

	if order != nil && int64(userId) != order.UserId {
		return nil, &handlers.RestError{
			Code:    http.StatusConflict,
			Message: "order has been already registered by another user",
			Err:     nil,
		}
	}

	if order != nil && int64(userId) == order.UserId {
		return &domain.SuccessResponse{
			Code:    http.StatusOK,
			Message: "order has been already registered by user",
		}, nil
	}

	accrual := float64(100)
	newOrder, err := s.accrualRepo.Create(ctx, domain.OrderAccrual{UserId: int64(userId), OrderId: orderId, Amount: &accrual, Status: domain.OrderStatusNew})
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}

	return &domain.SuccessResponse{
		Code:         http.StatusAccepted,
		Message:      "order accrual has been accepted to processing",
		OrderAccrual: newOrder,
	}, nil

}
func (s *OrderService) GetAllAccrual(ctx context.Context) ([]domain.OrderAccrual, *handlers.RestError) {
	userId, err := util.GetIdFromToken(ctx)
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}

	orders, err := s.accrualRepo.GetAllByUser(ctx, userId)
	fmt.Println("orders")
	fmt.Println(orders)
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}

	if orders == nil {
		return nil, &handlers.RestError{
			Code:    http.StatusNoContent,
			Message: "user does not have any accrual orders registered",
			Err:     nil,
		}
	}
	fmt.Println("ordersssss")
	fmt.Println(orders)
	return orders, nil

}
func (s *OrderService) RegisterWithdrawal(ctx context.Context, request *domain.OrderWithdrawalRequest) (*domain.SuccessResponse, *handlers.RestError) {
	userId, err := util.GetIdFromToken(ctx)
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}

	orderId, err := strconv.ParseUint(request.OrderId, 10, 64)
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusUnprocessableEntity,
			Message: err.Error(),
			Err:     err,
		}
	}

	order, err := s.withdrawalRepo.GetOrderByUser(ctx, userId, orderId)
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}
	if order != nil {
		return &domain.SuccessResponse{
			Code:    http.StatusOK,
			Message: "order withdrawal has been already registered by user",
		}, nil
	}

	user, err := s.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}

	if (user.Balance - request.Sum) < 0 {
		return nil, &handlers.RestError{
			Code:    http.StatusPaymentRequired,
			Message: "low balance",
			Err:     err,
		}
	}

	_, err = s.withdrawalRepo.Create(ctx, domain.OrderWithdrawal{UserId: int64(userId), OrderId: orderId, Amount: &request.Sum})
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}
	user.Balance -= request.Sum
	_, err = s.userRepo.UpdateUserBalance(ctx, *user)
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}

	return &domain.SuccessResponse{
		Code:    http.StatusOK,
		Message: "order withdrawal has been processed",
	}, nil

}
func (s *OrderService) GetAllWithdrawal(ctx context.Context) ([]domain.OrderWithdrawal, *handlers.RestError) {
	userId, err := util.GetIdFromToken(ctx)
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}

	orders, err := s.withdrawalRepo.GetAllByUser(ctx, userId)
	fmt.Println("orders")
	fmt.Println(orders)
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}

	if orders == nil {
		return nil, &handlers.RestError{
			Code:    http.StatusNoContent,
			Message: "user does not have any withdrawal orders registered",
			Err:     nil,
		}
	}

	return orders, nil

}
func (s *OrderService) GetUserBalance(ctx context.Context) (*domain.UserBalanceResponse, *handlers.RestError) {
	fmt.Println("GetUserBalance")

	userId, err := util.GetIdFromToken(ctx)
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}
	user, err := s.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}
	fmt.Println(user)

	totalWithdrawn, err := s.withdrawalRepo.GetAllSumByUser(ctx, userId)
	fmt.Println("totalWithdrawal")
	fmt.Println(totalWithdrawn)
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}

	return &domain.UserBalanceResponse{
		Current:   user.Balance,
		Withdrawn: totalWithdrawn,
	}, nil

}
