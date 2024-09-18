package services

import (
	"context"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/domain"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/handlers"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/logger"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/util"
	"go.uber.org/zap"
	"net/http"
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
	GetOrderByUser(ctx context.Context, user int, order string) (*domain.OrderAccrual, error)
	GetById(ctx context.Context, order string) (*domain.OrderAccrual, error)
	GetOrdersByStatus(ctx context.Context, status []string) ([]domain.OrderAccrual, error)
}

type OrderWithdrawalRepository interface {
	Create(ctx context.Context, order domain.OrderWithdrawal) (*domain.OrderWithdrawal, error)
	GetAllByUser(ctx context.Context, user int) ([]domain.OrderWithdrawal, error)
	GetOrderByUser(ctx context.Context, user int, order string) (*domain.OrderWithdrawal, error)
	GetAllSumByUser(ctx context.Context, user int) (float64, error)
}

func NewOrderService(accrualRepo OrderAccrualRepository, withdrawalRepo OrderWithdrawalRepository, userRepo UserRepository) *OrderService {
	return &OrderService{accrualRepo: accrualRepo, withdrawalRepo: withdrawalRepo, userRepo: userRepo}
}

func (s *OrderService) RegisterAccrual(ctx context.Context, request *domain.OrderAccrualRequest) (*domain.SuccessResponse, *handlers.RestError) {
	userId, err := util.GetIdFromToken(ctx)
	logger.Log.Info("user id from claims", zap.Any("userId", userId))

	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}

	order, err := s.accrualRepo.GetById(ctx, request.OrderId)
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}
	logger.Log.Info("order", zap.Any("order", order))
	logger.Log.Info("order user id ", zap.Any("user", userId))

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

	accrual := float64(0)
	newOrder, err := s.accrualRepo.Create(ctx, domain.OrderAccrual{UserId: int64(userId), OrderId: request.OrderId, Amount: &accrual, Status: domain.OrderStatusNew})
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

	logger.Log.Info("orders", zap.Any("order", orders))

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
	logger.Log.Info("orders", zap.Any("order", orders))

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

	order, err := s.withdrawalRepo.GetOrderByUser(ctx, userId, request.OrderId)
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

	if (*user.Balance - request.Sum) < 0 {
		return nil, &handlers.RestError{
			Code:    http.StatusPaymentRequired,
			Message: "low balance",
			Err:     err,
		}
	}

	_, err = s.withdrawalRepo.Create(ctx, domain.OrderWithdrawal{UserId: int64(userId), OrderId: request.OrderId, Amount: &request.Sum})
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}
	*user.Balance -= request.Sum
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
	logger.Log.Info("orders", zap.Any("order", orders))

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
	logger.Log.Info("GetUserBalance")

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

	totalWithdrawn, err := s.withdrawalRepo.GetAllSumByUser(ctx, userId)
	logger.Log.Info("totalWithdrawal", zap.Any("totalWithdrawal", totalWithdrawn))
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}

	return &domain.UserBalanceResponse{
		Current:   *user.Balance,
		Withdrawn: totalWithdrawn,
	}, nil

}
