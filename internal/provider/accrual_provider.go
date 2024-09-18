package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/config"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/domain"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strconv"
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
	GetOrderByUser(ctx context.Context, user int, order uint64) (*domain.OrderAccrual, error)
	GetById(ctx context.Context, order uint64) (*domain.OrderAccrual, error)
	GetOrdersByStatus(ctx context.Context, status []string) ([]domain.OrderAccrual, error)
}
type UserRepository interface {
	Create(ctx context.Context, user domain.User) (*domain.User, error)
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
	GetUserById(ctx context.Context, user int) (*domain.User, error)
	UpdateUserBalance(ctx context.Context, user domain.User) (*domain.User, error)
}

func CreateNewAccrualProvider(userRepository UserRepository, accrualRepository OrderAccrualRepository, config *config.AppConfig) *AccrualProvider {

	fmt.Println("creating accrual provider")
	return &AccrualProvider{
		accrualRepo: accrualRepository,
		userRepo:    userRepository,
		config:      config,
		client:      resty.New(),
	}
}

func (p *AccrualProvider) CreateWorkers(ctx context.Context, orders <-chan domain.OrderAccrual) {
	fmt.Println("creating workers")

	for i := 1; i <= p.config.Worker.Count; i++ {
		go p.worker(ctx, orders, i)
	}
}

func (p *AccrualProvider) Process(ctx context.Context, wg *sync.WaitGroup, ordersData chan<- domain.OrderAccrual) {
	fmt.Println("start processing data")

	defer wg.Done()
	defer close(ordersData)

	interval := 5
	var requestInterval = time.Duration(interval) * time.Second

	for {
		orders, err := p.accrualRepo.GetOrdersByStatus(ctx, []string{domain.OrderStatusNew, domain.OrderStatusProcessed})
		if err != nil {
			//log error
		}
		fmt.Println("orders")
		fmt.Println(orders)

		for _, o := range orders {
			// push to channel data
			ordersData <- o
		}

		time.Sleep(requestInterval)
	}

}

func (p *AccrualProvider) worker(ctx context.Context, orders <-chan domain.OrderAccrual, id int) {
	for order := range orders {
		fmt.Println("started worker")
		fmt.Println(id)
		fmt.Println(order)
		//send request to accrual
		accrualResponse, err := getAccrual(fmt.Sprintf("%s/%d", p.config.AccrualSystemAddress, order.OrderId), p.client)
		if err != nil {
			fmt.Println("accrual request error")
			fmt.Println(err)

			// log error
			continue
		}

		orderId, err := strconv.ParseUint(accrualResponse.Order, 10, 64)
		if err != nil {
			fmt.Println("order decode error")
			fmt.Println(err)
			// log error
			continue
		}
		// send request to accrual service, change status to processing
		// if status processed:  update order status and amount, finish processing
		// if status invalid: update order status , finish processing
		// if status processing: do nothing, push order to queue again
		// if status registered: do nothing, push order to queue again

		if orderId != order.OrderId {
			fmt.Println("wrong order  error")
			fmt.Println(err)
			// log error
			continue

		}
		//processed
		if accrualResponse.Status == domain.OrderStatusProcessed {
			fmt.Println("status processed")
			order.Status = domain.OrderStatusProcessed
			order.Amount = &accrualResponse.Accrual
			_, oErr := p.accrualRepo.Update(ctx, order)
			if oErr != nil {
				// log error
				fmt.Println("order update   error")
				fmt.Println(err)
				continue
			}

			user, uErr := p.userRepo.GetUserById(ctx, int(order.UserId))
			if uErr != nil {
				fmt.Println("user get   error")
				fmt.Println(err)
				// log error
				continue
			}
			user.Balance += accrualResponse.Accrual

			_, uErr = p.userRepo.UpdateUserBalance(ctx, *user)
			if uErr != nil {
				fmt.Println("user update   error")
				fmt.Println(err)
				// log error
				continue
			}
			continue
		}
		//invalid
		if accrualResponse.Status == domain.OrderStatusInvalid {
			fmt.Println("status invalid")

			order.Status = domain.OrderStatusInvalid
			_, oErr := p.accrualRepo.Update(ctx, order)
			if oErr != nil {
				// log error
				fmt.Println("order update  error")
				fmt.Println(err)
				continue
			}
			continue
		}
		if accrualResponse.Status == domain.OrderStatusProcessing ||
			accrualResponse.Status == domain.OrderStatusRegistered {
			fmt.Println("status processing")

			order.Status = domain.OrderStatusProcessing
			_, oErr := p.accrualRepo.Update(ctx, order)
			if oErr != nil {
				// log error
				fmt.Println("order update  error")
				fmt.Println(err)
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
