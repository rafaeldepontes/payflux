# PayFlux - Fintech Payment & Reconciliation System

PayFlux is a production-style fintech backend demonstration consisting of two microservices written in Go. It simulates a complete payment lifecycle: from request and idempotent processing to double-entry accounting, event-driven risk analysis, and transaction reconciliation.

## High-Level Architecture

```text
            Client
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
              |
              | (Postgres: Reconciliation & Risk Tables)
              | (External: Settlement Records)
```

## Getting Started

### Prerequisites
- Docker & Docker Compose
- Go 1.26+ (for local development)

### Running the System
Run the entire stack using the root `docker-compose.yml`:

```bash
docker compose up --build
```

This will start:
- **Ledger API** (`:8080`)
- **Reconciliation API** (`:8081`)
- **PostgreSQL** (Shared database with isolated schemas)
- **RabbitMQ** (Message Broker)
- **Redis** (Cache/Idempotency)

---

## Service Breakdown

### 1. Ledger Service
The core engine for moving money.
- **Double-Entry Accounting:** Every payment creates a debit and a credit entry, ensuring `Total Debits = Total Credits`.
- **Idempotency:** Requires an `Idempotency-Key` header to prevent duplicate charges.
- **Event Producer:** Emits `PaymentCompleted` and `PaymentRefunded` events to RabbitMQ.

### 2. Reconciliation & Risk Service
The watchdog for financial integrity.
- **Event Consumer:** Listens to ledger events.
- **Risk Engine:** Evaluates transactions against rules (e.g., `LargeTransactionRule`).
- **Reconciliation:** Compares internal ledger events with external "Settlement Records" (simulating bank statements).

---

## Database Architecture (PostgreSQL)

The services share a Postgres instance but use unique migration tables to manage their schemas:

### Ledger Schema
- `accounts`: Stores user and merchant wallet info.
- `payments`: Tracks the state of payment requests.
- `ledger_entries`: The immutable audit log of all money movements.

### Reconciliation Schema
- `settlement_records`: External data to be matched against our ledger.
- `reconciliation_results`: Stores the outcome (Matched/Mismatched) of comparisons.
- `risk_evaluations`: Stores risk scores and flags for each transaction.
- `exceptions`: Log of any discrepancies found.

---

## Example Use Flow

### 1. Create a Payment
Send a request to the Ledger service.
```bash
curl -X POST http://localhost:8080/payments \
  -H "Idempotency-Key: unique-key-12345" \
  -H "Content-Type: application/json" \
  -d '{
    "from_account": 1,
    "to_account": 2,
    "amount": 500,
    "currency": "USD"
  }'
```

### 2. Simulate External Settlement
The Reconciliation service needs "External" data to match. Use the helper endpoint:
```bash
curl -X POST http://localhost:8081/settlements \
  -H "Content-Type: application/json" \
  -d '{
    "transaction_id": "PASTE_PAYMENT_ID_HERE",
    "amount": 500,
    "status": "Settled"
  }'
```

### 3. Check Reconciliation & Risk
Verify how the system processed the event:
```bash
# Check Reconciliation Result
curl http://localhost:8081/reconciliation/PASTE_PAYMENT_ID_HERE

# Check Risk Score
curl http://localhost:8081/risk/PASTE_PAYMENT_ID_HERE
```

### 4. Check Balance
Verify the money moved correctly in the ledger:
```bash
curl http://localhost:8080/accounts/1/balance
```

## License
MIT
