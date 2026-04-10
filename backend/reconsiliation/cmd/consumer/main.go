package main

import (
	"context"
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
	"github.com/rafaeldepontes/reconsiliation/pkg/observability"
	"go.opentelemetry.io/otel"
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

	reconSvc := rs.NewService()
	riskSvc := risks.NewService()

	msgs := rabbitmq.GetConsumer()
	if msgs == nil {
		log.Fatal("could not get rabbitmq consumer")
	}

	tr := otel.Tracer("reconsiliation-consumer")

	go func() {
		for d := range *msgs {
			_, span := tr.Start(context.Background(), "consume-event")

			var event rm.PaymentEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("[ERROR] failed to unmarshal event: %v", err)
				span.End()
				continue
			}

			log.Printf("[INFO] Received event: %s for payment: %s", event.EventType, event.PaymentID)

			if err := reconSvc.ProcessEvent(event); err != nil {
				log.Printf("[ERROR] reconciliation failed: %v", err)
			}

			if err := riskSvc.ProcessEvent(event); err != nil {
				log.Printf("[ERROR] risk evaluation failed: %v", err)
			}
			span.End()
		}
	}()

	log.Println("Consumer started. Waiting for messages...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down consumer...")
}
