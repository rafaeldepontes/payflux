package payment

type Service interface {
	ProcessPayment() (string, error)
	CheckKey(key string) (string, error)
}
