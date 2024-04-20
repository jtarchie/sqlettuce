package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (c *Client) DBSize(ctx context.Context) (int64, error) {
	row := c.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM active_keys`,
	)
	if row.Err() != nil {
		return 0, fmt.Errorf("could not read key: %w", row.Err())
	}

	var value int64

	err := row.Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}

	return value, nil
}
