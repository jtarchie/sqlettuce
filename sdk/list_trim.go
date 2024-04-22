package sdk

import (
	"context"
	"database/sql"
	"fmt"
)

func (c *Client) ListTrim(ctx context.Context, name string, start int64, end int64) error {
	_, err := c.db.ExecContext(ctx, `
	UPDATE
		keys
	SET
		value = (
			SELECT
				json_group_array(json_each.value)
			FROM
				active_keys keys,
				json_each(keys.value)
			WHERE
				keys.name = :name AND
				keys.type = :type AND
				json_each.key >= IIF(:start >=0, :start, json_array_length(keys.value) + :start) AND
				json_each.key <= IIF(:end >=0, :end, json_array_length(keys.value) + :end)
		)
	WHERE
		keys.name = :name AND
		keys.type = :type
	`,
		sql.Named("name", name),
		sql.Named("type", ListType),
		sql.Named("start", start),
		sql.Named("end", end),
	)
	if err != nil {
		return fmt.Errorf("could not trim values: %w", err)
	}

	return nil
}
