package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/goplo/internal/handler"
	"github.com/rafaeldepontes/goplo/pkg/cache"
	"github.com/rafaeldepontes/goplo/pkg/db/postgres"
)

func init() {
	_ = godotenv.Load(".env", ".env.example")
	cache.GetCache()
	postgres.GetDb()
	if err := postgres.RunMigrations(); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
}

func main() {
	defer cache.Close()
	defer postgres.Close()

	port := os.Getenv("API_PORT")

	h := handler.NewHandler()

	log.Println("Application running on localhost:" + port)
	log.Fatalln(http.ListenAndServe(":"+port, h))
}
