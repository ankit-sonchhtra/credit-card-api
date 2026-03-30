-- name: CreateTransaction :one
INSERT INTO transactions (account_id, operation_type_id, amount)
VALUES ($1, $2, $3)
    RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE transaction_id = $1 LIMIT 1;

-- name: ListTransactionsByAccount :many
SELECT * FROM transactions
WHERE account_id = $1
ORDER BY created_at DESC;