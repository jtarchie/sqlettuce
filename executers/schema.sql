PRAGMA journal_mode = WAL;

PRAGMA synchronous = NORMAL;

PRAGMA journal_size_limit = 67108864;

PRAGMA mmap_size = 268435456;

PRAGMA cache_size = 2000;

PRAGMA busy_timeout = 5000;

CREATE TABLE IF NOT EXISTS keys (
  name TEXT NOT NULL PRIMARY KEY,
  value TEXT NOT NULL,
  type INTEGER NOT NULL,
  version INTEGER not null default 0,
  expires_at INTEGER,
  updated_at INTEGER not null
);

CREATE INDEX IF NOT EXISTS expires_at_idx ON keys(expires_at);