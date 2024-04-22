package sdk

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
)

func (c *Client) AddInt(ctx context.Context, name string, add int64) (int64, error) {
	var value int64

	err := sqlscan.Get(ctx, c.db, &value, `
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
	if err != nil {
		return 0, fmt.Errorf("could not set integer: %w", err)
	}

	return value, nil
}
