package sdk

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/jtarchie/sqlettus/executers"
	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrKeyAlreadyExists = errors.New("key already exists")
	ErrKeyDoesNotExist  = errors.New("key does not exist")
)

type Client struct {
	context context.Context
	db      executers.Executer
}

func NewClient(ctx context.Context, db executers.Executer) (*Client, error) {
	return &Client{
		db:      db,
		context: ctx,
	}, nil
}

func (c *Client) Close() error {
	err := c.db.Close()
	if err != nil {
		return fmt.Errorf("could not close db: %w", err)
	}

	return nil
}
