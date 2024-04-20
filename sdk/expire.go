package sdk

import (
	"database/sql"
	"fmt"
	"time"
)

func (c *Client) Expire(name string, ttl time.Duration) error {
	result, err := c.db.ExecContext(
		c.context,
		`UPDATE keys SET expires_at = :expires_at WHERE name = :name`,
		sql.Named("name", name),
		sql.Named("expires_at", time.Now().Add(ttl).UnixNano()),
	)
	if err != nil {
		return fmt.Errorf("could not expire key: %w", err)
	}

	count, _ := result.RowsAffected()
	if count == 0 {
		return ErrKeyDoesNotExist
	}

	return nil
}
