package sqlettus

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrKeyAlreadyExists = errors.New("key already exists")
	ErrKeyDoesNotExist  = errors.New("key does not exist")
)

//go:embed schema.sql
var schemaSQL string

type Executer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	WithTX(ctx context.Context, fun func(Executer) error) error
	Close() error
}

type Client struct {
	context context.Context
	db      Executer
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
		db: &preparedExecuter{
			db:       db,
			prepared: map[string]*sql.Stmt{},
		},
		context: ctx,
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

	_, err := c.db.ExecContext(c.context, `
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
		return fmt.Errorf("could not set key: %w", err)
	}

	return nil
}

func (c *Client) Get(name string) (string, error) {
	args := []any{
		sql.Named("name", name),
		sql.Named("now", time.Now().UnixMilli()),
	}

	row := c.db.QueryRowContext(c.context, `
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

	err := row.Err()
	if err != nil {
		return "", fmt.Errorf("could not find key: %w", err)
	}

	var value string

	err = row.Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}

	if err != nil {
		return "", fmt.Errorf("could not read value: %w", err)
	}

	return value, nil
}

func (c *Client) Delete(name string) (bool, error) {
	args := []any{
		sql.Named("name", name),
	}

	_, err := c.db.ExecContext(
		c.context,
		`DELETE FROM keys where name = :name;`,
		args...,
	)
	if err != nil {
		return false, fmt.Errorf("could not delete key: %w", err)
	}

	return true, nil
}

func (c *Client) Rename(current string, next string) error {
	err := c.db.WithTX(c.context, func(tx Executer) error {
		_, _ = tx.ExecContext(c.context, `DELETE FROM keys WHERE name = :new`, sql.Named("new", next))

		result, err := tx.ExecContext(c.context, `UPDATE keys SET name = :new WHERE name = :old`, sql.Named("new", next), sql.Named("old", current))
		if err != nil {
			return fmt.Errorf("could not rename key: %w", err)
		}

		count, _ := result.RowsAffected()
		if count == 0 {
			return ErrKeyDoesNotExist
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("could not rename: %w", err)
	}

	return nil
}

func (c *Client) RenameIfNotExists(current string, next string) error {
	err := c.db.WithTX(c.context, func(tx Executer) error {
		row := tx.QueryRowContext(c.context, `SELECT 1 FROM keys WHERE name = :new`, sql.Named("new", next))
		if row.Err() != nil {
			return fmt.Errorf("could not find new: %w", row.Err())
		}

		var value int

		err := row.Scan(&value)
		if errors.Is(err, sql.ErrNoRows) {
			_, err = tx.ExecContext(c.context, `UPDATE keys SET name = :new WHERE name = :old`, sql.Named("new", next), sql.Named("old", current))
			if err != nil {
				return fmt.Errorf("could not rename key: %w", err)
			}
		}

		if err != nil {
			return fmt.Errorf("could not scan: %w", err)
		}

		if value == 1 {
			return ErrKeyAlreadyExists
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("could not rename: %w", err)
	}

	return nil
}

func (c *Client) RandomKey() (string, error) {
	row := c.db.QueryRowContext(
		c.context,
		`SELECT name FROM keys ORDER BY RANDOM() LIMIT 1`,
	)
	if row.Err() != nil {
		return "", fmt.Errorf("could select random value: %w", row.Err())
	}

	var name string

	err := row.Scan(&name)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrKeyDoesNotExist
	}

	if err != nil {
		return "", fmt.Errorf("could not scan value: %w", err)
	}

	return name, nil
}

func (c *Client) FlushDB() error {
	_, err := c.db.ExecContext(c.context, `
		DELETE FROM keys;
		VACUUM;
		PRAGMA OPTIMIZE;
	`)
	if err != nil {
		return fmt.Errorf("could not flush db: %w", err)
	}

	return nil
}

func (c *Client) Exists(name string) (bool, error) {
	row := c.db.QueryRowContext(
		c.context,
		`SELECT 1 FROM keys WHERE name = :name`,
		sql.Named("name", name),
	)
	if row.Err() != nil {
		return false, fmt.Errorf("could not read key: %w", row.Err())
	}

	var value int

	err := row.Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	return true, nil
}

func (c *Client) Close() error {
	err := c.db.Close()
	if err != nil {
		return fmt.Errorf("could not close db: %w", err)
	}

	return nil
}
