package sdk

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
)

func (c *Client) SetRandomMember(ctx context.Context, name string, length int64) ([]string, error) {
	values := make([]string, 0, length)

	err := sqlscan.Select(ctx, c.db, &values, `
			SELECT
				json_each.key
			FROM
				active_keys keys,
				json_each(keys.payload)
			WHERE
				keys.name = :name AND
				keys.type = :type
			LIMIT :limit
			`,
		sql.Named("name", name),
		sql.Named("type", SetType),
		sql.Named("limit", length),
	)
	if err != nil {
		return nil, fmt.Errorf("could not pop values: %w", err)
	}

	return values, nil
}
