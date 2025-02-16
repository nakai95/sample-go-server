CREATE TABLE chat_messages (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);