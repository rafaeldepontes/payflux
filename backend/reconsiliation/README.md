# Golang - Reconciliation & Risk Engine

Project built with [Gini](https://gini-webserver.up.railway.app/)

---

A **production-style fintech backend service written in Go** that
performs **transaction reconciliation and basic risk analysis** on
financial events.

This project is designed for **engineering evaluation and portfolio
review**.\
It demonstrates how financial systems validate ledger activity, detect
inconsistencies, and apply automated risk rules.

------------------------------------------------------------------------

# Overview

The system consumes financial events (payments, refunds, ledger entries)
and compares them against external settlement records or simulated bank
statements.

Key objectives:

-   Detect **transaction mismatches**
-   Perform **automated reconciliation**
-   Apply **risk rules and anomaly detection**
-   Produce **exception reports**
-   Demonstrate **event-driven processing in Go**

------------------------------------------------------------------------

# Tech Stack

-   **Language:** Go
-   **Database:** PostgreSQL
-   **Message Broker:** RabbitMQ
-   **Cache:** Redis
-   **Observability:** Prometheus + OpenTelemetry
-   **Containerization:** Docker

------------------------------------------------------------------------

# High-Level Architecture

                RabbitMQ
                  |
                  v
         Reconciliation Consumer
                  |
          ----------------------
          |                    |
          v                    v
     Reconciliation Engine   Risk Engine
          |                    |
          -----------+---------
                     |
                     v
                PostgreSQL
                     |
                     v
               Exception Events
                     |
                     v
            Alerting / Reporting

Flow:

1.  Payment or ledger event is emitted to RabbitMQ
2.  Reconciliation service consumes the event
3.  Transaction is matched with settlement records
4.  Risk engine evaluates predefined rules
5.  Exceptions are stored and emitted as alerts

------------------------------------------------------------------------

# Reconciliation Model

Reconciliation ensures the internal ledger matches external records.

Example comparison:

Internal Ledger:

    Transaction ID: tx_1001
    Amount: 100 USD
    Status: Completed

External Settlement:

    Transaction ID: tx_1001
    Amount: 100 USD
    Status: Settled

If any mismatch occurs:

-   status mismatch
-   amount mismatch
-   missing transaction

The system generates an **exception record**.

------------------------------------------------------------------------

# Risk Evaluation

The risk engine evaluates transactions against simple rules.

Example rules:

    LargeTransactionRule
    DuplicateTransactionRule
    VelocityRule

Example:

    If transaction amount > 10,000 USD → flag as high risk
    If same user performs > 10 payments in 60 seconds → flag velocity risk

Risk scores are attached to each transaction.

------------------------------------------------------------------------

# API Endpoints

Although most processing is event-driven, the service exposes endpoints
for querying results.

------------------------------------------------------------------------

## Get Reconciliation Result

    GET /reconciliation/{transaction_id}

Response

``` json
{
  "transaction_id": "tx_1001",
  "status": "matched",
  "ledger_amount": 100,
  "settlement_amount": 100
}
```

------------------------------------------------------------------------

## Create Settlement Record

    POST /settlements

Request

``` json
{
  "transaction_id": "tx_1001",
  "amount": 100,
  "status": "Settled"
}
```

------------------------------------------------------------------------

## Get Risk Evaluation

    GET /risk/{transaction_id}

Response

``` json
{
  "transaction_id": "tx_1001",
  "risk_score": 12,
  "flags": [
    "LargeTransactionRule"
  ]
}
```

------------------------------------------------------------------------

## List Exceptions

Returns all detected reconciliation mismatches.

    GET /exceptions

Response

``` json
[
  {
    "transaction_id": "tx_2042",
    "type": "AmountMismatch",
    "ledger_amount": 120,
    "settlement_amount": 100
  }
]
```

------------------------------------------------------------------------

# Event Consumption

The service subscribes to RabbitMQ topics:

    payments.created
    payments.completed
    payments.refunded
    ledger.entries

Example event:

``` json
{
  "event_type": "PaymentCompleted",
  "transaction_id": "tx_1001",
  "amount": 100,
  "currency": "USD",
  "timestamp": "2026-04-05T12:00:00Z"
}
```

------------------------------------------------------------------------

# Exception Handling

When reconciliation fails:

1.  Exception stored in PostgreSQL
2.  Alert event emitted
3.  Transaction flagged for manual review

Example:

    ReconciliationFailed
    RiskFlagged
    ManualReviewRequired

------------------------------------------------------------------------

# Running the Project

Requirements

-   Docker
-   Docker Compose
-   Go 1.22+

Run:

    docker compose up --build

Services started:

-   Reconciliation service
-   PostgreSQL
-   RabbitMQ
-   Redis

------------------------------------------------------------------------

# Repository Structure

    cmd/
        api/
        consumer/
    internal/
        reconciliation/
        risk/
        rules/
        exceptions/
    pkg/
        events/
        rabbitmq/
        observability/
    deploy/
        docker/
        docker-compose.yml

------------------------------------------------------------------------

# Observability

Metrics exposed:

    /metrics

Prometheus metrics examples:

-   reconciliation_processed_total
-   reconciliation_failures_total
-   risk_flags_total

Tracing:

    RabbitMQ Consumer → Reconciliation Engine → Risk Engine → Database

------------------------------------------------------------------------

# Testing

Run tests:

    go test ./...

Test scenarios include:

-   matched transactions
-   missing settlement records
-   amount mismatches
-   risk rule triggering

------------------------------------------------------------------------

# Design Goals

This project demonstrates:

-   event-driven backend architecture
-   financial reconciliation logic
-   automated risk evaluation
-   observability and reliability patterns used in fintech systems

------------------------------------------------------------------------

# License

MIT

