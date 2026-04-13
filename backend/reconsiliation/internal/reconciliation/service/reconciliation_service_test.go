package service

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/rafaeldepontes/reconsiliation/internal/reconciliation/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of reconciliation.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetSettlementRecord(txID uuid.UUID) (model.SettlementRecord, error) {
	args := m.Called(txID)
	return args.Get(0).(model.SettlementRecord), args.Error(1)
}

func (m *MockRepository) CreateSettlementRecord(record model.SettlementRecord) error {
	args := m.Called(record)
	return args.Error(0)
}

func (m *MockRepository) CreateReconciliationResult(res model.ReconciliationResult) error {
	args := m.Called(res)
	return args.Error(0)
}

func (m *MockRepository) GetReconciliationResult(txID uuid.UUID) (model.ReconciliationResult, error) {
	args := m.Called(txID)
	return args.Get(0).(model.ReconciliationResult), args.Error(1)
}

func (m *MockRepository) CreateException(exc model.Exception) error {
	args := m.Called(exc)
	return args.Error(0)
}

func (m *MockRepository) ListExceptions() ([]model.Exception, error) {
	args := m.Called()
	return args.Get(0).([]model.Exception), args.Error(1)
}

func TestProcessEvent(t *testing.T) {
	txID := uuid.New()
	tests := []struct {
		name           string
		event          model.PaymentEvent
		mockSettlement model.SettlementRecord
		mockSettErr    error
		expectedStatus string
		expectedErr    bool
	}{
		{
			name: "Success Matched",
			event: model.PaymentEvent{
				PaymentID: txID.String(),
				Amount:    100,
			},
			mockSettlement: model.SettlementRecord{
				TransactionID: txID,
				Amount:        100,
			},
			mockSettErr:    nil,
			expectedStatus: "matched",
			expectedErr:    false,
		},
		{
			name: "Mismatched Amount",
			event: model.PaymentEvent{
				PaymentID: txID.String(),
				Amount:    100,
			},
			mockSettlement: model.SettlementRecord{
				TransactionID: txID,
				Amount:        120,
			},
			mockSettErr:    nil,
			expectedStatus: "mismatched",
			expectedErr:    false,
		},
		{
			name: "Missing Settlement",
			event: model.PaymentEvent{
				PaymentID: txID.String(),
				Amount:    100,
			},
			mockSettErr:    sql.ErrNoRows,
			expectedStatus: "missing_settlement",
			expectedErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockRepo.On("GetSettlementRecord", txID).Return(tt.mockSettlement, tt.mockSettErr)
			
			if tt.mockSettErr == sql.ErrNoRows {
				mockRepo.On("CreateException", mock.Anything).Return(nil)
			} else if tt.expectedStatus == "mismatched" {
				mockRepo.On("CreateException", mock.Anything).Return(nil)
			}
			mockRepo.On("CreateReconciliationResult", mock.MatchedBy(func(res model.ReconciliationResult) bool {
				return res.Status == tt.expectedStatus
			})).Return(nil)

			s := NewService(mockRepo)
			err := s.ProcessEvent(tt.event)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetResult(t *testing.T) {
	txID := uuid.New()
	tests := []struct {
		name        string
		txID        string
		mockRes     model.ReconciliationResult
		mockErr     error
		expectedErr bool
	}{
		{
			name: "Success",
			txID: txID.String(),
			mockRes: model.ReconciliationResult{
				TransactionID: txID,
				Status:        "matched",
			},
			mockErr:     nil,
			expectedErr: false,
		},
		{
			name:        "Invalid ID",
			txID:        "invalid",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			if tt.name == "Success" {
				mockRepo.On("GetReconciliationResult", txID).Return(tt.mockRes, tt.mockErr)
			}

			s := NewService(mockRepo)
			res, err := s.GetResult(tt.txID)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockRes.Status, res.Status)
			}
		})
	}
}
