package sdk

import (
	"context"
	"fmt"

	"github.com/jtarchie/sqlettuce/executers"
)

func (c *Client) GetSet(ctx context.Context, name string, value string) (string, error) {
	var prevValue *string

	err := c.db.WithTX(ctx, func(tx executers.Executer) error {
		var err error

		client := NewClient(tx)

		prevValue, err = client.Get(ctx, name)
		if err != nil {
			return fmt.Errorf("could not read value: %w", err)
		}

		err = client.Set(ctx, name, value, 0)
		if err != nil {
			return fmt.Errorf("could not set value: %w", err)
		}

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("could not read values: %w", err)
	}

	return *prevValue, nil
}
