package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/goplo/internal/handler"
	"github.com/rafaeldepontes/goplo/pkg/cache"
	"github.com/rafaeldepontes/goplo/pkg/db/postgres"
	"github.com/rafaeldepontes/goplo/pkg/message-broker/rabbitmq"
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

	log.Println("Application running on localhost:" + port)
	log.Fatalln(http.ListenAndServe(":"+port, h))
}
