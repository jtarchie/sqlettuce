package sdk

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
)

func (c *Client) ListRange(ctx context.Context, name string, start int64, end int64) ([]string, error) {
	var values []string

	err := sqlscan.Select(ctx, c.db, &values, `
	SELECT
		json_each.value
	FROM
		active_keys keys,
		json_each(keys.value)
	WHERE
		keys.name = :name AND
		keys.type = :type AND
		json_each.key >= IIF(:start >=0, :start, json_array_length(keys.value) + :start) AND
		json_each.key <= IIF(:end >=0, :end, json_array_length(keys.value) + :end);
	`,
		sql.Named("name", name),
		sql.Named("type", ListType),
		sql.Named("start", start),
		sql.Named("end", end),
	)
	if err != nil {
		return nil, fmt.Errorf("could not get values: %w", err)
	}

	return values, nil
}
