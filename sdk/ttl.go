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
		etime
	from
		keys
	where
		name = :name
		and (
			etime is null
			or etime > :now
		);
	`, args...)

	err := row.Err()
	if err != nil {
		return nil, fmt.Errorf("could not find key: %w", err)
	}

	var value *int64

	err = row.Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrKeyDoesNotExist
	}

	if err != nil {
		return nil, fmt.Errorf("could not read value: %w", err)
	}

	delta := int64(time.Until(time.Unix(0, *value)))

	return &delta, nil
}
