package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/jtarchie/sqlettuce/executers"
)

func (c *Client) ListRemove(ctx context.Context, name string, count int64, element string) (int64, error) {
	orderBy := "ASC"

	if 0 > count {
		orderBy = "DESC"
		count = -count
	}

	if count == 0 {
		count = math.MaxInt64
	}

	var positions []int64

	err := c.db.WithTX(ctx, func(tx executers.Executer) error {
		err := sqlscan.Select(ctx, tx, &positions, fmt.Sprintf(`
		SELECT
			json_each.key AS pos
		FROM
			active_keys keys,
			json_each(keys.payload)
		WHERE
			keys.name = :name
			AND keys.type = :type
			AND json_each.value = :element
		ORDER BY
			pos %s
		LIMIT
			:count;
	`, orderBy),
			sql.Named("name", name),
			sql.Named("type", ListType),
			sql.Named("element", element),
			sql.Named("count", count),
		)
		if err != nil {
			return fmt.Errorf("could not determine positions: %w", err)
		}

		for _, position := range positions {
			_, err = tx.ExecContext(ctx, `
			UPDATE
				keys
			SET
				payload = json_remove(payload, '$[' || :position || ']')
			WHERE
				keys.name = :name AND
				keys.type = :type
		`,
				sql.Named("name", name),
				sql.Named("type", ListType),
				sql.Named("position", position),
			)
			if err != nil {
				return fmt.Errorf("could not remove element: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("could not remove elements: %w", err)
	}

	return int64(len(positions)), nil
}
