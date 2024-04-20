package sdk

import (
	"database/sql"
	"fmt"
	"time"
)

func (c *Client) Set(name string, value any, ttl time.Duration) error {
	now := time.Now()

	var expiresAt *int64
	if ttl > 0 {
		expiresAt = new(int64)
		*expiresAt = now.Add(ttl).UnixNano()
	}

	args := []any{
		sql.Named("name", name),
		sql.Named("value", value),
		sql.Named("expires_at", expiresAt),
		sql.Named("updated_at", now.UnixNano()),
	}

	_, err := c.db.ExecContext(c.context, `
		INSERT INTO
			keys (name, value, expires_at, updated_at)
		values
			(:name, :value, :expires_at, :updated_at) ON CONFLICT (name) do
		UPDATE
		SET
			version = version + 1,
			value = excluded.value,
			expires_at = excluded.expires_at,
			updated_at = excluded.updated_at
	`, args...)
	if err != nil {
		return fmt.Errorf("could not set key: %w", err)
	}

	return nil
}
