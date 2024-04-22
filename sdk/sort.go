package sdk

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
)

func (c *Client) Sort(ctx context.Context, name string) ([]string, error) {
	var values []string

	err := sqlscan.Select(ctx, c.db, &values, `
		SELECT
			json_each.value
		FROM
			active_keys keys,
			json_each(keys.value)
		WHERE
			keys.name = :name AND
			keys.type = :type
		ORDER BY
			CAST(json_each.value AS REAL)
	`,
		sql.Named("name", name),
		sql.Named("type", ListType),
	)
	if err != nil {
		return nil, fmt.Errorf("could sort values: %w", err)
	}

	return values, nil
}
