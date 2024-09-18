package repositories

import (
	"context"
	"database/sql"
	"errors"
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
	var userId int
	rows, err := r.store.DB.NamedQueryContext(ctx, userCreateQuery, user)
	if err != nil {
		return nil, &Error{
			Time:    time.Now(),
			Message: err.Error(),
			Err:     err,
		}
	}
	if rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			return nil, &Error{
				Time:    time.Now(),
				Message: err.Error(),
				Err:     err,
			}
		}
		user.ID = userId
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
func (r *PGUserRepository) GetUserById(ctx context.Context, id int) (*domain.User, error) {

	var user domain.User
	err := r.store.DB.GetContext(ctx, &user,
		userGetByIdQuery, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return &user, nil
}
