package account

import "github.com/rafaeldepontes/goplo/internal/account/model"

type Service interface {
	GetAccountBalance(accountID int) (model.BalanceRes, error)
}
