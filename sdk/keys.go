package sdk

import (
	"context"
	"database/sql"
	"fmt"
)

func (c *Client) Keys(ctx context.Context, glob string) ([]string, error) {
	rows, err := c.db.QueryContext(ctx, `
	select
		name
	from
		active_keys
	where
		name GLOB :glob;
	`,
		sql.Named("glob", glob),
	)
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
