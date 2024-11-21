CREATE TABLE IF NOT EXISTS users (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    user_name VARCHAR(255) NOT NULL,
    message_type INT NOT NULL,
    command VARCHAR(255) NOT NULL,
    status INT NOT NULL,
    history_user_name VARCHAR(255) NOT NULL,
    balance INT NOT NULL,
    balance_all_time INT NOT NULL,
    auto_pay TINYINT(1) NOT NULL,
    test_used TINYINT(1) NOT NULL,
    last_msg_id BIGINT NOT NULL,
    referal_id BIGINT NOT NULL,
    time_create TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    time_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
