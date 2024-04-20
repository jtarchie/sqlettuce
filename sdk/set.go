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
		*expiresAt = now.Add(ttl).UnixNano()
	}

	args := []any{
		sql.Named("name", name),
		sql.Named("value", value),
		sql.Named("type", StringType),
		sql.Named("expires_at", expiresAt),
		sql.Named("updated_at", now.UnixNano()),
	}

	_, err := c.db.ExecContext(ctx, `
		INSERT INTO
			keys (name, value, type, expires_at, updated_at)
		values
			(:name, :value, :type, :expires_at, :updated_at) ON CONFLICT (name) do
		UPDATE
		SET
			version = version + 1,
			value = excluded.value,
			expires_at = excluded.expires_at,
			updated_at = excluded.updated_at
	`, args...)
	if err != nil {
		return fmt.Errorf("could not set key: %w", err)
	}

	return nil
}

func (c *Client) MSet(ctx context.Context, pairs ...[2]string) error {
	err := c.db.WithTX(ctx, func(tx executers.Executer) error {
		client, _ := NewClient(tx)

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
