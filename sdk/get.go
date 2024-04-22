package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
)

func (c *Client) Get(ctx context.Context, name string) (string, error) {
	var value string

	err := sqlscan.Get(ctx, c.db, &value, `
	SELECT
		value
	FROM
		active_keys
	WHERE
		name = :name AND
		type = :type;
	`,
		sql.Named("name", name),
		sql.Named("type", StringType),
	)

	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}

	if err != nil {
		return "", fmt.Errorf("could not read value: %w", err)
	}

	return value, nil
}
