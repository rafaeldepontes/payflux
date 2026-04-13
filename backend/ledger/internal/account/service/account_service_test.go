package service

import (
	"errors"
	"testing"

	"github.com/rafaeldepontes/ledger/internal/account/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of account.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetAccountBalance(accountID int) (int64, error) {
	args := m.Called(accountID)
	return args.Get(0).(int64), args.Error(1)
}

func TestGetAccountBalance(t *testing.T) {
	tests := []struct {
		name          string
		accountID     int
		mockBalance   int64
		mockErr       error
		expectedRes   model.BalanceRes
		expectedErr   bool
	}{
		{
			name:        "Success",
			accountID:   1,
			mockBalance: 1000,
			mockErr:     nil,
			expectedRes: model.BalanceRes{
				AccountID: 1,
				Balance:   1000,
			},
			expectedErr: false,
		},
		{
			name:        "Repository Error",
			accountID:   2,
			mockBalance: 0,
			mockErr:     errors.New("db error"),
			expectedRes: model.BalanceRes{},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockRepo.On("GetAccountBalance", tt.accountID).Return(tt.mockBalance, tt.mockErr)

			s := NewService(mockRepo)
			res, err := s.GetAccountBalance(tt.accountID)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRes, res)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
