package sdk

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
)

func (c *Client) SetMembers(ctx context.Context, name string) ([]string, error) {
	var values []string

	err := sqlscan.Select(ctx, c.db, &values, `
	SELECT
		json_each.key
	FROM
		active_keys keys,
		json_each(keys.payload)
	WHERE
		keys.name = :name AND
		keys.type = :type;
	`,
		sql.Named("name", name),
		sql.Named("type", SetType),
	)
	if err != nil {
		return nil, fmt.Errorf("could not get values: %w", err)
	}

	return values, nil
}
