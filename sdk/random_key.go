package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
)

func (c *Client) RandomKey(ctx context.Context) (string, error) {
	var name string

	err := sqlscan.Get(ctx, c.db, &name,
		`SELECT name FROM active_keys ORDER BY RANDOM() LIMIT 1`,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrKeyDoesNotExist
	}

	if err != nil {
		return "", fmt.Errorf("could select random value: %w", err)
	}

	return name, nil
}
