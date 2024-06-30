package postgres

const (
	userBalanceQuery = `SELECT 
    (SELECT COALESCE(SUM(accrual), 0) FROM orders WHERE user_id = $1) AS balance,
    (SELECT COALESCE(SUM(sum), 0) FROM withdrawals WHERE user_id = $1) AS total_withdrawals;
`

	userWithdrawalsQuery = `SELECT
    w.order_id,
    w.sum,
    w.processed_at
FROM
    withdrawals w
WHERE
    w.user_id = $1;
`

	checkWithdrawQuery = `SELECT EXISTS (
    SELECT 1
    FROM withdrawals
    WHERE order_id = $1
) AS withdrawal_exists;
`

	processWithdrawQuery = `INSERT INTO withdrawals (order_id, user_id, sum, processed_at)
VALUES ($1, (SELECT user_id FROM orders WHERE id = $1), $2, CURRENT_TIMESTAMP);
`

	checkOrderQuery = `SELECT user_id
FROM orders
WHERE id = $1;
`

	addOrderQuery = `INSERT INTO orders (id, user_id, status, uploaded_at)
VALUES ($1, $2, $3, CURRENT_TIMESTAMP);
`

	getNotUpdatedOrdersQuery = `SELECT id
FROM orders
WHERE status IN ('NEW', 'PROCESSING');
`

	updateOrderQuery = `UPDATE orders
SET status = $2,
    accrual = $3
WHERE id = $1;
`

	getUserOrdersQuery = `SELECT id, status, accrual, uploaded_at
FROM orders
WHERE user_id = $1;
`

	checkUserQuery = `SELECT id
FROM users
WHERE username = $1
  AND pass_hash = $2;
`

	addUserQuery = `INSERT INTO users (username, pass_hash)
VALUES ($1, $2)
RETURNING id;
`

	checkCreateTablesQuery = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY NOT NULL,
    username VARCHAR UNIQUE NOT NULL,
    pass_hash VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
    id INTEGER PRIMARY KEY NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id),
    status VARCHAR NOT NULL,
    accrual DOUBLE PRECISION,
    uploaded_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS withdrawls (
    id SERIAL PRIMARY KEY NOT NULL,
    order_id BIGINT NOT NULL REFERENCES orders(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    sum DOUBLE PRECISION NOT NULL
);
`
)
