CREATE TABLE IF NOT EXISTS users (
    id SERIAL NOT NULL UNIQUE,
    role VARCHAR(55) NOT NULL,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL,
    registered_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL NOT NULL UNIQUE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token UUID NOT NULL DEFAULT gen_random_uuid(),
    expires_at TIMESTAMP NOT NULL
);
