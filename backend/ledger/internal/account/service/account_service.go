package service

import (
	"github.com/rafaeldepontes/ledger/internal/account"
	"github.com/rafaeldepontes/ledger/internal/account/model"
)

type svc struct {
	repo account.Repository
}

func NewService(repo account.Repository) account.Service {
	return svc{
		repo: repo,
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
