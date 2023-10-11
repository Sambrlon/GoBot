CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    chat_id BIGINT,
    username TEXT,
    text TEXT,
    is_admin BOOLEAN,
    timestamp TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    username VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS clients (
    id SERIAL PRIMARY KEY,
    chat_id BIGINT
);

-- Удалите уникальное ограничение по имени столбца chat_id (замените clients на имя вашей таблицы)
ALTER TABLE messages DROP CONSTRAINT IF EXISTS clients_chat_id_key;
