package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/reconsiliation/internal/handler"
	"github.com/rafaeldepontes/reconsiliation/pkg/db/postgres"
	"github.com/rafaeldepontes/reconsiliation/pkg/observability"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func init() {
	_ = godotenv.Load(".env", ".env.example")
}

// @title Reconciliation & Risk API
// @version 1.0
// @description A production-style fintech backend service for transaction reconciliation and risk analysis.
// @host localhost:8081
// @BasePath /
func main() {
	defer postgres.Close()

	tp, err := observability.InitTracer("reconsiliation-api")
	if err != nil {
		log.Printf("failed to init tracer: %v", err)
	} else {
		defer tp.Shutdown(context.Background())
	}

	h := handler.NewHandler()

	otelHandler := otelhttp.NewHandler(h, "reconsiliation-api")

	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		otelHandler.ServeHTTP(w, r)
	})

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Reconciliation API running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
