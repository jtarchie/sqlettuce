package sdk

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (c *Client) TTL(name string) (*int64, error) {
	args := []any{
		sql.Named("name", name),
		sql.Named("now", time.Now().UnixNano()),
	}

	row := c.db.QueryRowContext(c.context, `
	select
		expires_at
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
		return nil, fmt.Errorf("could not find key: %w", err)
	}

	var value sql.NullInt64

	err = row.Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrKeyDoesNotExist
	}

	if err != nil {
		return nil, fmt.Errorf("could not read value: %w", err)
	}

	if !value.Valid {
		return nil, nil
	}

	delta := int64(time.Until(time.Unix(0, value.Int64)))

	return &delta, nil
}
