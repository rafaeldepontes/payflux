# Golang - Payment Ledger & Orchestrator

Project built with [Gini](https://gini-webserver.up.railway.app/)

---

A **production-style fintech backend service written in Go** that
simulates a payment processing system with a **double-entry ledger**,
**idempotent payment execution**, and **event-driven architecture**.

This project is designed for **engineering evaluation purposes**
(portfolio / technical review).\
It demonstrates how financial systems guarantee **consistency,
auditability, and reliability** when moving money.

------------------------------------------------------------------------

# Overview

The system receives payment requests through an API and processes them
using an orchestrator service.\
Each payment produces immutable ledger entries and emits events to a
message broker.

Key goals of the project:

-   Guarantee **double-entry accounting correctness**
-   Provide **idempotent payment processing**
-   Demonstrate **event-driven architecture**
-   Maintain a fully **auditable ledger**
-   Show **production-style backend structure**

------------------------------------------------------------------------

# Tech Stack

-   **Language:** Go
-   **Database:** PostgreSQL
-   **Message Broker:** RabbitMQ
-   **Cache / Idempotency store:** Redis
-   **Containerization:** Docker
-   **Observability:** Prometheus + OpenTelemetry

------------------------------------------------------------------------

# High-Level Architecture

            Client
              |
              v
          API Gateway
              |
              v
       Payment Orchestrator
          |           |
          |           v
          |      Redis (Idempotency)
          |
          v
       PostgreSQL (Ledger)
          |
          v
      RabbitMQ (Events)
          |
          v
     Downstream Consumers
    (Notifications, Risk, Analytics)

Flow:

1.  Client sends payment request
2.  Orchestrator validates and checks idempotency
3.  Ledger entries are written in a DB transaction
4.  Event is emitted to RabbitMQ
5.  Downstream services consume events

------------------------------------------------------------------------

# Ledger Model

The ledger follows **double-entry accounting**.

Example:

    User Wallet      -100
    Merchant Wallet  +100

Every transaction must satisfy:

    Total Debits = Total Credits

Ledger entries are **immutable** and cannot be updated.

------------------------------------------------------------------------

# API Endpoints

## Create Payment

Creates a new payment transaction.

    POST /payments

Request:

> Headers: ... | Idempotency-Key: "abc-123" | ....

``` json
{
  "from_account": "user_wallet",
  "to_account": "merchant_wallet",
  "amount": 100,
  "currency": "USD",
}
```

Response:

``` json
{
  "payment_id": "pay_9f21ab",
  "status": "processed"
}
```

------------------------------------------------------------------------

## Get Payment

Returns payment information.

    GET /payments/{id}

Response

``` json
{
  "payment_id": "pay_9f21ab",
  "status": "processed",
  "amount": 100,
  "currency": "USD"
}
```

------------------------------------------------------------------------

## Refund Payment

Creates a refund transaction.

    POST /payments/{id}/refund

Request

``` json
{
  "amount": 100
}
```

Response

``` json
{
  "refund_id": "ref_88412",
  "status": "processed"
}
```

------------------------------------------------------------------------

## Get Account Balance

Returns the computed account balance.

    GET /accounts/{id}/balance

Response

``` json
{
  "account_id": "merchant_wallet",
  "balance": 1500
}
```

------------------------------------------------------------------------

# Event Model

Each payment emits events to RabbitMQ.

Example event:

    PaymentCreated
    PaymentCompleted
    PaymentRefunded

Example message:

``` json
{
  "event_type": "PaymentCompleted",
  "payment_id": "pay_9f21ab",
  "amount": 100,
  "currency": "USD",
  "timestamp": "2026-04-05T12:00:00Z"
}
```

------------------------------------------------------------------------

# Idempotency

To prevent duplicate payments the API requires an `idempotency_key`.

If the same key is reused:

-   the original response is returned
-   no duplicate transaction occurs

Implementation:

    Client Request
          |
          v
    Redis Idempotency Check
          |
          v
    Process Payment (if new)

------------------------------------------------------------------------

# Running the Project

Requirements

-   Docker
-   Docker Compose
-   Go 1.26+

Run:

    docker compose up --build

Services started:

-   API server
-   PostgreSQL
-   RabbitMQ
-   Redis

------------------------------------------------------------------------

# Repository Structure

    internal/
        util/
        handler/
        payment/
        ledger/
        account/
        idempotency/
    pkg/
        events/
        rabbitmq/
        observability/
    deploy/
        docker/
        docker-compose.yml

------------------------------------------------------------------------

# Observability

The service exposes:

Metrics

    /metrics

Prometheus examples:

-   payment_requests_total
-   payment_failures_total
-   ledger_transactions_total

Tracing

OpenTelemetry traces:

    API Request → Orchestrator → Ledger → RabbitMQ

------------------------------------------------------------------------

# Testing

Run tests:

    go test ./...

Integration tests simulate:

-   payment processing
-   refund flows
-   idempotency checks
-   ledger correctness

------------------------------------------------------------------------

# Design Goals

This project demonstrates:

-   financial correctness
-   event-driven microservice architecture
-   production-style Go backend patterns
-   distributed system design for fintech environments

------------------------------------------------------------------------

# License

MIT
