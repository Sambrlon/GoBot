CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    chat_id BIGINT,
    message_id BIGINT,
    username TEXT,
    text TEXT,
    is_admin BOOLEAN,
    timestamp TIMESTAMP
);

ALTER TABLE messages DROP CONSTRAINT IF EXISTS clients_chat_id_key;
