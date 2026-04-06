package payment

type Repository interface {
	ProcessPayment(any) (string, error)
}