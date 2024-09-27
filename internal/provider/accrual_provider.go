package provider

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/config"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/domain"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/logger"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type AccrualProvider struct {
	accrualRepo OrderAccrualRepository
	userRepo    UserRepository
	config      *config.AppConfig
	client      *resty.Client
}
type OrderAccrualRepository interface {
	Create(ctx context.Context, order domain.OrderAccrual) (*domain.OrderAccrual, error)
	Update(ctx context.Context, order domain.OrderAccrual) (*domain.OrderAccrual, error)
	GetAllByUser(ctx context.Context, user int) ([]domain.OrderAccrual, error)
	GetOrderByUser(ctx context.Context, user int, order string) (*domain.OrderAccrual, error)
	GetByID(ctx context.Context, order string) (*domain.OrderAccrual, error)
	GetOrdersByStatus(ctx context.Context, status []string) ([]domain.OrderAccrual, error)
}
type UserRepository interface {
	Create(ctx context.Context, user domain.User) (*domain.User, error)
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
	GetUserByID(ctx context.Context, user int) (*domain.User, error)
	UpdateUserBalance(ctx context.Context, user domain.User) (*domain.User, error)
	RunInTx(fn func(tx *sql.Tx) *domain.Error) *domain.Error
}

func CreateNewAccrualProvider(userRepository UserRepository, accrualRepository OrderAccrualRepository, config *config.AppConfig) *AccrualProvider {
	logger.Log.Info("creating accrual provider")

	return &AccrualProvider{
		accrualRepo: accrualRepository,
		userRepo:    userRepository,
		config:      config,
		client:      resty.New(),
	}
}

func (p *AccrualProvider) CreateWorkers(ctx context.Context, orders <-chan domain.OrderAccrual) {
	logger.Log.Info("creating workers")

	for i := 1; i <= p.config.Worker.Count; i++ {
		go p.worker(ctx, orders, i)
	}
}

func (p *AccrualProvider) Process(ctx context.Context, ordersData chan<- domain.OrderAccrual) error {
	logger.Log.Info("start processing data")

	defer close(ordersData)

	var interval = time.Duration(p.config.PollInterval) * time.Second

	for {
		orders, err := p.accrualRepo.GetOrdersByStatus(ctx, []string{domain.OrderStatusNew, domain.OrderStatusProcessed})
		if err != nil {
			logger.Log.Info("error fetching orders by status in provider", zap.Error(err))
			return &Error{Message: "error fetching orders by status in provider", Time: time.Now(), Err: err}
		}

		logger.Log.Info("orders by status", zap.Any("order", orders))

		for _, o := range orders {
			ordersData <- o
		}

		time.Sleep(interval)
	}

}

func (p *AccrualProvider) worker(ctx context.Context, orders <-chan domain.OrderAccrual, id int) {
	for order := range orders {
		logger.Log.Info("started worker", zap.Any("id", id), zap.Any("order", order))

		accrualResponse, err := getAccrual(fmt.Sprintf("%s/api/orders/%s", p.config.AccrualSystemAddress, order.OrderID), p.client, p.config.RetryAfter, p.config.RetryCount)
		if err != nil {
			logger.Log.Info("accrual request error", zap.Error(err))
			continue
		}
		if accrualResponse == nil {
			logger.Log.Info("order not registered in accrual")
			continue
		}

		if accrualResponse.Order != order.OrderID {
			logger.Log.Info("wrong order  error", zap.Error(err))
			continue
		}
		if accrualResponse.Status == domain.OrderStatusProcessed {
			logger.Log.Info("status processed")

			tErr := p.userRepo.RunInTx(func(tx *sql.Tx) *domain.Error {
				order.Status = domain.OrderStatusProcessed
				order.Amount = &accrualResponse.Accrual
				_, oErr := p.accrualRepo.Update(ctx, order)
				if oErr != nil {
					return &domain.Error{Message: "order update error", Err: err}
				}

				user, uErr := p.userRepo.GetUserByID(ctx, int(order.UserID))
				if uErr != nil {
					return &domain.Error{Message: "user get error", Err: err}
				}
				*user.Balance += accrualResponse.Accrual

				_, uErr = p.userRepo.UpdateUserBalance(ctx, *user)
				if uErr != nil {
					return &domain.Error{Message: "user update error", Err: err}
				}
				return nil
			})

			if tErr != nil {
				logger.Log.Info("error updating balance", zap.Error(err))
			}
			continue

		}
		if accrualResponse.Status == domain.OrderStatusInvalid {
			logger.Log.Info("status invalid")

			order.Status = domain.OrderStatusInvalid
			_, oErr := p.accrualRepo.Update(ctx, order)
			if oErr != nil {
				logger.Log.Info("user update error", zap.Error(err))

				continue
			}
			continue
		}
		if accrualResponse.Status == domain.OrderStatusProcessing ||
			accrualResponse.Status == domain.OrderStatusRegistered {
			logger.Log.Info("status processing")

			order.Status = domain.OrderStatusProcessing
			_, oErr := p.accrualRepo.Update(ctx, order)
			if oErr != nil {
				logger.Log.Info("order update  error", zap.Error(err))
				continue
			}

		}

	}
}

func getAccrual(url string, client *resty.Client, retryAfter int, retryCount int) (*domain.AccrualResponse, error) {
	res, err := client.
		R().
		Get(url)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("client: could not create request: %s\n", err.Error()), Time: time.Now(), Err: err}
	}

	if res.StatusCode() == http.StatusNoContent {
		return nil, nil
	}
	if res.StatusCode() == http.StatusTooManyRequests {
		time.Sleep(time.Duration(retryAfter) * time.Second)
		if retryCount <= 0 {
			return nil, &Error{Message: fmt.Sprintf("client: responded with error: %s\n", err), Time: time.Now(), Err: err}
		}
		return getAccrual(url, client, retryAfter, retryCount-1)
	}
	if res.StatusCode() != http.StatusOK {
		return nil, &Error{Message: fmt.Sprintf("client: responded with error: %s\n", err), Time: time.Now(), Err: err}
	}
	var response domain.AccrualResponse

	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		return nil, &Error{
			Message: err.Error(),
			Err:     err,
		}
	}
	return &response, nil
}
