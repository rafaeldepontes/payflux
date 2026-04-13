package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/rafaeldepontes/reconsiliation/internal/reconciliation/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of reconciliation.Service
type MockService struct {
	mock.Mock
}

func (m *MockService) ProcessEvent(event model.PaymentEvent) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockService) GetResult(txID string) (model.ReconciliationResult, error) {
	args := m.Called(txID)
	return args.Get(0).(model.ReconciliationResult), args.Error(1)
}

func (m *MockService) ListExceptions() ([]model.Exception, error) {
	args := m.Called()
	return args.Get(0).([]model.Exception), args.Error(1)
}

func (m *MockService) CreateSettlementRecord(txID string, amount int64, status string) error {
	args := m.Called(txID, amount, status)
	return args.Error(0)
}

func TestGetReconciliationResult(t *testing.T) {
	txID := uuid.New().String()
	tests := []struct {
		name           string
		idParam        string
		mockSvcRes     model.ReconciliationResult
		mockSvcErr     error
		expectedStatus int
	}{
		{
			name:    "Success",
			idParam: txID,
			mockSvcRes: model.ReconciliationResult{
				Status: "matched",
			},
			mockSvcErr:     nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Not Found",
			idParam:        txID,
			mockSvcRes:     model.ReconciliationResult{},
			mockSvcErr:     errors.New("not found"),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			mockSvc.On("GetResult", tt.idParam).Return(tt.mockSvcRes, tt.mockSvcErr)

			c := NewController(mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/reconciliation/"+tt.idParam, nil)
			req.SetPathValue("transaction_id", tt.idParam)

			w := httptest.NewRecorder()
			c.GetReconciliationResult(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestCreateSettlementRecord(t *testing.T) {
	txID := uuid.New().String()
	tests := []struct {
		name           string
		reqBody        interface{}
		mockSvcErr     error
		expectedStatus int
	}{
		{
			name: "Success",
			reqBody: map[string]interface{}{
				"transaction_id": txID,
				"amount":         100,
				"status":         "Settled",
			},
			mockSvcErr:     nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid Body",
			reqBody:        "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.name == "Success" {
				mockSvc.On("CreateSettlementRecord", txID, int64(100), "Settled").Return(tt.mockSvcErr)
			}

			c := NewController(mockSvc)

			body, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/settlements", bytes.NewBuffer(body))

			w := httptest.NewRecorder()
			c.CreateSettlementRecord(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
