CREATE TABLE IF NOT EXISTS transactions (
    transaction_id BIGSERIAL PRIMARY KEY,
    source_account_id BIGINT NOT NULL,
    destination_account_id BIGINT NOT NULL,
    amount NUMERIC(20,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_source_account FOREIGN KEY (source_account_id) REFERENCES accounts(account_id),
    CONSTRAINT fk_destination_account FOREIGN KEY (destination_account_id) REFERENCES accounts(account_id)
);