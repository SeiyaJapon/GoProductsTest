CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    code VARCHAR(32) UNIQUE NOT NULL,
    name VARCHAR(64) NOT NULL
);

INSERT INTO categories (code, name) VALUES
    ('Clothing', 'Clothing'),
    ('Shoes', 'Shoes'),
    ('Accessories', 'Accessories');

ALTER TABLE products ADD COLUMN IF NOT EXISTS category_id INTEGER REFERENCES categories(id);

UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'Clothing') WHERE code IN ('PROD001', 'PROD004', 'PROD007');
UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'Shoes') WHERE code IN ('PROD002', 'PROD006');
UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'Accessories') WHERE code IN ('PROD003', 'PROD005', 'PROD008');