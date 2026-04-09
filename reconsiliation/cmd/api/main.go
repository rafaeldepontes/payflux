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

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Reconciliation API running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, h))
}
