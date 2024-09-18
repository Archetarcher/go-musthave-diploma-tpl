package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/domain"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/logger"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/store/pgsql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type PGOrderAccrualRepository struct {
	store *pgsql.Store
}

func NewPGOrderAccrualRepository(store *pgsql.Store) *PGOrderAccrualRepository {
	return &PGOrderAccrualRepository{store: store}
}
func (r *PGOrderAccrualRepository) Create(ctx context.Context, order domain.OrderAccrual) (*domain.OrderAccrual, error) {
	logger.Log.Info("order data repo", zap.Any("order", order))
	var orderId string
	rows, err := r.store.DB.NamedQueryContext(ctx, orderAccrualCreateQuery, order)
	if err != nil {
		return nil, &Error{
			Message: fmt.Sprintf("%s, in %s", err.Error(), "PGOrderAccrualRepository Create()"),
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
func (r *PGOrderAccrualRepository) GetById(ctx context.Context, id string) (*domain.OrderAccrual, error) {
	var order domain.OrderAccrual

	err := r.store.DB.GetContext(ctx, &order, orderAccrualGetByIdQuery, id)

	logger.Log.Info("order id", zap.Any("order", order))

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, &Error{
			Message: fmt.Sprintf("%s, in %s", err.Error(), "PGOrderAccrualRepository GetById()"),
			Err:     err,
		}
	}
	return &order, nil
}

func (r *PGOrderAccrualRepository) GetOrderByUser(ctx context.Context, userId int, orderId string) (*domain.OrderAccrual, error) {
	var order domain.OrderAccrual

	err := r.store.DB.GetContext(ctx, &order,
		orderAccrualGetByUserIdQuery, userId, orderId)

	logger.Log.Info("order user", zap.Any("order", order))

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, &Error{
			Message: fmt.Sprintf("%s, in %s", err.Error(), "PGOrderAccrualRepository GetOrderByUser()"),
			Err:     err,
		}
	}
	return &order, nil
}
func (r *PGOrderAccrualRepository) GetAllByUser(ctx context.Context, id int) ([]domain.OrderAccrual, error) {
	var orders []domain.OrderAccrual

	err := r.store.DB.SelectContext(ctx, &orders,
		orderAccrualGetAllByUserIdQuery, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, &Error{
			Message: fmt.Sprintf("%s, in %s", err.Error(), "PGOrderAccrualRepository GetAllByUser()"),
			Err:     err,
		}
	}

	return orders, nil
}
func (r *PGOrderAccrualRepository) GetOrdersByStatus(ctx context.Context, statuses []string) ([]domain.OrderAccrual, error) {
	var orders []domain.OrderAccrual

	q, args, err := sqlx.In(orderAccrualGetOrdersByStatusQuery, statuses)
	if err != nil {
		return nil, &Error{
			Message: fmt.Sprintf("%s, in %s", err.Error(), "PGOrderAccrualRepository GetAllByUser()"),
			Err:     err,
		}
	}

	q = sqlx.Rebind(sqlx.DOLLAR, q)

	err = r.store.DB.SelectContext(ctx, &orders,
		q, args...)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, &Error{
			Message: fmt.Sprintf("%s, in %s", err.Error(), "PGOrderAccrualRepository GetAllByUser()"),
			Err:     err,
		}
	}

	return orders, nil
}
func (r *PGOrderAccrualRepository) Update(ctx context.Context, order domain.OrderAccrual) (*domain.OrderAccrual, error) {

	_, err := r.store.DB.NamedExecContext(ctx, orderAccrualUpdateQuery, order)
	if err != nil {
		return nil, &Error{
			Message: fmt.Sprintf("%s, in %s", err.Error(), "PGOrderAccrualRepository Update()"),
			Err:     err,
		}
	}

	return &order, nil
}
