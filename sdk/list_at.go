package sdk

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
)

func (c *Client) ListAt(ctx context.Context, name string, index int64) (string, error) {
	var value string

	err := sqlscan.Get(ctx, c.db, &value, `
		SELECT
			json_extract(value, '$[#' || :index || ']')
		FROM
			active_keys keys
		WHERE
			keys.name = :name AND
			keys.type = :type;
	`,
		sql.Named("name", name),
		sql.Named("index", index),
		sql.Named("type", ListType),
	)
	if err != nil {
		return "", fmt.Errorf("could not extract value: %w", err)
	}

	return value, nil
}
