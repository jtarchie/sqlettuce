package sdk

import (
	"context"
	"database/sql"
	"fmt"
)

func (c *Client) Sort(ctx context.Context, name string) ([]string, error) {
	var values []string

	rows, err := c.db.QueryContext(ctx, `
		SELECT
			json_each.value
		FROM
			keys,
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
		return nil, fmt.Errorf("could not query for sorted values: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var value string

		err := rows.Scan(&value)
		if err != nil {
			return nil, fmt.Errorf("could not scan for sorted values: %w", err)
		}

		values = append(values, value)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("could not scan rows: %w", rows.Err())
	}

	return values, nil
}
