package sdk

import (
	"context"
	"database/sql"
	"fmt"
)

func (c *Client) AddInt(ctx context.Context, name string, add int64) (int64, error) {
	row := c.db.QueryRowContext(ctx, `
		INSERT INTO
			keys (name, value, type)
		VALUES
			(:name, :value, :type) ON CONFLICT(name) DO
		UPDATE
		SET
			value = CAST(value AS INTEGER) + CAST(:value AS INTEGER)
		WHERE
			printf("%d", value) = value AND
			type = :type
		RETURNING
			CAST(value AS INTEGER);
	`,
		sql.Named("name", name),
		sql.Named("value", add),
		sql.Named("type", StringType),
	)

	err := row.Err()
	if err != nil {
		return 0, fmt.Errorf("could not set integer: %w", err)
	}

	var value int64

	err = row.Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("could not scan: %w", err)
	}

	return value, nil
}
