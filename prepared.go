package sqlettus

import (
	"context"
	"database/sql"
	"fmt"
)

type preparedExecuter struct {
	db       *sql.DB
	prepared map[string]*sql.Stmt
}

func (p *preparedExecuter) WithTX(ctx context.Context, fun func(Executer) error) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	err = fun(&txExecuter{tx})
	if err != nil {
		return fmt.Errorf("could execute within transaction: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func (p *preparedExecuter) PrepareContext(context.Context, string) (*sql.Stmt, error) {
	return nil, fmt.Errorf("PreparedContext unsupported")
}

func (p *preparedExecuter) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if _, ok := p.prepared[query]; !ok {
		statement, err := p.db.PrepareContext(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("could not prepare statement: %w", err)
		}

		p.prepared[query] = statement
	}

	return p.prepared[query].ExecContext(ctx, args...)
}

func (p *preparedExecuter) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if _, ok := p.prepared[query]; !ok {
		statement, err := p.db.PrepareContext(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("could not prepare statement: %w", err)
		}

		p.prepared[query] = statement
	}

	return p.prepared[query].QueryContext(ctx, args...)
}

func (p *preparedExecuter) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if _, ok := p.prepared[query]; !ok {
		statement, err := p.db.PrepareContext(ctx, query)
		if err != nil {
			return nil
		}

		p.prepared[query] = statement
	}

	return p.prepared[query].QueryRowContext(ctx, args...)
}

func (p *preparedExecuter) Close() error {
	return p.db.Close()
}

var _ Executer = &preparedExecuter{}
