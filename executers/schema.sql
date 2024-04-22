PRAGMA journal_mode = WAL;

PRAGMA synchronous = NORMAL;

PRAGMA journal_size_limit = 67108864;

PRAGMA mmap_size = 268435456;

PRAGMA cache_size = 2000;

PRAGMA busy_timeout = 5000;

CREATE TABLE IF NOT EXISTS keys (
  name TEXT NOT NULL PRIMARY KEY,
  value TEXT,
  payload JSONB,
  type INTEGER NOT NULL,
  version INTEGER NOT NULL DEFAULT 0,
  expires_at INTEGER,
  updated_at INTEGER NOT NULL default (unixepoch('now', 'subsec') * 1000)
);

CREATE TRIGGER IF NOT EXISTS trigger_keys_before_update BEFORE
UPDATE
  ON keys BEGIN
SELECT
  NEW.updated_at = (unixepoch('now', 'subsec') * 1000),
  NEW.version = OLD.version + 1;

END;

CREATE INDEX IF NOT EXISTS expires_at_idx ON keys(expires_at);

CREATE VIEW IF NOT EXISTS active_keys AS
SELECT
  *
FROM
  keys
WHERE
  expires_at IS NULL
  OR expires_at > (unixepoch('now', 'subsec') * 1000);