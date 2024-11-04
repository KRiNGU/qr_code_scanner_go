CREATE TABLE IF NOT EXISTS receipts (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    price NUMERIC(12, 2),
    amount NUMERIC(10, 3),
    receipt_id integer REFERENCES receipts,
    product_id integer REFERENCES products,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);
