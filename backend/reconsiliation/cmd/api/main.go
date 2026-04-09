package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/reconsiliation/internal/handler"
	"github.com/rafaeldepontes/reconsiliation/pkg/db/postgres"
)

func init() {
	_ = godotenv.Load(".env", ".env.example")
}

func main() {
	defer postgres.Close()

	h := handler.NewHandler()

	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Reconciliation API running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
