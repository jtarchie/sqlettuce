package sdk

import (
	"context"
	"fmt"
)

func (c *Client) FlushDB(ctx context.Context) error {
	_, err := c.db.ExecContext(ctx, `
		DELETE FROM keys;
		VACUUM;
		PRAGMA OPTIMIZE;
	`)
	if err != nil {
		return fmt.Errorf("could not flush db: %w", err)
	}

	return nil
}
