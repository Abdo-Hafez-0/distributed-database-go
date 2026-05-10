-- =============================================================
-- sql/seed/seed.sql
-- Sample data for development and testing.
-- Run AFTER schema.sql.
-- =============================================================

-- Users (passwords are bcrypt hashes of "password123")
INSERT INTO `users` (`username`, `email`, `password`, `role`, `active`) VALUES
    ('admin',   'admin@example.com',   '$2a$12$examplehashADMIN000000000000000000000000000000000000000', 'admin', TRUE),
    ('alice',   'alice@example.com',   '$2a$12$examplehashALICE000000000000000000000000000000000000000', 'user',  TRUE),
    ('bob',     'bob@example.com',     '$2a$12$examplehashBOB0000000000000000000000000000000000000000', 'user',  TRUE),
    ('charlie', 'charlie@example.com', '$2a$12$examplehashCHARLIE00000000000000000000000000000000000', 'user',  FALSE);

-- Products
INSERT INTO `products` (`name`, `description`, `price`, `stock`, `category`) VALUES
    ('Laptop Pro 15',   'High-performance laptop with 16 GB RAM',   1299.99, 50,  'Electronics'),
    ('Wireless Mouse',  'Ergonomic wireless mouse, 2.4 GHz',           29.99, 200, 'Electronics'),
    ('Desk Chair',      'Adjustable lumbar-support office chair',      249.99, 30,  'Furniture'),
    ('Standing Desk',   'Height-adjustable sit/stand desk, 120 cm',   499.99, 15,  'Furniture'),
    ('USB-C Hub 7-in-1','Multi-port USB-C hub with HDMI and PD',        49.99, 100, 'Electronics');

-- Orders
INSERT INTO `orders` (`user_id`, `total`, `status`) VALUES
    (2, 1329.98, 'completed'),   -- alice
    (3,   49.99, 'pending'),     -- bob
    (2,  249.99, 'processing');  -- alice second order

-- Order items
INSERT INTO `order_items` (`order_id`, `product_id`, `quantity`, `unit_price`) VALUES
    (1, 1, 1, 1299.99),   -- alice: Laptop Pro 15
    (1, 2, 1,   29.99),   -- alice: Wireless Mouse
    (2, 5, 1,   49.99),   -- bob:   USB-C Hub
    (3, 3, 1,  249.99);   -- alice: Desk Chair