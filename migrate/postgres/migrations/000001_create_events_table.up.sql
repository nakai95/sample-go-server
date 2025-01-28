CREATE TABLE events (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(300) NOT NULL,
    image_url VARCHAR(2048) NOT NULL
);