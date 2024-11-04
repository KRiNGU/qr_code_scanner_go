ALTER TABLE transactions ADD COLUMN product_name_fk VARCHAR(255);
UPDATE transactions SET product_name_fk = (SELECT product_name FROM products WHERE products.id = transactions.product_id);

ALTER TABLE transactions DROP CONSTRAINT transactions_product_id_fkey;
ALTER TABLE transactions ADD CONSTRAINT transactions_product_name_fk_fkey FOREIGN KEY (product_name_fk) REFERENCES products(product_name);
ALTER TABLE transactions DROP COLUMN product_id;

ALTER TABLE products DROP CONSTRAINT products_pkey;
ALTER TABLE products DROP id;
ALTER TABLE products ADD PRIMARY KEY (product_name);
