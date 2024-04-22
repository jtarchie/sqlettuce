package sdk

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
)

func (c *Client) Keys(ctx context.Context, glob string) ([]string, error) {
	var names []string

	err := sqlscan.Select(ctx, c.db, &names, `
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

	return names, nil
}
