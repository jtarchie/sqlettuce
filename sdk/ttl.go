package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/georgysavva/scany/v2/sqlscan"
)

func (c *Client) TTL(ctx context.Context, name string) (*time.Duration, error) {
	var value sql.NullInt64

	err := sqlscan.Get(ctx, c.db, &value, `
	select
		expires_at
	from
		active_keys
	where
		name = :name;
	`,
		sql.Named("name", name),
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrKeyDoesNotExist
	}

	if err != nil {
		return nil, fmt.Errorf("could not find key: %w", err)
	}

	if !value.Valid {
		return nil, nil
	}

	delta := time.Until(time.Unix(0, value.Int64*1_000_000))

	return &delta, nil
}
