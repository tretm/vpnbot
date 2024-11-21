CREATE TABLE IF NOT EXISTS payment_history (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    amount INT NOT NULL,
    transaction_type VARCHAR(255) NOT NULL,
    comment TEXT,
    time_create TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);