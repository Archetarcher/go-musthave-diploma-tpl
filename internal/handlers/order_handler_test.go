package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/domain"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupToken(srv *httptest.Server) (error, *domain.AuthResponse) {
	res, err := resty.New().R().SetBody(authRequest).Post(fmt.Sprintf("%s/api/user/login", srv.URL))
	if err != nil {
		return err, nil
	}
	var response domain.AuthResponse

	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		return err, nil
	}
	return nil, &response
}
func init() {

}
func TestOrderHandler_RegisterAccrualOrder(t *testing.T) {
	srv := setup()

	err, response := setupToken(srv)
	require.NoError(t, err, "Не удалось авторизовать пользователя", srv, err)
	orderId := 656505266020341
	invalidOrderId := 12345678902

	defer srv.Close()

	type request struct {
		query   string
		method  string
		headers map[string]string
		body    string
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
			name: "202 — новый номер заказа принят в обработку",
			request: request{
				query:   "/api/user/orders",
				method:  http.MethodPost,
				headers: map[string]string{"Authorization": "Bearer " + response.Token},
				body:    fmt.Sprintf("%d", orderId),
			},
			want: want{
				code: http.StatusAccepted,
			},
		},
		{
			name: "200 — номер заказа уже был загружен этим пользователем",
			request: request{
				query:   "/api/user/orders",
				method:  http.MethodPost,
				headers: map[string]string{"Authorization": "Bearer " + response.Token},
				body:    fmt.Sprintf("%d", orderId),
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "400 — неверный формат запроса",
			request: request{
				query:   "/api/user/orders",
				method:  http.MethodPost,
				headers: map[string]string{"Authorization": "Bearer " + response.Token},
				body:    fmt.Sprintf("\"%d\"", orderId),
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "401 — пользователь не аутентифицирован",
			request: request{
				query:   "/api/user/orders",
				method:  http.MethodPost,
				headers: map[string]string{},
				body:    fmt.Sprintf("%d", orderId),
			},
			want: want{
				code: http.StatusUnauthorized,
			},
		},
		{
			name: "422 — неверный формат номера заказа",
			request: request{
				query:   "/api/user/orders",
				method:  http.MethodPost,
				headers: map[string]string{"Authorization": "Bearer " + response.Token},
				body:    fmt.Sprintf("%d", invalidOrderId),
			},
			want: want{
				code: http.StatusUnprocessableEntity,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.request.method
			req.URL = srv.URL + tt.request.query
			req.SetBody(tt.request.body)
			req.SetHeaders(tt.request.headers)

			resp, err := req.Send()

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.want.code, resp.StatusCode(), "Код ответа не совпадает с ожидаемым", srv.URL+tt.request.query)

		})
	}
}

func TestOrderHandler_RegisterWithdrawalOrder(t *testing.T) {
	srv := setup()
	err, response := setupToken(srv)

	require.NoError(t, err, "Не удалось авторизовать пользователя")
	invalidOrderId := 12345678902
	orderId := 656505266020341

	defer srv.Close()

	type request struct {
		query   string
		method  string
		headers map[string]string
		body    domain.OrderWithdrawalRequest
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
			name: "402 — на счету недостаточно средств",
			request: request{
				query:   "/api/user/balance/withdraw",
				method:  http.MethodPost,
				headers: map[string]string{"Authorization": "Bearer " + response.Token},
				body:    domain.OrderWithdrawalRequest{OrderID: fmt.Sprintf("%d", orderId), Sum: rand.Float64()},
			},
			want: want{
				code: http.StatusPaymentRequired,
			},
		},
		{
			name: "401 — пользователь не аутентифицирован",
			request: request{
				query:   "/api/user/balance/withdraw",
				method:  http.MethodPost,
				headers: map[string]string{},
				body:    domain.OrderWithdrawalRequest{},
			},
			want: want{
				code: http.StatusUnauthorized,
			},
		},
		{
			name: "422 — неверный формат номера заказа",
			request: request{
				query:   "/api/user/balance/withdraw",
				method:  http.MethodPost,
				headers: map[string]string{"Authorization": "Bearer " + response.Token},
				body:    domain.OrderWithdrawalRequest{OrderID: fmt.Sprintf("%d", invalidOrderId), Sum: rand.Float64()},
			},
			want: want{
				code: http.StatusUnprocessableEntity,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.request.method
			req.URL = srv.URL + tt.request.query

			req.SetBody(tt.request.body)
			req.SetHeaders(tt.request.headers)

			resp, err := req.Send()

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.want.code, resp.StatusCode(), "Код ответа не совпадает с ожидаемым", resp.String())

		})
	}
}

func TestOrderHandler_GetAllAccrual(t *testing.T) {
	srv := setup()

	err, response := setupToken(srv)
	require.NoError(t, err, "Не удалось получить токен")

	defer srv.Close()

	type request struct {
		query   string
		method  string
		headers map[string]string
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
			name: "200 — успешная обработка запроса",
			request: request{
				query:   "/api/user/orders",
				method:  http.MethodGet,
				headers: map[string]string{"Authorization": "Bearer " + response.Token},
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "401 — пользователь не авторизован",
			request: request{
				query:   "/api/user/orders",
				method:  http.MethodGet,
				headers: map[string]string{},
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
			req.SetHeaders(tt.request.headers)

			resp, err := req.Send()

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.want.code, resp.StatusCode(), "Код ответа не совпадает с ожидаемым")

		})
	}
}

func TestOrderHandler_GetAllWithdrawal(t *testing.T) {
	srv := setup()

	err, response := setupToken(srv)
	require.NoError(t, err, "Не удалось получить токен")

	defer srv.Close()

	type request struct {
		query   string
		method  string
		headers map[string]string
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
			name: "204 — нет ни одного списания",
			request: request{
				query:   "/api/user/withdrawals",
				method:  http.MethodGet,
				headers: map[string]string{"Authorization": "Bearer " + response.Token},
			},
			want: want{
				code: http.StatusNoContent,
			},
		},
		{
			name: "401 — пользователь не авторизован",
			request: request{
				query:   "/api/user/withdrawals",
				method:  http.MethodGet,
				headers: map[string]string{},
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
			req.SetHeaders(tt.request.headers)

			resp, err := req.Send()

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.want.code, resp.StatusCode(), "Код ответа не совпадает с ожидаемым")

		})
	}
}
