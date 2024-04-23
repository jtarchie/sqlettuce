package sdk

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/jtarchie/sqlettuce/executers"
)

func (c *Client) ListPush(ctx context.Context, name string, values ...string) (int64, error) {
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
				(:name, jsonb_array(), :type);
			`,

				sql.Named("name", name),
				sql.Named("type", ListType),
			)
			if err != nil {
				return fmt.Errorf("could not set key: %w", err)
			}
		}

		for _, value := range values {
			_, err := tx.ExecContext(ctx, `
			UPDATE
				keys
			SET
				payload = jsonb_insert(payload, '$[#]', :value)
			WHERE
				name = :name AND
				type = :type
		`,
				sql.Named("name", name),
				sql.Named("type", ListType),
				sql.Named("value", value),
			)
			if err != nil {
				return fmt.Errorf("could not set value: %w", err)
			}
		}

		err = sqlscan.Get(ctx, tx, &length, `
			SELECT json_array_length(payload) FROM active_keys WHERE name = :name
		`,
			sql.Named("name", name),
		)
		if err != nil {
			return fmt.Errorf("could not determine list length: %w", err)
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("transaction failed: %w", err)
	}

	return length, nil
}
