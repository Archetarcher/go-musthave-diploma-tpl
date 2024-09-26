package pgsql

import (
	"context"
	"github.com/jmoiron/sqlx"
	"reflect"
	"testing"
)

func TestNewStore(t *testing.T) {
	type args struct {
		ctx  context.Context
		conf *Config
	}
	tests := []struct {
		name    string
		args    args
		want    *Store
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStore(tt.args.ctx, tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStore() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_CheckConnection(t *testing.T) {
	type fields struct {
		DB     *sqlx.DB
		config *Config
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				DB:     tt.fields.DB,
				config: tt.fields.config,
			}
			if err := s.CheckConnection(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("CheckConnection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStore_Close(t *testing.T) {
	type fields struct {
		DB     *sqlx.DB
		config *Config
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				DB:     tt.fields.DB,
				config: tt.fields.config,
			}
			s.Close(tt.args.ctx)
		})
	}
}

func TestStore_Restore(t *testing.T) {
	type fields struct {
		DB     *sqlx.DB
		config *Config
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				DB:     tt.fields.DB,
				config: tt.fields.config,
			}
			if err := s.Restore(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStore_RunMigrations(t *testing.T) {
	type fields struct {
		DB     *sqlx.DB
		config *Config
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				DB:     tt.fields.DB,
				config: tt.fields.config,
			}
			if err := s.RunMigrations(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("RunMigrations() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
