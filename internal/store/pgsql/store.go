package pgsql

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"time"
)

type Store struct {
	DB     *sqlx.DB
	config *Config
}

func NewStore(ctx context.Context, conf *Config) (*Store, error) {
	db := sqlx.MustOpen("pgx", conf.DatabaseUri)

	store := Store{
		DB:     db,
		config: conf,
	}

	err := store.CheckConnection(ctx)
	if err != nil {
		return nil, err
	}

	if err = store.RunMigrations(ctx); err != nil {
		return nil, err
	}

	return &store, nil
}

func (s *Store) CheckConnection(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := s.DB.PingContext(ctx); err != nil {
		return &Error{
			Time:    time.Now(),
			Message: fmt.Sprintf("%s, %s", ErrorStatusText(StatusDBConnectionException), err.Error()),
			Err:     err,
		}
	}
	return nil
}
func (s *Store) Close(ctx context.Context) {
	err := s.DB.Close()
	if err != nil {
		//log
	}
}
func (s *Store) RunMigrations(ctx context.Context) error {
	db, err := goose.OpenDBWithDriver("pgx", s.config.DatabaseUri)
	if err != nil {
		// log fatal
	}

	if err = goose.RunContext(ctx, "up", db, s.config.MigrationsPath); err != nil {
		return &Error{
			Time:    time.Now(),
			Message: fmt.Sprintf("%s, %s", ErrorStatusText(StatusDBMigrationException), err.Error()),
			Err:     err,
		}
	}
	return nil
}
func (s *Store) Restore(ctx context.Context) error {
	return nil
}
