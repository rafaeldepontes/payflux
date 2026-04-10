package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ReconciliationProcessedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "reconciliation_processed_total",
		Help: "The total number of reconciliation requests processed",
	})

	ReconciliationFailuresTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "reconciliation_failures_total",
		Help: "The total number of reconciliation failures",
	})

	RiskFlagsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "risk_flags_total",
		Help: "The total number of risk flags generated",
	})
)
