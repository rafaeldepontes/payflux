package payment

import "net/http"

type Controller interface {
	ProcessPayment(w http.ResponseWriter, r *http.Request)
	GetPayment(w http.ResponseWriter, r *http.Request)
	RefundPayment(w http.ResponseWriter, r *http.Request)
}
