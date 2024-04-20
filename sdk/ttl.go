package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (c *Client) TTL(ctx context.Context, name string) (*time.Duration, error) {
	row := c.db.QueryRowContext(ctx, `
	select
		expires_at
	from
		keys
	where
		name = :name
		and (
			expires_at is null
			or expires_at > :now
		);
	`,
		sql.Named("name", name),
		sql.Named("now", time.Now().UnixMilli()),
	)

	err := row.Err()
	if err != nil {
		return nil, fmt.Errorf("could not find key: %w", err)
	}

	var value sql.NullInt64

	err = row.Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrKeyDoesNotExist
	}

	if err != nil {
		return nil, fmt.Errorf("could not read value: %w", err)
	}

	if !value.Valid {
		return nil, nil
	}

	delta := time.Until(time.Unix(0, value.Int64*1_000_000))

	return &delta, nil
}
