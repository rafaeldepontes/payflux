package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	reconciliation_model "github.com/rafaeldepontes/reconsiliation/internal/reconciliation/model"
	"github.com/rafaeldepontes/reconsiliation/internal/risk/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of risk.Service
type MockService struct {
	mock.Mock
}

func (m *MockService) ProcessEvent(event reconciliation_model.PaymentEvent) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockService) GetResult(txID string) (model.RiskEvaluation, error) {
	args := m.Called(txID)
	return args.Get(0).(model.RiskEvaluation), args.Error(1)
}

func TestGetRiskEvaluation(t *testing.T) {
	txID := uuid.New().String()
	tests := []struct {
		name           string
		idParam        string
		mockSvcRes     model.RiskEvaluation
		mockSvcErr     error
		expectedStatus int
	}{
		{
			name:    "Success",
			idParam: txID,
			mockSvcRes: model.RiskEvaluation{
				RiskScore: 50,
			},
			mockSvcErr:     nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Not Found",
			idParam:        txID,
			mockSvcRes:     model.RiskEvaluation{},
			mockSvcErr:     errors.New("not found"),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			mockSvc.On("GetResult", tt.idParam).Return(tt.mockSvcRes, tt.mockSvcErr)

			c := NewController(mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/risk/"+tt.idParam, nil)
			req.SetPathValue("transaction_id", tt.idParam)

			w := httptest.NewRecorder()
			c.GetRiskEvaluation(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
