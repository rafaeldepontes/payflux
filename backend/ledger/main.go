package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	ar "github.com/rafaeldepontes/ledger/internal/account/repository"
	as "github.com/rafaeldepontes/ledger/internal/account/server"
	asvc "github.com/rafaeldepontes/ledger/internal/account/service"
	cs "github.com/rafaeldepontes/ledger/internal/cache/service"
	"github.com/rafaeldepontes/ledger/internal/handler"
	"github.com/rafaeldepontes/ledger/internal/rate/limit"
	pr "github.com/rafaeldepontes/ledger/internal/payment/repository"
	ps "github.com/rafaeldepontes/ledger/internal/payment/server"
	psvc "github.com/rafaeldepontes/ledger/internal/payment/service"
	"github.com/rafaeldepontes/ledger/pkg/cache"
	"github.com/rafaeldepontes/ledger/pkg/db/postgres"
	"github.com/rafaeldepontes/ledger/pkg/message-broker/rabbitmq"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func init() {
	_ = godotenv.Load(".env", ".env.example")
}

// @title Payment Ledger API
// @version 1.0
// @description A production-style fintech backend service for payment processing.
// @host localhost:8080
// @BasePath /
func main() {
	postgres.GetDb()
	defer postgres.Close()

	if err := postgres.RunMigrations(); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	redisClient := cache.GetCache()
	defer cache.Close()

	rabbitmq.GetConnection()
	rabbitmq.GetChannel()
	rabbitmq.GetQueue()
	defer rabbitmq.Close()

	accountRepo := ar.NewRepository()
	accountSvc := asvc.NewService(accountRepo)
	accountCtrl := as.NewController(accountSvc)

	cacheSvc := cs.NewService(redisClient)
	paymentRepo := pr.NewRepository()
	broker := &rabbitmq.Broker{}
	paymentSvc := psvc.NewService(paymentRepo, cacheSvc, broker)
	paymentCtrl := ps.NewController(paymentSvc)

	port := os.Getenv("API_PORT")

	rateLimit := limit.NewMiddleware()
	h := handler.NewHandler(paymentCtrl, accountCtrl, rateLimit)

	otelHandler := otelhttp.NewHandler(h, "ledger-api")

	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Idempotency-Key")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		otelHandler.ServeHTTP(w, r)
	})

	log.Println("Application running on localhost:" + port)
	log.Fatalln(http.ListenAndServe(":"+port, corsHandler))
}
