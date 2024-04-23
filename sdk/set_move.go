package sdk

import (
	"context"
	"errors"
	"fmt"

	"github.com/jtarchie/sqlettuce/executers"
)

var (
	ErrElementNotInSource           = errors.New("no element in source")
	ErrElementNotAddedInDestination = errors.New("could not add element to destination")
)

func (c *Client) SetMove(ctx context.Context, destination string, source string, element string) error {
	err := c.db.WithTX(ctx, func(tx executers.Executer) error {
		client := NewClient(tx)

		moved, err := client.SetRemove(ctx, source, element)
		if err != nil {
			return fmt.Errorf("could remove from source: %w", err)
		}

		if moved == 0 {
			return ErrElementNotInSource
		}

		added, err := client.SetAdd(ctx, destination, element)
		if err != nil {
			return fmt.Errorf("could add to destination: %w", err)
		}

		if added == 0 {
			return ErrElementNotAddedInDestination
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("could not move between sets: %w", err)
	}

	return nil
}
