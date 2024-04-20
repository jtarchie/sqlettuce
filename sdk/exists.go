package sdk

import (
	"database/sql"
	"errors"
	"fmt"
)

func (c *Client) Exists(name string) (bool, error) {
	row := c.db.QueryRowContext(
		c.context,
		`SELECT 1 FROM keys WHERE name = :name`,
		sql.Named("name", name),
	)
	if row.Err() != nil {
		return false, fmt.Errorf("could not read key: %w", row.Err())
	}

	var value int

	err := row.Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	return true, nil
}
