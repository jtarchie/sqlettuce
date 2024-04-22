package sdk

import (
	"context"
	"database/sql"
	"fmt"
)

func (c *Client) ListRange(ctx context.Context, name string, start int64, end int64) ([]string, error) {
	rows, err := c.db.QueryContext(ctx, `
	SELECT
		json_each.value
	FROM
		active_keys keys,
		json_each(keys.value)
	WHERE
		keys.name = :name AND
		keys.type = :type AND
		json_each.key >= IIF(:start >=0, :start, json_array_length(keys.value) + :start) AND
		json_each.key <= IIF(:end >=0, :end, json_array_length(keys.value) + :end);
	`,
		sql.Named("name", name),
		sql.Named("type", ListType),
		sql.Named("start", start),
		sql.Named("end", end),
	)
	if err != nil {
		return nil, fmt.Errorf("could not get values: %w", err)
	}

	defer rows.Close()

	var values []string

	for rows.Next() {
		var name string

		err := rows.Scan(&name)
		if err != nil {
			return nil, fmt.Errorf("could not scan value: %w", err)
		}

		values = append(values, name)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("could not scan rows: %w", rows.Err())
	}

	return values, nil
}
