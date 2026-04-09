package payment

import "net/http"

type Controller interface {
	ProcessPayment(w http.ResponseWriter, r *http.Request)

	// TODO:impl get payment and refund payment
}
