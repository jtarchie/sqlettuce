package sdk

import (
	"database/sql"
	"fmt"
	"time"
)

func (c *Client) Expire(name string, ttl time.Duration) error {
	result, err := c.db.ExecContext(
		c.context,
		`UPDATE keys SET etime = :etime WHERE name = :name`,
		sql.Named("name", name),
		sql.Named("etime", time.Now().Add(ttl).UnixNano()),
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
