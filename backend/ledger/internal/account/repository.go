package account

type Repository interface {
	GetAccountBalance(accountID int) (int64, error)
}
