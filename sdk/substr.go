package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (c *Client) Substr(ctx context.Context, name string, start int64, end int64) (string, error) {
	row := c.db.QueryRowContext(ctx, `
	SELECT SUBSTR(
    value,
    IIF(:start < 0, :start, :start + 1),
    IIF(
      :end < 0,
      LENGTH(value) - :end,
      :start + :end + 1
    )
  )
	FROM
		active_keys keys
	WHERE
		name = :name AND
		type = :type;
	`,
		sql.Named("name", name),
		sql.Named("type", StringType),
		sql.Named("start", start),
		sql.Named("end", end),
	)

	err := row.Err()
	if err != nil {
		return "", fmt.Errorf("could not find key: %w", err)
	}

	var value string

	err = row.Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}

	if err != nil {
		return "", fmt.Errorf("could not read value: %w", err)
	}

	return value, nil
}
