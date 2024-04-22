package sdk

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
)

func (c *Client) ListLength(ctx context.Context, name string) (int64, error) {
	var length int64

	err := sqlscan.Get(ctx, c.db, &length, `
		SELECT
			json_array_length(payload)
		FROM
			active_keys keys
		WHERE
			keys.name = :name AND
			keys.type = :type
	`,
		sql.Named("name", name),
		sql.Named("type", ListType),
	)
	if err != nil {
		return 0, fmt.Errorf("could not read list length: %w", err)
	}

	return length, nil
}
