CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(50) PRIMARY KEY,
    email VARCHAR(100) NOT NULL,
    phone VARCHAR(50),
    phone_country VARCHAR(25),
    password VARCHAR(255) NOT NULL,
    email_verify BOOLEAN DEFAULT FALSE,
    phone_verify BOOLEAN DEFAULT FALSE,
    avatar VARCHAR(255),
    user_name VARCHAR(30) NOT NULL,
    github_connect VARCHAR(100),
    google_connect VARCHAR(100),
    create_at TIMESTAMP,
    update_at TIMESTAMP
);