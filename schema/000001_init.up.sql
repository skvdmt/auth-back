CREATE TABLE IF NOT EXISTS refresh_tokens (
  id UUID NOT NULL DEFAULT uuidv7(),
  username VARCHAR(15) NOT NULL,
  refresh_token VARCHAR(4096) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'Europe/Moscow'),
  PRIMARY KEY(id)
);

COMMENT ON TABLE refresh_tokens IS 'Таблица токенов обновления.';
COMMENT ON COLUMN refresh_tokens.id IS 'Уникальный идентификатор токена обновления.';
COMMENT ON COLUMN refresh_tokens.username IS 'Имя пользователя, токена обновления.';
COMMENT ON COLUMN refresh_tokens.refresh_token IS 'Токен обновления.';
COMMENT ON COLUMN refresh_tokens.created_at IS 'Дата и время создания токена обновления.';
