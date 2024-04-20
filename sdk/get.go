package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (c *Client) Get(ctx context.Context, name string) (string, error) {
	row := c.db.QueryRowContext(ctx, `
	SELECT
		value
	FROM
		active_keys
	WHERE
		name = :name AND
		type = :type;
	`,
		sql.Named("name", name),
		sql.Named("type", StringType),
	)

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
