package sdk

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
)

func (c *Client) SetLength(ctx context.Context, name string) (int64, error) {
	var length int64

	err := sqlscan.Get(ctx, c.db, &length, `
		SELECT
			COUNT(json_each.value)
		FROM
			active_keys keys,
			json_each(keys.payload)
		WHERE
			keys.name = :name AND
			keys.type = :type
	`,
		sql.Named("name", name),
		sql.Named("type", SetType),
	)
	if err != nil {
		return 0, fmt.Errorf("could not read list length: %w", err)
	}

	return length, nil
}
