package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/v2/sqlscan"
)

func (c *Client) SetDiff(ctx context.Context, names ...string) ([]string, error) {
	args := []any{
		sql.Named("type", SetType),
		sql.Named("name", names[0]),
	}

	placeholders := &strings.Builder{}

	for index, name := range names {
		placeholderName := fmt.Sprintf("p%d", index)
		args = append(args, sql.Named(placeholderName, name))

		placeholders.WriteByte(':')
		placeholders.WriteString(placeholderName)

		if index < len(names)-1 {
			placeholders.WriteByte(',')
		}
	}

	var values []string

	err := sqlscan.Select(ctx, c.db, &values, fmt.Sprintf(`
		WITH counts AS (
			SELECT
				keys.name AS name,
				json_each.key AS value,
				COUNT(json_each.key) AS count
			FROM
				active_keys keys,
				json_each(keys.payload)
			WHERE
				keys.type = :type AND
				keys.name IN (%s)
			GROUP BY
				json_each.key
		)
		SELECT
			value
		FROM
			counts c
		WHERE
			c.count = 1 AND
			c.name = :name
	`,
		placeholders.String(),
	),
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("could not diff values: %w", err)
	}

	return values, nil
}
