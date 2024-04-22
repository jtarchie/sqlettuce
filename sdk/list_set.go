package sdk

import (
	"context"
	"database/sql"
	"fmt"
)

func (c *Client) ListSet(ctx context.Context, name string, index int64, element string) error {
	_, err := c.db.ExecContext(ctx, `
		UPDATE
			keys
		SET
		payload = json_replace(payload, '$[' || :index || ']', :element)
		WHERE
			name = :name AND
			type = :type
	`,

		sql.Named("name", name),
		sql.Named("type", ListType),
		sql.Named("index", index),
		sql.Named("element", element),
	)
	if err != nil {
		return fmt.Errorf("could not set element: %w", err)
	}

	return nil
}
