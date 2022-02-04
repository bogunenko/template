
/* name: CreateAccount :execresult */
INSERT INTO account () VALUES ();

/* name: GetAccount :one */
SELECT * FROM account
WHERE id = ?;

/* name: GetAccountForUpdate :one */
SELECT * FROM account
WHERE id = ? FOR UPDATE;

/* name: Deposit :execresult */
UPDATE account
SET balance = balance + sqlc.arg(amount)
WHERE id = ?;

/* name: Withdraw :execresult */
UPDATE account
SET balance = balance - sqlc.arg(amount)
WHERE id = ?;

/* name: ListAccounts :many */
SELECT * FROM account;

/* name: CreateTransaction :execresult */
INSERT INTO transaction (from_id, to_id, amount)
VALUES (?, ?, ?);

/* name: GetTransaction :one */
SELECT * FROM transaction
WHERE id = ?;