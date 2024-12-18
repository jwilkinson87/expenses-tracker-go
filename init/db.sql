CREATE TABLE users (
    id UUID PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
    id UUID PRIMARY KEY,
    label VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE expenses (
    id UUID PRIMARY KEY,
    amount BIGINT NOT NULL,
    category_id UUID REFERENCES categories (id) ON DELETE SET NULL,
    user_id UUID REFERENCES users (id) ON DELETE CASCADE,
    expense_date DATE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE reset_tokens (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expiry_time TIMESTAMP,
    reset_token VARCHAR(64) NOT NULL,
    user_id UUID REFERNECES users (id) ON DELETE CASCADE
)
