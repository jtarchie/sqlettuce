package sdk

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/jtarchie/sqlettus/executers"
)

func (c *Client) ListPop(ctx context.Context, name string, length int64) ([]string, error) {
	values := make([]string, 0, length)

	err := c.db.WithTX(ctx, func(tx executers.Executer) error {
		for range length {
			var value string

			err := sqlscan.Get(ctx, tx, &value, `
			SELECT
				jsonb_extract(payload, '$[#-1]')
			FROM
				active_keys keys
			WHERE
				name = :name AND type = :type
			`,
				sql.Named("name", name),
				sql.Named("type", ListType),
			)
			if err != nil {
				return fmt.Errorf("could not pop value: %w", err)
			}

			values = append(values, value)

			_, err = tx.ExecContext(ctx, `
			UPDATE
				keys
			SET
				payload = jsonb_remove(payload, '$[#-1]')
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
