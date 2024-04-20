package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func (c *Client) Keys(ctx context.Context, glob string) ([]string, error) {
	args := []any{
		sql.Named("glob", glob),
		sql.Named("now", time.Now().UnixMilli()),
	}

	rows, err := c.db.QueryContext(ctx, `
	select
		name
	from
		keys
	where
		name GLOB :glob
		and (
			expires_at is null
			or expires_at > :now
		);
	`, args...)
	if err != nil {
		return nil, fmt.Errorf("could not glob names: %w", err)
	}

	defer rows.Close()

	var names []string

	for rows.Next() {
		var name string

		err := rows.Scan(&name)
		if err != nil {
			return nil, fmt.Errorf("could not scan name: %w", err)
		}

		names = append(names, name)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("could not scan rows: %w", rows.Err())
	}

	return names, nil
}
