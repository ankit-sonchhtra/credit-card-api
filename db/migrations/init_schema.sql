CREATE TABLE accounts
(
    account_id      BIGSERIAL PRIMARY KEY,
    document_number VARCHAR(20) NOT NULL UNIQUE,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE operation_types
(
    operation_type_id INT PRIMARY KEY,
    description       VARCHAR(50) NOT NULL
);

CREATE TABLE transactions
(
    transaction_id    BIGSERIAL PRIMARY KEY,
    account_id        BIGINT            NOT NULL REFERENCES accounts (account_id),
    operation_type_id BIGINT            NOT NULL REFERENCES operation_types (operation_type_id),
    amount            NUMERIC(15, 2) NOT NULL,
    created_at        TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

-- Seed data
INSERT INTO operation_types (operation_type_id, description)
VALUES (1, 'Normal Purchase'),
       (2, 'Purchase with installments'),
       (3, 'Withdrawal'),
       (4, 'Credit Voucher');