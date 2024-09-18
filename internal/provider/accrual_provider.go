package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/config"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/domain"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/logger"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"net/http"
	"sync"
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
	GetById(ctx context.Context, order string) (*domain.OrderAccrual, error)
	GetOrdersByStatus(ctx context.Context, status []string) ([]domain.OrderAccrual, error)
}
type UserRepository interface {
	Create(ctx context.Context, user domain.User) (*domain.User, error)
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
	GetUserById(ctx context.Context, user int) (*domain.User, error)
	UpdateUserBalance(ctx context.Context, user domain.User) (*domain.User, error)
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

func (p *AccrualProvider) Process(ctx context.Context, wg *sync.WaitGroup, ordersData chan<- domain.OrderAccrual) {
	logger.Log.Info("start processing data")

	defer wg.Done()
	defer close(ordersData)

	interval := 5
	var requestInterval = time.Duration(interval) * time.Second

	for {
		orders, err := p.accrualRepo.GetOrdersByStatus(ctx, []string{domain.OrderStatusNew, domain.OrderStatusProcessed})
		if err != nil {
			//log error
			logger.Log.Info("error fetching orders by status in provider", zap.Error(err))
		}

		logger.Log.Info("orders by status", zap.Any("order", orders))

		for _, o := range orders {
			// push to channel data
			ordersData <- o
		}

		time.Sleep(requestInterval)
	}

}

func (p *AccrualProvider) worker(ctx context.Context, orders <-chan domain.OrderAccrual, id int) {
	for order := range orders {
		logger.Log.Info("started worker", zap.Any("id", id), zap.Any("order", order))

		//send request to accrual
		accrualResponse, err := getAccrual(fmt.Sprintf("%s/api/orders/%s", p.config.AccrualSystemAddress, order.OrderId), p.client)
		if err != nil {
			logger.Log.Info("accrual request error", zap.Error(err))
			// log error
			continue
		}
		if accrualResponse == nil {
			logger.Log.Info("order not registered in accrual")
			continue
		}

		// send request to accrual service, change status to processing
		// if status processed:  update order status and amount, finish processing
		// if status invalid: update order status , finish processing
		// if status processing: do nothing, push order to queue again
		// if status registered: do nothing, push order to queue again

		if accrualResponse.Order != order.OrderId {
			logger.Log.Info("wrong order  error", zap.Error(err))
			// log error
			continue

		}
		//processed
		if accrualResponse.Status == domain.OrderStatusProcessed {
			logger.Log.Info("status processed")

			order.Status = domain.OrderStatusProcessed
			order.Amount = &accrualResponse.Accrual
			_, oErr := p.accrualRepo.Update(ctx, order)
			if oErr != nil {
				// log error
				logger.Log.Info("order update error", zap.Error(err))
				continue
			}

			user, uErr := p.userRepo.GetUserById(ctx, int(order.UserId))
			if uErr != nil {
				logger.Log.Info("user get error", zap.Error(err))

				// log error
				continue
			}
			*user.Balance += accrualResponse.Accrual

			_, uErr = p.userRepo.UpdateUserBalance(ctx, *user)
			if uErr != nil {
				logger.Log.Info("user update error", zap.Error(err))

				// log error
				continue
			}
			continue
		}
		//invalid
		if accrualResponse.Status == domain.OrderStatusInvalid {
			logger.Log.Info("status invalid")

			order.Status = domain.OrderStatusInvalid
			_, oErr := p.accrualRepo.Update(ctx, order)
			if oErr != nil {
				// log error
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
				// log error
				logger.Log.Info("order update  error", zap.Error(err))
				continue
			}

		}

	}
}

func getAccrual(url string, client *resty.Client) (*domain.AccrualResponse, error) {
	res, err := client.
		R().
		Get(url)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("client: could not create request: %s\n", err.Error()), Time: time.Now(), Err: err}
	}

	if res.StatusCode() == http.StatusNoContent {
		return nil, nil
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
