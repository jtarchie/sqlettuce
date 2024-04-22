package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
)

func (c *Client) Exists(ctx context.Context, name string) (bool, error) {
	var value int64

	err := sqlscan.Get(ctx, c.db, &value,
		`SELECT 1 FROM active_keys WHERE name = :name`,
		sql.Named("name", name),
	)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("could not read key: %w", err)
	}

	return true, nil
}
