package postgres

const (
	userBalanceQuery = `SELECT 
    COALESCE(SUM(o.accrual), 0) - COALESCE(SUM(w.sum), 0) AS balance,
    COALESCE(SUM(w.sum), 0) AS total_withdrawals
FROM 
    users u
LEFT JOIN 
    orders o ON u.id = o.user_id
LEFT JOIN 
    withdrawals w ON u.id = w.user_id
WHERE 
    u.id = $1
GROUP BY 
    u.id;
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
VALUES ($1, $2, $3, CURRENT_TIMESTAMP);
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
SET status = $1,
    accrual = $2
WHERE id = $3;
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

	createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY NOT NULL,
    username VARCHAR UNIQUE NOT NULL,
    pass_hash VARCHAR NOT NULL
);`
	createOrdersTable = `
CREATE TABLE IF NOT EXISTS orders (
    id BIGINT PRIMARY KEY NOT NULL,
    user_id BIGINT NOT NULL,
    status VARCHAR NOT NULL,
    accrual DOUBLE PRECISION,
    uploaded_at TIMESTAMP NOT NULL
);`
	createWithdrawlsTable = `
CREATE TABLE IF NOT EXISTS withdrawals (
    id SERIAL PRIMARY KEY NOT NULL,
    order_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    sum DOUBLE PRECISION NOT NULL,
    processed_at TIMESTAMP NOT NULL
);`
)
