package handlers

import (
	"context"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/config"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/domain"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/repositories"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/services"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/store/pgsql"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var authRequest = domain.AuthRequest{Login: "test", Password: "password"}
var c = config.NewConfig()

func setup() *httptest.Server {
	c.ParseConfig()
	util.GenerateAuthToken(c)

	ctx := context.Background()

	storage, err := pgsql.NewStore(ctx, &pgsql.Config{DatabaseURI: c.DatabaseURI, MigrationsPath: c.MigrationsPath})
	if err != nil {
		log.Fatal("failed to init storage with error", err)
		return nil
	}
	userRepository := repositories.NewPGUserRepository(storage)
	orderAccrualRepository := repositories.NewPGOrderAccrualRepository(storage)
	orderWithdrawalRepository := repositories.NewPGOrderWithdrawalRepository(storage)

	orderService := services.NewOrderService(orderAccrualRepository, orderWithdrawalRepository, userRepository)
	authService := services.NewAuthService(userRepository, c.Token)

	orderHandler := NewOrderHandler(orderService)
	authHandler := NewAuthHandler(authService)

	r := chi.NewRouter()
	r.Post("/api/user/register", authHandler.Register)
	r.Post("/api/user/login", authHandler.Login)
	r.Group(func(userRouter chi.Router) {
		userRouter.Use(jwtauth.Verifier(c.Token.AuthToken))
		userRouter.Use(jwtauth.Authenticator(c.Token.AuthToken))

		userRouter.Post("/api/user/orders", orderHandler.RegisterAccrualOrder)
		userRouter.Get("/api/user/orders", orderHandler.GetAllAccrual)

		userRouter.Get("/api/user/balance", orderHandler.GetUserBalance)
		userRouter.Post("/api/user/balance/withdraw", orderHandler.RegisterWithdrawalOrder)
		userRouter.Get("/api/user/withdrawals", orderHandler.GetAllWithdrawal)

	})
	srv := httptest.NewServer(r)

	return srv

}

func TestAuthHandler_Register(t *testing.T) {
	srv := setup()

	defer srv.Close()
	type request struct {
		query  string
		method string
		body   map[string]string
	}
	type want struct {
		code int
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "200 — пользователь успешно зарегистрирован и аутентифицирован",
			request: request{
				query:  "/api/user/register",
				method: http.MethodPost,
				body:   map[string]string{"login": authRequest.Login, "password": authRequest.Password},
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "400 — неверный формат запроса",
			request: request{
				query:  "/api/user/register",
				method: http.MethodPost,
				body:   map[string]string{"login": authRequest.Login},
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "409 — логин уже занят",
			request: request{
				query:  "/api/user/register",
				method: http.MethodPost,
				body:   map[string]string{"login": authRequest.Login, "password": authRequest.Password},
			},
			want: want{
				code: http.StatusConflict,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.request.method
			req.URL = srv.URL + tt.request.query
			req.SetBody(tt.request.body)

			resp, err := req.Send()

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.want.code, resp.StatusCode(), "Код ответа не совпадает с ожидаемым")
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	srv := setup()

	defer srv.Close()
	type request struct {
		query  string
		method string
		body   map[string]string
	}
	type want struct {
		code int
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "200 — пользователь успешно аутентифицирован",
			request: request{
				query:  "/api/user/login",
				method: http.MethodPost,
				body:   map[string]string{"login": authRequest.Login, "password": authRequest.Password},
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "400 — неверный формат запроса",
			request: request{
				query:  "/api/user/login",
				method: http.MethodPost,
				body:   map[string]string{"login": authRequest.Login},
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "401 — неверная пара логин/пароль",
			request: request{
				query:  "/api/user/login",
				method: http.MethodPost,
				body:   map[string]string{"login": authRequest.Login, "password": authRequest.Login},
			},
			want: want{
				code: http.StatusUnauthorized,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.request.method
			req.URL = srv.URL + tt.request.query
			req.SetBody(tt.request.body)

			resp, err := req.Send()

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.want.code, resp.StatusCode(), "Код ответа не совпадает с ожидаемым")

		})
	}
}
