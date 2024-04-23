package sdk

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jtarchie/sqlettuce/executers"
)

func (c *Client) SetRemove(ctx context.Context, name string, values ...string) (int64, error) {
	var length int64

	err := c.db.WithTX(ctx, func(tx executers.Executer) error {
		for _, value := range values {
			result, err := tx.ExecContext(ctx, `
			UPDATE
				keys
			SET
				payload = jsonb_remove(payload, '$.' || :value)
			WHERE
				name = :name AND
				type = :type AND
				jsonb_extract(payload, '$.' || :value) IS NOT NULL;
		`,
				sql.Named("name", name),
				sql.Named("type", SetType),
				sql.Named("value", value),
			)
			if err != nil {
				return fmt.Errorf("could not set value: %w", err)
			}

			if effected, err := result.RowsAffected(); err == nil {
				length += effected
			}
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("transaction failed: %w", err)
	}

	return length, nil
}
