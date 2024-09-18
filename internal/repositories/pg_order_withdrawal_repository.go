package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/domain"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/logger"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/store/pgsql"
	"go.uber.org/zap"
)

type PGOrderWithdrawalRepository struct {
	store *pgsql.Store
}

func NewPGOrderWithdrawalRepository(store *pgsql.Store) *PGOrderWithdrawalRepository {
	return &PGOrderWithdrawalRepository{store: store}
}
func (r *PGOrderWithdrawalRepository) Create(ctx context.Context, order domain.OrderWithdrawal) (*domain.OrderWithdrawal, error) {
	var orderId string
	rows, err := r.store.DB.NamedQueryContext(ctx, orderWithdrawalCreateQuery, order)
	if err != nil {
		return nil, &Error{
			Message: fmt.Sprintf("%s, in %s", err.Error(), "PGOrderWithdrawalRepository Create()"),
			Err:     err,
		}
	}
	if rows.Next() {
		err = rows.Scan(&orderId)
		if err != nil {
			return nil, &Error{
				Message: fmt.Sprintf("%s, in %s", err.Error(), "PGOrderAccrualRepository Create()"),
				Err:     err,
			}
		}
		order.OrderId = orderId
	}

	return &order, nil
}
func (r *PGOrderWithdrawalRepository) GetOrderByUser(ctx context.Context, userId int, orderId string) (*domain.OrderWithdrawal, error) {
	var order domain.OrderWithdrawal

	err := r.store.DB.GetContext(ctx, &order,
		orderWithdrawalGetByUserIdQuery, userId, orderId)

	logger.Log.Info("order user", zap.Any("order", order))

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, &Error{
			Message: fmt.Sprintf("%s, in %s", err.Error(), "PGOrderWithdrawalRepository GetOrderByUser()"),
			Err:     err,
		}
	}
	return &order, nil
}

func (r *PGOrderWithdrawalRepository) GetAllByUser(ctx context.Context, id int) ([]domain.OrderWithdrawal, error) {
	var orders []domain.OrderWithdrawal

	err := r.store.DB.SelectContext(ctx, &orders,
		orderWithdrawalGetAllByUserIdQuery, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, &Error{
			Message: fmt.Sprintf("%s, in %s", err.Error(), "PGOrderWithdrawalRepository GetAllByUser()"),
			Err:     err,
		}
	}

	return orders, nil
}

func (r *PGOrderWithdrawalRepository) GetAllSumByUser(ctx context.Context, id int) (float64, error) {
	var sum float64

	err := r.store.DB.GetContext(ctx, &sum,
		orderWithdrawalGetAllByUserSumIdQuery, id)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}

	if err != nil {
		return 0, &Error{
			Message: fmt.Sprintf("%s, in %s", err.Error(), "PGOrderWithdrawalRepository GetAllSumByUser()"),
			Err:     err,
		}
	}

	return sum, nil
}
