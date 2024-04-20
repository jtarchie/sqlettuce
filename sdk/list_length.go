package sdk

import (
	"context"
	"database/sql"
	"fmt"
)

func (c *Client) ListLength(ctx context.Context, name string) (int64, error) {
	row := c.db.QueryRowContext(ctx, `
		SELECT
			json_array_length(value)
		FROM
			active_keys keys
		WHERE
			keys.name = :name AND
			keys.type = :type
	`,
		sql.Named("name", name),
		sql.Named("type", ListType),
	)

	err := row.Err()
	if err != nil {
		return 0, fmt.Errorf("could not read list length: %w", err)
	}

	var length int64

	err = row.Scan(&length)
	if err != nil {
		return 0, fmt.Errorf("could not scan: %w", err)
	}

	return length, nil
}
