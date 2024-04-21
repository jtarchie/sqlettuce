//nolint: sqlclosecheck, wrapcheck
package executers

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
)

//go:embed schema.sql
var schemaSQL string

type PreparedExecuter struct {
	db       *sql.DB
	prepared map[string]*sql.Stmt
}

func NewPrepared(db *sql.DB) *PreparedExecuter {
	return &PreparedExecuter{
		db:       db,
		prepared: map[string]*sql.Stmt{},
	}
}

func FromDB(filename string) (*PreparedExecuter, error) {
	// set config based on: https://github.com/mattn/go-sqlite3/issues/1179#issuecomment-1638083995
	db, err := sql.Open("sqlite3", filename+"?_tx=immediate")
	if err != nil {
		return nil, fmt.Errorf("could open sqlite3: %w", err)
	}

	_, err = db.ExecContext(context.TODO(), schemaSQL)
	if err != nil {
		return nil, fmt.Errorf("could not create schema: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)
	db.SetConnMaxIdleTime(0)

	return NewPrepared(db), nil
}

func (p *PreparedExecuter) WithTX(ctx context.Context, fun func(Executer) error) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	defer func() { _ = tx.Rollback() }()

	err = fun(NewTX(tx))
	if err != nil {
		return fmt.Errorf("could execute within transaction: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func (p *PreparedExecuter) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return p.db.PrepareContext(ctx, query)
}

func (p *PreparedExecuter) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if _, ok := p.prepared[query]; !ok {
		statement, err := p.db.PrepareContext(ctx, query)
		if err != nil {
			return p.db.ExecContext(ctx, query, args...)
		}

		p.prepared[query] = statement
	}

	return p.prepared[query].ExecContext(ctx, args...)
}

func (p *PreparedExecuter) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if _, ok := p.prepared[query]; !ok {
		statement, err := p.db.PrepareContext(ctx, query)
		if err != nil {
			return p.db.QueryContext(ctx, query, args...)
		}

		p.prepared[query] = statement
	}

	return p.prepared[query].QueryContext(ctx, args...)
}

func (p *PreparedExecuter) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if _, ok := p.prepared[query]; !ok {
		statement, err := p.db.PrepareContext(ctx, query)
		if err != nil {
			return p.db.QueryRowContext(ctx, query, args...)
		}

		p.prepared[query] = statement
	}

	return p.prepared[query].QueryRowContext(ctx, args...)
}

func (p *PreparedExecuter) Close() error {
	for _, prepared := range p.prepared {
		_ = prepared.Close()
	}

	return p.db.Close()
}

var _ Executer = &PreparedExecuter{}
