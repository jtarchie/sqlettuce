package sdk

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jtarchie/sqlettus/executers"
)

func (c *Client) ListShift(ctx context.Context, name string, length int64) ([]string, error) {
	values := make([]string, 0, length)

	err := c.db.WithTX(ctx, func(tx executers.Executer) error {
		for range length {
			row := tx.QueryRowContext(ctx, `
			SELECT
				json_extract(value, '$[0]')
			FROM
				active_keys keys
			WHERE
				name = :name AND type = :type
			`,
				sql.Named("name", name),
				sql.Named("type", ListType),
			)

			err := row.Err()
			if err != nil {
				return fmt.Errorf("could not pop value: %w", err)
			}

			var value string

			err = row.Scan(&value)
			if err != nil {
				return fmt.Errorf("could not scan value: %w", err)
			}

			values = append(values, value)

			_, err = tx.ExecContext(ctx, `
			UPDATE
				keys
			SET
				value = json_remove(value, '$[0]')
			WHERE
				name = :name AND type = :type
			`,
				sql.Named("name", name),
				sql.Named("type", ListType),
			)
			if err != nil {
				return fmt.Errorf("could not pop value: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	return values, nil
}
