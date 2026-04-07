-- name: CreateTransaction :one
INSERT INTO transactions (account_id, operation_type_id, amount, balance)
VALUES ($1, $2, $3, $4)
    RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE transaction_id = $1 LIMIT 1;

-- name: GetAllTransactionById :many
SELECT *
FROM transactions
WHERE account_id = $1
ORDER BY created_at ASC;


-- name: UpdateTransaction :one
UPDATE transactions
SET balance = $2
WHERE transaction_id = $1
    RETURNING *;

-- name: ListTransactionsByAccount :many
SELECT * FROM transactions
WHERE account_id = $1
ORDER BY created_at DESC;