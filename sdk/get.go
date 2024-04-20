package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (c *Client) Get(ctx context.Context, name string) (string, error) {
	args := []any{
		sql.Named("name", name),
		sql.Named("now", time.Now().UnixNano()),
	}

	row := c.db.QueryRowContext(ctx, `
	select
		value
	from
		keys
	where
		name = :name
		and (
			expires_at is null
			or expires_at > :now
		);
	`, args...)

	err := row.Err()
	if err != nil {
		return "", fmt.Errorf("could not find key: %w", err)
	}

	var value string

	err = row.Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}

	if err != nil {
		return "", fmt.Errorf("could not read value: %w", err)
	}

	return value, nil
}
