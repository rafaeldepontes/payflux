package payment

type Service interface {
	ProcessPayment() (string, error)
}
