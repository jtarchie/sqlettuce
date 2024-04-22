package sdk

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
)

func (c *Client) DBSize(ctx context.Context) (int64, error) {
	var value int64

	err := sqlscan.Get(ctx, c.db, &value,
		`SELECT COUNT(*) FROM active_keys`,
	)
	if err != nil {
		return 0, fmt.Errorf("could not count key: %w", err)
	}

	return value, nil
}
