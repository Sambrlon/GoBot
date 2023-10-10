-- Создание таблицы для хранения сообщений
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    chat_id BIGINT,
    username TEXT,
    text TEXT,
    is_admin BOOLEAN,
    timestamp TIMESTAMP
);
