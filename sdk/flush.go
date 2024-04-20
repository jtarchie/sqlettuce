package sdk

import "fmt"

func (c *Client) FlushDB() error {
	_, err := c.db.ExecContext(c.context, `
		DELETE FROM keys;
		VACUUM;
		PRAGMA OPTIMIZE;
	`)
	if err != nil {
		return fmt.Errorf("could not flush db: %w", err)
	}

	return nil
}
