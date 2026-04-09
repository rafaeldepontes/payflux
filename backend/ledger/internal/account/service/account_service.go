package service

import (
	"github.com/rafaeldepontes/ledger/internal/account"
	"github.com/rafaeldepontes/ledger/internal/account/model"
	ar "github.com/rafaeldepontes/ledger/internal/account/repository"
)

type svc struct {
	repo account.Repository
}

func NewService() account.Service {
	return svc{
		repo: ar.NewRepository(),
	}
}

func (s svc) GetAccountBalance(accountID int) (model.BalanceRes, error) {
	balance, err := s.repo.GetAccountBalance(accountID)
	if err != nil {
		return model.BalanceRes{}, err
	}
	return model.BalanceRes{
		AccountID: accountID,
		Balance:   balance,
	}, nil
}
