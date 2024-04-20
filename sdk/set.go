package sdk

import (
	"database/sql"
	"fmt"
	"time"
)

func (c *Client) Set(name string, value any, ttl time.Duration) error {
	now := time.Now()

	var etime *int64
	if ttl > 0 {
		etime = new(int64)
		*etime = now.Add(ttl).UnixMilli()
	}

	args := []any{
		sql.Named("name", name),
		sql.Named("value", value),
		sql.Named("etime", etime),
		sql.Named("mtime", now.UnixMilli()),
	}

	_, err := c.db.ExecContext(c.context, `
		INSERT INTO
			keys (name, value, etime, mtime)
		values
			(:name, :value, :etime, :mtime) ON CONFLICT (name) do
		UPDATE
		SET
			version = version + 1,
			value = excluded.value,
			etime = excluded.etime,
			mtime = excluded.mtime
	`, args...)
	if err != nil {
		return fmt.Errorf("could not set key: %w", err)
	}

	return nil
}
