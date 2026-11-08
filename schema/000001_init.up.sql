CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username VARCHAR(15) NOT NULL UNIQUE,
  password_hash VARCHAR(64) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMENT ON TABLE users IS 'Таблица пользователей.';
COMMENT ON COLUMN users.id IS 'Уникальный идентификатор пользователя.';
COMMENT ON COLUMN users.username IS 'Имя пользователя.';
COMMENT ON COLUMN users.password_hash IS 'Хеш пароля пользователя.';
COMMENT ON COLUMN users.created_at IS 'Дата и время создания пользователя.';
COMMENT ON COLUMN users.updated_at IS 'Дата и время обновления информации о пользователе.';
