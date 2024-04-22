package sdk

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jtarchie/sqlettus/executers"
)

func (c *Client) SetAdd(ctx context.Context, name string, values ...string) (int64, error) {
	var length int64

	err := c.db.WithTX(ctx, func(tx executers.Executer) error {
		client := NewClient(tx)

		found, err := client.Exists(ctx, name)
		if err != nil {
			return fmt.Errorf("could not lookup key: %w", err)
		}

		if !found {
			_, err := tx.ExecContext(ctx, `
			INSERT INTO
				keys (name, payload, type)
			values
				(:name, jsonb_object(), :type);
			`,

				sql.Named("name", name),
				sql.Named("type", SetType),
			)
			if err != nil {
				return fmt.Errorf("could not set key: %w", err)
			}
		}

		for _, value := range values {
			result, err := tx.ExecContext(ctx, `
			UPDATE
				keys
			SET
				payload = jsonb_insert(payload, '$.' || :value, 0)
			WHERE
				name = :name AND
				type = :type AND
				jsonb_extract(payload, '$.' || :value) IS NULL
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
