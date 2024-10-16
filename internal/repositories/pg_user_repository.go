package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/domain"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/store/pgsql"
	"time"
)

type PGUserRepository struct {
	store *pgsql.Store
}

func NewPGUserRepository(store *pgsql.Store) *PGUserRepository {
	return &PGUserRepository{store: store}
}
func (r *PGUserRepository) Create(ctx context.Context, user domain.User) (*domain.User, error) {
	var userID int
	rows, err := r.store.DB.NamedQueryContext(ctx, userCreateQuery, user)
	if err != nil {
		return nil, &Error{
			Time:    time.Now(),
			Message: err.Error(),
			Err:     err,
		}
	}
	if rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			return nil, &Error{
				Time:    time.Now(),
				Message: err.Error(),
				Err:     err,
			}
		}
		user.ID = userID
	}
	err = rows.Err()
	if err != nil {
		return nil, &Error{
			Message: fmt.Sprintf("%s, in %s", err.Error(), "PGOrderAccrualRepository Create()"),
			Err:     err,
		}
	}

	return &user, nil
}
func (r *PGUserRepository) UpdateUserBalance(ctx context.Context, user domain.User) (*domain.User, error) {

	_, err := r.store.DB.NamedExecContext(ctx, userUpdateQuery, map[string]interface{}{
		"balance": user.Balance,
		"id":      user.ID,
	})
	if err != nil {
		return nil, &Error{
			Time:    time.Now(),
			Message: err.Error(),
			Err:     err,
		}
	}

	return &user, nil
}
func (r *PGUserRepository) GetUserByLogin(ctx context.Context, login string) (*domain.User, error) {

	var user domain.User
	err := r.store.DB.GetContext(ctx, &user,
		userGetByLoginQuery, login)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *PGUserRepository) GetUserByID(ctx context.Context, id int) (*domain.User, error) {

	var user domain.User
	err := r.store.DB.GetContext(ctx, &user,
		userGetByIDQuery, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PGUserRepository) RunInTx(fn func(tx *sql.Tx) *domain.Error) *domain.Error {
	tx, err := r.store.DB.Begin()
	if err != nil {
		return &domain.Error{
			Message: fmt.Sprintf("%s, %s", ErrorStatusText(StatusDBTransactionException), err.Error()),
			Code:    StatusDBTransactionException,
		}
	}

	tErr := fn(tx)
	if tErr == nil {
		cErr := tx.Commit()
		if cErr != nil {
			return &domain.Error{
				Message: fmt.Sprintf("%s, %s", ErrorStatusText(StatusDBTransactionException), cErr.Error()),
				Code:    StatusDBTransactionException,
			}
		}
		return nil
	}

	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		jErr := errors.Join(tErr, rollbackErr)
		return &domain.Error{Code: tErr.Code, Message: tErr.Message, Err: jErr}
	}

	return nil
}
