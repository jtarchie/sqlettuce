package sdk

import (
	"context"
	"database/sql"
	"fmt"
)

func (c *Client) ListAt(ctx context.Context, name string, index int64) (string, error) {
	row := c.db.QueryRowContext(ctx, `
		SELECT
			json_extract(value, '$[#' || :index || ']')
		FROM
			active_keys keys
		WHERE
			keys.name = :name AND
			keys.type = :type;
	`,
		sql.Named("name", name),
		sql.Named("index", index),
		sql.Named("type", ListType),
	)

	err := row.Err()
	if err != nil {
		return "", fmt.Errorf("could not extract value: %w", err)
	}

	var value string

	err = row.Scan(&value)
	if err != nil {
		return "", fmt.Errorf("could not scan value: %w", err)
	}

	return value, nil
}
