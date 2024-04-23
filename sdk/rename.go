package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/jtarchie/sqlettuce/executers"
)

func (c *Client) Rename(ctx context.Context, current string, next string) error {
	err := c.db.WithTX(ctx, func(tx executers.Executer) error {
		_, _ = tx.ExecContext(ctx, `DELETE FROM keys WHERE name = :new`, sql.Named("new", next))

		result, err := tx.ExecContext(ctx, `UPDATE keys SET name = :new WHERE name = :old`, sql.Named("new", next), sql.Named("old", current))
		if err != nil {
			return fmt.Errorf("could not rename key: %w", err)
		}

		count, _ := result.RowsAffected()
		if count == 0 {
			return ErrKeyDoesNotExist
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("could not rename: %w", err)
	}

	return nil
}

func (c *Client) RenameIfNotExists(ctx context.Context, current string, next string) error {
	err := c.db.WithTX(ctx, func(tx executers.Executer) error {
		var value int

		err := sqlscan.Get(ctx, tx, &value,
			`SELECT 1 FROM keys WHERE name = :new`,
			sql.Named("new", next),
		)

		if errors.Is(err, sql.ErrNoRows) {
			_, err = tx.ExecContext(ctx, `UPDATE keys SET name = :new WHERE name = :old`, sql.Named("new", next), sql.Named("old", current))
			if err != nil {
				return fmt.Errorf("could not rename key: %w", err)
			}
		}

		if err != nil {
			return fmt.Errorf("could not find new: %w", err)
		}

		if value == 1 {
			return ErrKeyAlreadyExists
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("could not rename: %w", err)
	}

	return nil
}
