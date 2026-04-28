CREATE TABLE IF NOT EXISTS cart_items (
    user_id BIGINT REFERENCES users(id),
    product_id BIGINT REFERENCES products(id),
    quantity INT NOT NULL,
    PRIMARY KEY (user_id, product_id)
);
