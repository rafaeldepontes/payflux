package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	PaymentRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "payment_requests_total",
		Help: "The total number of payment requests processed",
	})

	PaymentFailuresTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "payment_failures_total",
		Help: "The total number of payment failures",
	})

	LedgerTransactionsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ledger_transactions_total",
		Help: "The total number of ledger transactions written",
	})
)
