package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rafaeldepontes/ledger/internal/account/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of account.Service
type MockService struct {
	mock.Mock
}

func (m *MockService) GetAccountBalance(accountID int) (model.BalanceRes, error) {
	args := m.Called(accountID)
	return args.Get(0).(model.BalanceRes), args.Error(1)
}

func TestGetAccountBalance(t *testing.T) {
	tests := []struct {
		name           string
		idParam        string
		mockSvcRes     model.BalanceRes
		mockSvcErr     error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:    "Success",
			idParam: "1",
			mockSvcRes: model.BalanceRes{
				AccountID: 1,
				Balance:   1000,
			},
			mockSvcErr:     nil,
			expectedStatus: http.StatusOK,
			expectedBody: model.BalanceRes{
				AccountID: 1,
				Balance:   1000,
			},
		},
		{
			name:           "Invalid ID",
			idParam:        "abc",
			mockSvcRes:     model.BalanceRes{},
			mockSvcErr:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "invalid account id"},
		},
		{
			name:           "Service Error",
			idParam:        "1",
			mockSvcRes:     model.BalanceRes{},
			mockSvcErr:     errors.New("internal error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{"message": "internal error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.idParam == "1" {
				mockSvc.On("GetAccountBalance", 1).Return(tt.mockSvcRes, tt.mockSvcErr)
			}

			c := NewController(mockSvc)
			
			req := httptest.NewRequest(http.MethodGet, "/accounts/"+tt.idParam+"/balance", nil)
			// Mock the path value for Go 1.22+ mux
			req.SetPathValue("id", tt.idParam)
			
			w := httptest.NewRecorder()

			c.GetAccountBalance(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			
			var actualBody interface{}
			if tt.expectedStatus == http.StatusOK {
				var res model.BalanceRes
				json.NewDecoder(w.Body).Decode(&res)
				actualBody = res
			} else {
				var res map[string]string
				json.NewDecoder(w.Body).Decode(&res)
				actualBody = res
			}
			assert.Equal(t, tt.expectedBody, actualBody)
			
			if tt.idParam == "1" {
				mockSvc.AssertExpectations(t)
			}
		})
	}
}
