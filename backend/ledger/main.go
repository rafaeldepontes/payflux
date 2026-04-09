package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/ledger/internal/handler"
	"github.com/rafaeldepontes/ledger/pkg/cache"
	"github.com/rafaeldepontes/ledger/pkg/db/postgres"
	"github.com/rafaeldepontes/ledger/pkg/message-broker/rabbitmq"
)

func init() {
	_ = godotenv.Load(".env", ".env.example")
	cache.GetCache()
	postgres.GetDb()
	rabbitmq.GetConnection()
	rabbitmq.GetChannel()
	rabbitmq.GetQueue()
	if err := postgres.RunMigrations(); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
}

func main() {
	defer cache.Close()
	defer postgres.Close()
	defer rabbitmq.Close()

	port := os.Getenv("API_PORT")

	h := handler.NewHandler()

	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Idempotency-Key")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})

	log.Println("Application running on localhost:" + port)
	log.Fatalln(http.ListenAndServe(":"+port, corsHandler))
}
