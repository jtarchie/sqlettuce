package sdk

import (
	"database/sql"
	"errors"
	"fmt"
)

func (c *Client) RandomKey() (string, error) {
	row := c.db.QueryRowContext(
		c.context,
		`SELECT name FROM keys ORDER BY RANDOM() LIMIT 1`,
	)
	if row.Err() != nil {
		return "", fmt.Errorf("could select random value: %w", row.Err())
	}

	var name string

	err := row.Scan(&name)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrKeyDoesNotExist
	}

	if err != nil {
		return "", fmt.Errorf("could not scan value: %w", err)
	}

	return name, nil
}
