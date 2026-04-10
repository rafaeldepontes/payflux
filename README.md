# PayFlux - Fintech Payment & Reconciliation System

PayFlux is a production-style fintech backend demonstration consisting of microservices written in Go and React frontend. It simulates a complete payment lifecycle: from request and idempotent processing to double-entry accounting, event-driven risk analysis, and transaction reconciliation.

## High-Level Architecture

```text
            Client (React Frontend)
              |
              v
      +-------+-------+
      |     Ledger    | <--- (API: Payments, Refunds, Balances)
      |    Service    |
      +-------+-------+
              |
              | (Postgres: Double-Entry Ledger)
              | (Redis: Idempotency)
              |
              v (RabbitMQ Event: PaymentCompleted)
              |
      +-------+-------+
      | Reconciliation| <--- (API: Risk Scores, Exceptions)
      |  & Risk Svc   |
      +-------+-------+
```

## Getting Started

### Prerequisites
- Docker & Docker Compose
- Node.js
- Go

### Running the System
Run the entire stack using the root `docker-compose.yml`:

```bash
docker compose up --build
```

This will start:
- **Frontend App** (`:3000`)
- **Ledger API** (`:8080`)
- **Reconciliation API** (`:8081`)
- **Prometheus** (`:9090`) - Metrics dashboard
- **OpenTelemetry Collector** (`:4317`, `:4318`) - Trace collection
- **PostgreSQL** (Shared database with isolated schemas)
- **RabbitMQ** (Message Broker)
- **Redis** (Cache/Idempotency)

---

## Observability & Documentation

### Metrics & Traces
- **Prometheus UI:** `http://localhost:9090`
- **Ledger Metrics:** `http://localhost:8080/metrics`
- **Reconciliation Metrics:** `http://localhost:8081/metrics`
- **Tracing:** Exported via OTLP to the OpenTelemetry Collector.

### API Documentation (Swagger)
- **Ledger API Docs:** `http://localhost:8080/swagger/`
- **Reconciliation API Docs:** `http://localhost:8081/swagger/`

---

## Service Breakdown

### 1. Frontend (React + TypeScript + Vite)
- A modern dashboard to interact with all backend features.
- Create payments, check status, and view balances.
- Styled with Tailwind CSS.

### 2. Ledger Service (Go)
The core engine for moving money.
- **Double-Entry Accounting:** Ensures `Total Debits = Total Credits`.
- **Idempotency:** Prevents duplicate charges via `Idempotency-Key`.

### 3. Reconciliation & Risk Service (Go)
The watchdog for financial integrity.
- **Risk Engine:** Automated rule evaluation for every transaction.
- **Reconciliation:** Matches internal events with external settlement records.

---

## Database Architecture (PostgreSQL)

The services share a Postgres instance but use unique migration tables:
- `ledger_schema_migrations`: Tracks ledger-specific tables (`accounts`, `payments`, `ledger_entries`).
- `reconciliation_schema_migrations`: Tracks reconciliation tables (`settlement_records`, `reconciliation_results`, `risk_evaluations`, `exceptions`).

---

## Example Use Flow

1. **Access the Frontend:** Open `http://localhost:3000`.
2. **Create a Payment:** Fill the "Create Payment" form and click "Send".
3. **Simulate Settlement:** After a successful payment, click "Simulate Settlement" to create the external record.
4. **Verify:** Use the "Check Transaction" section to see the Reconciliation (should be `matched`) and Risk Score.
5. **Check Balance:** Enter the Account ID (1 or 2) in the "Account Balance" section to see the updated funds.

## License
MIT
