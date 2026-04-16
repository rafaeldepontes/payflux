package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/reconsiliation/internal/app"
	"github.com/rafaeldepontes/reconsiliation/pkg/db/postgres"
	"github.com/rafaeldepontes/reconsiliation/pkg/observability"
)

func init() {
	_ = godotenv.Load(".env", ".env.example")
}

func main() {
	defer postgres.Close()

	tp, err := observability.InitTracer("reconsiliation-consumer")
	if err != nil {
		log.Printf("failed to init tracer: %v", err)
	} else {
		defer tp.Shutdown(context.Background())
	}

	if err := postgres.RunMigrations(); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	app := app.New()
	app.Run()

	log.Println("Consumer started. Waiting for messages...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down consumer...")
}
