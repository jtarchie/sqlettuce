package sdk

import (
	"context"
	"database/sql"
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

//go:embed schema.sql
var schemaSQL string

type Client struct {
	context context.Context
	db      executers.Executer
}

func NewClient(ctx context.Context, filename string) (*Client, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, fmt.Errorf("could open sqlite3: %w", err)
	}

	_, err = db.ExecContext(ctx, schemaSQL)
	if err != nil {
		return nil, fmt.Errorf("could not create schema: %w", err)
	}

	db.SetMaxOpenConns(1)

	return &Client{
		db:      executers.NewPrepared(db),
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
