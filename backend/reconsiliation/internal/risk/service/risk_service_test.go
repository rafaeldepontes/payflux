package service

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	reconciliation_model "github.com/rafaeldepontes/reconsiliation/internal/reconciliation/model"
	"github.com/rafaeldepontes/reconsiliation/internal/risk/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of risk.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateRiskEvaluation(evaluation model.RiskEvaluation) error {
	args := m.Called(evaluation)
	return args.Error(0)
}

func (m *MockRepository) GetRiskEvaluation(txID uuid.UUID) (model.RiskEvaluation, error) {
	args := m.Called(txID)
	return args.Get(0).(model.RiskEvaluation), args.Error(1)
}

func TestProcessEvent(t *testing.T) {
	txID := uuid.New()
	tests := []struct {
		name          string
		event         reconciliation_model.PaymentEvent
		expectedScore int
		expectedFlags []string
		expectedErr   bool
	}{
		{
			name: "Low Risk",
			event: reconciliation_model.PaymentEvent{
				PaymentID: txID.String(),
				Amount:    100,
			},
			expectedScore: 0,
			expectedFlags: []string{},
			expectedErr:   false,
		},
		{
			name: "High Risk - Large Amount",
			event: reconciliation_model.PaymentEvent{
				PaymentID: txID.String(),
				Amount:    20000,
			},
			expectedScore: 50,
			expectedFlags: []string{"LargeTransactionRule"},
			expectedErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockRepo.On("CreateRiskEvaluation", mock.Anything).Return(nil)

			s := NewService(mockRepo)
			err := s.ProcessEvent(tt.event)

			assert.NoError(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetResult(t *testing.T) {
	txID := uuid.New()
	tests := []struct {
		name        string
		txID        string
		mockRes     model.RiskEvaluation
		mockErr     error
		expectedErr bool
	}{
		{
			name: "Success",
			txID: txID.String(),
			mockRes: model.RiskEvaluation{
				TransactionID: txID,
				RiskScore:     50,
			},
			mockErr:     nil,
			expectedErr: false,
		},
		{
			name:        "Not Found",
			txID:        txID.String(),
			mockRes:     model.RiskEvaluation{},
			mockErr:     errors.New("not found"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockRepo.On("GetRiskEvaluation", txID).Return(tt.mockRes, tt.mockErr)

			s := NewService(mockRepo)
			res, err := s.GetResult(tt.txID)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockRes.RiskScore, res.RiskScore)
			}
		})
	}
}
