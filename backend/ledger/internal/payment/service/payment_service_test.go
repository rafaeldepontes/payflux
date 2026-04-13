package service

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rafaeldepontes/ledger/internal/payment/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of payment.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) ProcessPayment(p model.Payment) error {
	args := m.Called(mock.AnythingOfType("model.Payment"))
	return args.Error(0)
}

func (m *MockRepository) GetPaymentByID(id uuid.UUID) (model.Payment, error) {
	args := m.Called(id)
	return args.Get(0).(model.Payment), args.Error(1)
}

func (m *MockRepository) RefundPayment(p model.Payment) error {
	args := m.Called(p)
	return args.Error(0)
}

// MockCache is a mock implementation of cache.Cache
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Add(key string, value string) { m.Called(key, value) }
func (m *MockCache) AddWithTTL(t time.Duration, key string, value ...string) {
	m.Called(t, key, value)
}
func (m *MockCache) Set(key string, value string) { m.Called(key, value) }
func (m *MockCache) Remove(key string)           { m.Called(key) }
func (m *MockCache) Get(key string) (string, bool) {
	args := m.Called(key)
	return args.String(0), args.Bool(1)
}
func (m *MockCache) Clear()     { m.Called() }
func (m *MockCache) FullClear() { m.Called() }

// MockBroker is a mock implementation of payment.MessageBroker
type MockBroker struct {
	mock.Mock
}

func (m *MockBroker) Publish(body []byte) error {
	args := m.Called(mock.Anything)
	return args.Error(0)
}

func TestProcessPayment(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		req         model.PaymentReq
		mockRepoErr error
		mockBrkErr  error
		expectedErr bool
	}{
		{
			name: "Success",
			key:  "1234567890123456",
			req: model.PaymentReq{
				FromAccount: 1,
				ToAccount:   2,
				Amount:      100,
				Currency:    "USD",
			},
			mockRepoErr: nil,
			mockBrkErr:  nil,
			expectedErr: false,
		},
		{
			name: "Validation Error",
			key:  "1234567890123456",
			req: model.PaymentReq{
				FromAccount: 0,
				ToAccount:   2,
				Amount:      100,
				Currency:    "USD",
			},
			expectedErr: true,
		},
		{
			name: "Repository Error",
			key:  "1234567890123456",
			req: model.PaymentReq{
				FromAccount: 1,
				ToAccount:   2,
				Amount:      100,
				Currency:    "USD",
			},
			mockRepoErr: errors.New("db error"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockCache := new(MockCache)
			mockBrk := new(MockBroker)

			if !tt.expectedErr || tt.mockRepoErr != nil {
				if tt.mockRepoErr != nil {
					mockRepo.On("ProcessPayment", mock.Anything).Return(tt.mockRepoErr)
				} else if tt.name == "Success" {
					mockRepo.On("ProcessPayment", mock.Anything).Return(nil)
					mockBrk.On("Publish", mock.Anything).Return(nil)
					mockCache.On("Add", tt.key, mock.Anything).Return()
				}
			}

			s := NewService(mockRepo, mockCache, mockBrk)
			res, err := s.ProcessPayment(tt.key, tt.req)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.req.Amount, res.Amount)
				assert.Equal(t, CompletedStatus, res.Status)
			}
		})
	}
}

func TestCheckKey(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		mockVal     string
		mockHas     bool
		expectedErr bool
	}{
		{
			name:        "Cache Hit",
			key:         "key1",
			mockVal:     `{"payment_id":"id1","status":"completed"}`,
			mockHas:     true,
			expectedErr: false,
		},
		{
			name:        "Cache Miss",
			key:         "key2",
			mockVal:     "",
			mockHas:     false,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := new(MockCache)
			mockCache.On("Get", tt.key).Return(tt.mockVal, tt.mockHas)

			s := NewService(nil, mockCache, nil)
			res, err := s.CheckKey(tt.key)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "completed", res.Status)
			}
		})
	}
}

func TestGetPayment(t *testing.T) {
	id := uuid.New()
	tests := []struct {
		name        string
		id          string
		mockRes     model.Payment
		mockErr     error
		expectedErr bool
	}{
		{
			name: "Success",
			id:   id.String(),
			mockRes: model.Payment{
				ID:     id,
				Status: "completed",
				Amount: 100,
			},
			mockErr:     nil,
			expectedErr: false,
		},
		{
			name:        "Not Found",
			id:          id.String(),
			mockRes:     model.Payment{},
			mockErr:     errors.New("not found"),
			expectedErr: true,
		},
		{
			name:        "Invalid ID",
			id:          "invalid",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			if tt.name != "Invalid ID" {
				mockRepo.On("GetPaymentByID", id).Return(tt.mockRes, tt.mockErr)
			}

			s := NewService(mockRepo, nil, nil)
			res, err := s.GetPayment(tt.id)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, id.String(), res.ID)
			}
		})
	}
}

func TestRefundPayment(t *testing.T) {
	id := uuid.New()
	tests := []struct {
		name        string
		id          string
		req         model.RefundReq
		mockPayment model.Payment
		mockGetErr  error
		mockRefErr  error
		expectedErr bool
	}{
		{
			name: "Success",
			id:   id.String(),
			req:  model.RefundReq{Amount: 100},
			mockPayment: model.Payment{
				ID:     id,
				Status: "completed",
				Amount: 100,
			},
			mockGetErr:  nil,
			mockRefErr:  nil,
			expectedErr: false,
		},
		{
			name: "Already Refunded",
			id:   id.String(),
			mockPayment: model.Payment{
				ID:     id,
				Status: RefundedStatus,
			},
			mockGetErr:  nil,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockBrk := new(MockBroker)
			
			if tt.name != "Invalid ID" {
				mockRepo.On("GetPaymentByID", id).Return(tt.mockPayment, tt.mockGetErr)
				if tt.name == "Success" {
					mockRepo.On("RefundPayment", mock.Anything).Return(nil)
					mockBrk.On("Publish", mock.Anything).Return(nil)
				}
			}

			s := NewService(mockRepo, nil, mockBrk)
			res, err := s.RefundPayment(tt.id, tt.req)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, RefundedStatus, res.Status)
			}
		})
	}
}
