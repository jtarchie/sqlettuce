package sdk

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (c *Client) Get(name string) (string, error) {
	args := []any{
		sql.Named("name", name),
		sql.Named("now", time.Now().UnixMilli()),
	}

	row := c.db.QueryRowContext(c.context, `
	select
		value
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
