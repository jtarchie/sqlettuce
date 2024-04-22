package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
)

func (c *Client) SetContains(ctx context.Context, name string, element string) (bool, error) {
	var length sql.NullInt64

	err := sqlscan.Get(ctx, c.db, &length, `
		SELECT
			1
		FROM
			active_keys keys,
			json_each(keys.payload)
		WHERE
			keys.name = :name AND
			keys.type = :type AND
			json_each.key = :element
		LIMIT 1
	`,
		sql.Named("name", name),
		sql.Named("type", SetType),
		sql.Named("element", element),
	)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("could not read list length: %w", err)
	}

	return true, nil
}
