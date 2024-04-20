package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jtarchie/sqlettus/executers"
)

func (c *Client) Set(ctx context.Context, name string, value any, ttl time.Duration) error {
	now := time.Now()

	var expiresAt *int64
	if ttl > 0 {
		expiresAt = new(int64)
		*expiresAt = now.Add(ttl).UnixMilli()
	}

	_, err := c.db.ExecContext(ctx, `
		INSERT INTO
			keys (name, value, type, expires_at)
		values
			(:name, :value, :type, :expires_at) ON CONFLICT (name) do
		UPDATE
		SET
			value = excluded.value,
			type = excluded.type,
			expires_at = excluded.expires_at
	`,
		sql.Named("name", name),
		sql.Named("value", value),
		sql.Named("type", StringType),
		sql.Named("expires_at", expiresAt),
	)
	if err != nil {
		return fmt.Errorf("could not set key: %w", err)
	}

	return nil
}

func (c *Client) MSet(ctx context.Context, pairs ...[2]string) error {
	err := c.db.WithTX(ctx, func(tx executers.Executer) error {
		client := NewClient(tx)

		for _, pair := range pairs {
			err := client.Set(ctx, pair[0], pair[1], 0)
			if err != nil {
				return fmt.Errorf("could not use set: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("could not use transaction: %w", err)
	}

	return nil
}
