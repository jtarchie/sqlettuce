package sdk

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jtarchie/sqlettus/executers"
)

func (c *Client) Rename(current string, next string) error {
	err := c.db.WithTX(c.context, func(tx executers.Executer) error {
		_, _ = tx.ExecContext(c.context, `DELETE FROM keys WHERE name = :new`, sql.Named("new", next))

		result, err := tx.ExecContext(c.context, `UPDATE keys SET name = :new WHERE name = :old`, sql.Named("new", next), sql.Named("old", current))
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

func (c *Client) RenameIfNotExists(current string, next string) error {
	err := c.db.WithTX(c.context, func(tx executers.Executer) error {
		row := tx.QueryRowContext(c.context, `SELECT 1 FROM keys WHERE name = :new`, sql.Named("new", next))
		if row.Err() != nil {
			return fmt.Errorf("could not find new: %w", row.Err())
		}

		var value int

		err := row.Scan(&value)
		if errors.Is(err, sql.ErrNoRows) {
			_, err = tx.ExecContext(c.context, `UPDATE keys SET name = :new WHERE name = :old`, sql.Named("new", next), sql.Named("old", current))
			if err != nil {
				return fmt.Errorf("could not rename key: %w", err)
			}
		}

		if err != nil {
			return fmt.Errorf("could not scan: %w", err)
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
