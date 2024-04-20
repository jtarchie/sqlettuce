package sdk

import (
	"context"
	"database/sql"
	"fmt"
)

func (c *Client) Delete(ctx context.Context, name string) (bool, error) {
	_, err := c.db.ExecContext(
		ctx,
		`DELETE FROM keys where name = :name;`,
		sql.Named("name", name),
	)
	if err != nil {
		return false, fmt.Errorf("could not delete key: %w", err)
	}

	return true, nil
}
