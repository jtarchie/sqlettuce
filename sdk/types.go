package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
)

type Type int

const (
	StringType Type = Type(1)
	ListType   Type = Type(2)
	SetType    Type = Type(3)
)

func (t Type) String() string {
	switch t {
	case StringType:
		return "string"
	case ListType:
		return "list"
	case SetType:
		return "set"
	default:
		return "unknown"
	}
}

func (c *Client) Type(ctx context.Context, name string) (Type, error) {
	var value Type

	err := sqlscan.Get(ctx, c.db, &value,
		`SELECT type FROM active_keys WHERE name = :name`,
		sql.Named("name", name),
	)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrKeyDoesNotExist
	}

	if err != nil {
		return 0, fmt.Errorf("could not read key: %w", err)
	}

	return value, nil
}
