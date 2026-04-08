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
}

func main() {
	defer cache.Close()
	defer postgres.Close()

	h := handler.NewHandler()

	log.Fatalln(http.ListenAndServe(":"+os.Getenv("API_PORT"), h))
}
