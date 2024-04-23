package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/jtarchie/sqlettuce/executers"
)

func (c *Client) Get(ctx context.Context, name string) (*string, error) {
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
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("could not read value: %w", err)
	}

	return &value, nil
}

func (c *Client) MGet(ctx context.Context, names ...string) ([]*string, error) {
	values := make([]*string, 0, len(names))

	err := c.db.WithTX(ctx, func(tx executers.Executer) error {
		client := NewClient(tx)

		for _, name := range names {
			value, err := client.Get(ctx, name)
			if err != nil {
				return fmt.Errorf("could not read value: %w", err)
			}

			values = append(values, value)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not read values: %w", err)
	}

	return values, nil
}
