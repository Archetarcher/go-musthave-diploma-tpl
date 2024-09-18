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
	GetByID(ctx context.Context, order string) (*domain.OrderAccrual, error)
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
	userID, err := util.GetIDFromToken(ctx)
	logger.Log.Info("user id from claims", zap.Any("userID", userID))

	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}

	order, err := s.accrualRepo.GetByID(ctx, request.OrderID)
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}
	logger.Log.Info("order", zap.Any("order", order))
	logger.Log.Info("order user id ", zap.Any("user", userID))

	if order != nil && int64(userID) != order.UserID {
		return nil, &handlers.RestError{
			Code:    http.StatusConflict,
			Message: "order has been already registered by another user",
			Err:     nil,
		}
	}

	if order != nil && int64(userID) == order.UserID {
		return &domain.SuccessResponse{
			Code:    http.StatusOK,
			Message: "order has been already registered by user",
		}, nil
	}

	accrual := float64(0)
	newOrder, err := s.accrualRepo.Create(ctx, domain.OrderAccrual{UserID: int64(userID), OrderID: request.OrderID, Amount: &accrual, Status: domain.OrderStatusNew})
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
	userID, err := util.GetIDFromToken(ctx)
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}

	orders, err := s.accrualRepo.GetAllByUser(ctx, userID)

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
	userID, err := util.GetIDFromToken(ctx)
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}

	order, err := s.withdrawalRepo.GetOrderByUser(ctx, userID, request.OrderID)
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

	user, err := s.userRepo.GetUserByID(ctx, userID)
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

	_, err = s.withdrawalRepo.Create(ctx, domain.OrderWithdrawal{UserID: int64(userID), OrderID: request.OrderID, Amount: &request.Sum})
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
	userID, err := util.GetIDFromToken(ctx)
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}

	orders, err := s.withdrawalRepo.GetAllByUser(ctx, userID)
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

	userID, err := util.GetIDFromToken(ctx)
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, &handlers.RestError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Err:     err,
		}
	}

	totalWithdrawn, err := s.withdrawalRepo.GetAllSumByUser(ctx, userID)
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
