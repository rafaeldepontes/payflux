package app

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rafaeldepontes/reconsiliation/internal/reconciliation"
	rm "github.com/rafaeldepontes/reconsiliation/internal/reconciliation/model"
	rr "github.com/rafaeldepontes/reconsiliation/internal/reconciliation/repository"
	rs "github.com/rafaeldepontes/reconsiliation/internal/reconciliation/service"
	"github.com/rafaeldepontes/reconsiliation/internal/risk"
	rk "github.com/rafaeldepontes/reconsiliation/internal/risk/repository"
	risks "github.com/rafaeldepontes/reconsiliation/internal/risk/service"
	"github.com/rafaeldepontes/reconsiliation/pkg/message-broker/rabbitmq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type App struct {
	reconSvc reconciliation.Service
	riskSvc  risk.Service
	tr       trace.Tracer
}

func New() *App {
	reconRepo := rr.NewRepository()
	riskRepo := rk.NewRepository()

	return &App{
		reconSvc: rs.NewService(reconRepo),
		riskSvc:  risks.NewService(riskRepo),
		tr:       otel.Tracer("reconsiliation-consumer"),
	}
}

func (a *App) Run() {
	msgs := rabbitmq.GetConsumer()
	if msgs == nil {
		log.Fatal("could not get rabbitmq consumer")
	}

	dlqMsgs := rabbitmq.GetDLQConsumer()
	if dlqMsgs == nil {
		log.Fatal("could not get rabbitmq DLQ consumer")
	}

	go a.consumeMain(*msgs)
	go a.consumeDLQ(*dlqMsgs)
}

func (a *App) consumeMain(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		if err := a.handlePaymentEvent(d); err != nil {
			log.Printf("[ERROR] main consumer failed: %v", err)
			_ = d.Nack(false, false)
			continue
		}
		_ = d.Ack(false)
	}
}

func (a *App) consumeDLQ(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		log.Printf("[DLQ] dead letter received: %s", string(d.Body))

		// Not sure what I have to do here...
		// Just inspect / log / store / alert??
		_ = d.Ack(false)
	}
}

func (a *App) handlePaymentEvent(d amqp.Delivery) error {
	_, span := a.tr.Start(context.Background(), "consume-event")
	defer span.End()

	var event rm.PaymentEvent
	if err := json.Unmarshal(d.Body, &event); err != nil {
		return err
	}

	log.Printf("[INFO] Received event: %s for payment: %s", event.EventType, event.PaymentID)

	if err := a.reconSvc.ProcessEvent(event); err != nil {
		return err
	}
	if err := a.riskSvc.ProcessEvent(event); err != nil {
		return err
	}

	return nil
}
