-- Create payments first
WITH seeded_payments AS (
    INSERT INTO payments (
        id,
        idempotency_key,
        from_account_id,
        to_account_id,
        amount,
        currency,
        status
    )
    VALUES
    (
        gen_random_uuid(),
        'seed-user-1',
        1,
        1,
        100000,
        'USD',
        'completed'
    ),
    (
        gen_random_uuid(),
        'seed-user-2',
        2,
        2,
        250000,
        'USD',
        'completed'
    )
    RETURNING id, to_account_id, amount, currency
)

-- Create ledger entries
INSERT INTO ledger_entries (
    payment_id,
    account_id,
    amount,
    currency
)
SELECT
    id,
    to_account_id,
    amount,
    currency
FROM seeded_payments;