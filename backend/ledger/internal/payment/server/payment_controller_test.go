package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rafaeldepontes/ledger/internal/payment/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of payment.Service
type MockService struct {
	mock.Mock
}

func (m *MockService) ProcessPayment(key string, payment model.PaymentReq) (model.PaymentRes, error) {
	args := m.Called(key, payment)
	return args.Get(0).(model.PaymentRes), args.Error(1)
}

func (m *MockService) CheckKey(key string) (model.PaymentRes, error) {
	args := m.Called(key)
	return args.Get(0).(model.PaymentRes), args.Error(1)
}

func (m *MockService) GetPayment(id string) (model.PaymentRes, error) {
	args := m.Called(id)
	return args.Get(0).(model.PaymentRes), args.Error(1)
}

func (m *MockService) RefundPayment(id string, req model.RefundReq) (model.PaymentRes, error) {
	args := m.Called(id, req)
	return args.Get(0).(model.PaymentRes), args.Error(1)
}

func TestProcessPayment(t *testing.T) {
	tests := []struct {
		name           string
		idempotencyKey string
		reqBody        interface{}
		mockCheckRes   model.PaymentRes
		mockCheckErr   error
		mockProcRes    model.PaymentRes
		mockProcErr    error
		expectedStatus int
	}{
		{
			name:           "Success",
			idempotencyKey: "1234567890123456",
			reqBody: model.PaymentReq{
				FromAccount: 1,
				ToAccount:   2,
				Amount:      100,
				Currency:    "USD",
			},
			mockCheckRes:   model.PaymentRes{},
			mockCheckErr:   errors.New("not on cache"),
			mockProcRes:    model.PaymentRes{ID: "pay1", Status: "completed"},
			mockProcErr:    nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing Idempotency-Key",
			idempotencyKey: "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid Idempotency-Key Length",
			idempotencyKey: "too-short",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Cache Hit",
			idempotencyKey: "1234567890123456",
			mockCheckRes:   model.PaymentRes{ID: "pay1", Status: "completed"},
			mockCheckErr:   nil,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.idempotencyKey != "" && len(tt.idempotencyKey) == 16 {
				mockSvc.On("CheckKey", tt.idempotencyKey).Return(tt.mockCheckRes, tt.mockCheckErr)
				if tt.mockCheckErr != nil && tt.name == "Success" {
					mockSvc.On("ProcessPayment", tt.idempotencyKey, tt.reqBody).Return(tt.mockProcRes, tt.mockProcErr)
				}
			}

			c := NewController(mockSvc)

			body, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer(body))
			if tt.idempotencyKey != "" {
				req.Header.Set("Idempotency-Key", tt.idempotencyKey)
			}

			w := httptest.NewRecorder()
			c.ProcessPayment(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestGetPayment(t *testing.T) {
	tests := []struct {
		name           string
		idParam        string
		mockSvcRes     model.PaymentRes
		mockSvcErr     error
		expectedStatus int
	}{
		{
			name:           "Success",
			idParam:        "pay1",
			mockSvcRes:     model.PaymentRes{ID: "pay1", Status: "completed"},
			mockSvcErr:     nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Not Found",
			idParam:        "pay2",
			mockSvcRes:     model.PaymentRes{},
			mockSvcErr:     errors.New("not found"),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			mockSvc.On("GetPayment", tt.idParam).Return(tt.mockSvcRes, tt.mockSvcErr)

			c := NewController(mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/payments/"+tt.idParam, nil)
			req.SetPathValue("id", tt.idParam)

			w := httptest.NewRecorder()
			c.GetPayment(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
