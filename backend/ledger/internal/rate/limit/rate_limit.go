package limit

import (
	"log"
	"net/http"

	"github.com/rafaeldepontes/ledger/internal/util"
	"golang.org/x/time/rate"
)

type Middleware interface {
	RateLimit(h http.Handler) http.Handler
}

type middleware struct {
	rl *rate.Limiter
}

func NewMiddleware() Middleware {
	return &middleware{
		rl: rate.NewLimiter(1, 5),
	}
}

func (m *middleware) RateLimit(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.rl.Allow() {
			log.Println("[WARN] Request limit reached.")
			util.HandleError(w, "too many requests.", http.StatusTooManyRequests)
			return
		}
		h.ServeHTTP(w, r)
	})
}
