ALTER TABLE products ADD COLUMN id SERIAL;
ALTER TABLE products DROP CONSTRAINT products_pkey;
ALTER TABLE products ADD PRIMARY KEY (id);

ALTER TABLE transactions ADD COLUMN product_id INTEGER;
UPDATE transactions SET product_id = (SELECT id FROM products WHERE products.product_name = transactions.product_name_fk);

ALTER TABLE transactions DROP CONSTRAINT transactions_product_name_fk_fkey;
ALTER TABLE transactions ADD CONSTRAINT transactions_product_id_fkey FOREIGN KEY (product_id) REFERENCES transactions(id);
ALTER TABLE transactions DROP COLUMN product_name_fk;
