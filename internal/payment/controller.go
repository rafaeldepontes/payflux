package payment

import "net/http"

type Controller interface {
	ProcessPayment(w http.ResponseWriter, r *http.Request)
}
