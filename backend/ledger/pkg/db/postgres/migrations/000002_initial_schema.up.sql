CREATE TABLE account_balance (
    id BIGSERIAL PRIMARY KEY,
    account_id INT NOT NULL REFERENCES accounts(id),
    balance BIGINT NOT NULL,
    version INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (account_id, version)
);

CREATE INDEX idx_account_balance_account_id
ON account_balance (account_id, created_at DESC);