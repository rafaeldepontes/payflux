package account

import "github.com/rafaeldepontes/ledger/internal/account/model"

type Service interface {
	GetAccountBalance(accountID int) (model.BalanceRes, error)
}
