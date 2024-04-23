package sdk

import (
	"context"
	"fmt"

	"github.com/jtarchie/sqlettuce/executers"
)

func (c *Client) SetPop(ctx context.Context, name string, count int64) ([]string, error) {
	var values []string

	err := c.db.WithTX(ctx, func(tx executers.Executer) error {
		var err error

		client := NewClient(tx)

		values, err = client.SetRandomMember(ctx, name, count)
		if err != nil {
			return fmt.Errorf("could not get random elements: %w", err)
		}

		_, err = client.SetRemove(ctx, name, values...)
		if err != nil {
			return fmt.Errorf("could not remove elements: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not run transaction: %w", err)
	}

	return values, nil
}
