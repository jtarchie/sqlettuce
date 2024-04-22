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
		sql.Named("first", names[0]),
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
			c.name = :first
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

func (c *Client) SetDiffAndStore(ctx context.Context, name string, names ...string) (int64, error) {
	args := []any{
		sql.Named("type", SetType),
		sql.Named("first", names[0]),
		sql.Named("name", name),
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

	var count int64

	err := sqlscan.Get(ctx, c.db, &count, fmt.Sprintf(`
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
		), payloads AS (
			SELECT
				jsonb_group_object(value, 0) AS payload
			FROM
				counts c
			WHERE
				c.count = 1 AND
				c.name = :first
		)
		INSERT INTO
			keys (name, payload, type)
		VALUES
			(:name, (SELECT payload FROM payloads), :type)
		ON CONFLICT(name) DO
		UPDATE
		SET
			payload = excluded.payload
		WHERE
			keys.name = :name AND
			keys.type = :type
		RETURNING
			(
				SELECT
					COUNT(json_each.key)
				FROM
					active_keys keys,
					json_each(keys.payload)
				WHERE
					keys.name = :name AND
					keys.type = :type
			)
	`,
		placeholders.String(),
	),
		args...,
	)
	if err != nil {
		return 0, fmt.Errorf("could not diff values: %w", err)
	}

	return count, nil
}
