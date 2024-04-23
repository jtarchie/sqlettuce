package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Type int

const (
	StringType = Type(1)
)

func (t Type) String() string {
	switch t {
	case StringType:
		return "string"
	default:
		return "unknown"
	}
}

func (c *Client) Type(ctx context.Context, name string) (Type, error) {
	row := c.db.QueryRowContext(
		ctx,
		`SELECT type FROM keys WHERE name = :name`,
		sql.Named("name", name),
	)
	if row.Err() != nil {
		return 0, fmt.Errorf("could not read key: %w", row.Err())
	}

	var value Type

	err := row.Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrKeyDoesNotExist
	}

	return value, nil
}
