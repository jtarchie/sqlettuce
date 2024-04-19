package sqlettus

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "embed"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schemaSQL string

type Client struct {
	context  context.Context
	db       *sql.DB
	prepared map[string]*sql.Stmt
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
		db:       db,
		context:  ctx,
		prepared: map[string]*sql.Stmt{},
	}, nil
}

func (c *Client) Set(name string, value any, ttl time.Duration) error {
	now := time.Now()
	var etime *int64
	if ttl > 0 {
		etime = new(int64)
		*etime = now.Add(ttl).UnixMilli()
	}

	args := []any{
		sql.Named("name", name),
		sql.Named("value", value),
		sql.Named("etime", etime),
		sql.Named("mtime", now.UnixMilli()),
	}

	_, err := c.exec(`
		INSERT INTO
			keys (name, value, etime, mtime)
		values
			(:name, :value, :etime, :mtime) ON CONFLICT (name) do
		UPDATE
		SET
			version = version + 1,
			value = excluded.value,
			etime = excluded.etime,
			mtime = excluded.mtime
	`, args...)
	if err != nil {
		return fmt.Errorf("could not set key: %q", err)
	}

	return nil
}

func (c *Client) Get(name string) (string, error) {
	args := []any{
		sql.Named("name", name),
		sql.Named("now", time.Now().UnixMilli()),
	}

	row, err := c.queryRow(`
	select
		value
	from
		keys
	where
		name = :name
		and (
			etime is null
			or etime > :now
		);
	`, args...)
	if err != nil {
		return "", fmt.Errorf("could not find key: %w", err)
	}

	var value string
	
	err = row.Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("could not read value: %w", err)
	}

	return value, nil
}

func (c *Client) exec(query string, args ...any) (sql.Result, error) {
	if _, ok := c.prepared[query]; !ok {
		statement, err := c.db.PrepareContext(c.context, query)
		if err != nil {
			return nil, fmt.Errorf("could not prepare statement: %w", err)
		}

		c.prepared[query] = statement
	}

	return c.prepared[query].ExecContext(c.context, args...)
}

func (c *Client) queryRow(query string, args ...any) (*sql.Row, error) {
	if _, ok := c.prepared[query]; !ok {
		statement, err := c.db.PrepareContext(c.context, query)
		if err != nil {
			return nil, fmt.Errorf("could not prepare statement: %w", err)
		}

		c.prepared[query] = statement
	}

	row := c.prepared[query].QueryRowContext(c.context, args...)
	
	return row, row.Err()
}

func (c *Client) Close() error {
	err := c.db.Close()
	if err != nil {
		return fmt.Errorf("could not close db: %w", err)
	}

	return nil
}
