package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	rm "github.com/rafaeldepontes/reconsiliation/internal/reconciliation/model"
	rs "github.com/rafaeldepontes/reconsiliation/internal/reconciliation/service"
	risks "github.com/rafaeldepontes/reconsiliation/internal/risk/service"
	"github.com/rafaeldepontes/reconsiliation/pkg/db/postgres"
	"github.com/rafaeldepontes/reconsiliation/pkg/message-broker/rabbitmq"
)

func init() {
	_ = godotenv.Load(".env", ".env.example")
}

func main() {
	defer postgres.Close()

	if err := postgres.RunMigrations(); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	reconSvc := rs.NewService()
	riskSvc := risks.NewService()

	msgs := rabbitmq.GetConsumer()
	if msgs == nil {
		log.Fatal("could not get rabbitmq consumer")
	}

	go func() {
		for d := range *msgs {
			var event rm.PaymentEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("[ERROR] failed to unmarshal event: %v", err)
				continue
			}

			log.Printf("[INFO] Received event: %s for payment: %s", event.EventType, event.PaymentID)

			if err := reconSvc.ProcessEvent(event); err != nil {
				log.Printf("[ERROR] reconciliation failed: %v", err)
			}

			if err := riskSvc.ProcessEvent(event); err != nil {
				log.Printf("[ERROR] risk evaluation failed: %v", err)
			}
		}
	}()

	log.Println("Consumer started. Waiting for messages...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down consumer...")
}
