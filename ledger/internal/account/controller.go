package account

import "net/http"

type Controller interface {
	GetAccountBalance(w http.ResponseWriter, r *http.Request)
}
