package payment

type Service interface {
	ProcessPayment(key string) (string, error)
	CheckKey(key string) (string, error)
}
