package sqlettus

import (
	"context"
	"database/sql"
	"fmt"
)

type preparedExecuter struct {
	executer Executer
	prepared map[string]*sql.Stmt
}

// PrepareContext implements Executer.
func (p *preparedExecuter) PrepareContext(context.Context, string) (*sql.Stmt, error) {
	return nil, fmt.Errorf("PreparedContext unsupported")
}

// ExecContext implements Executer.
func (p *preparedExecuter) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if _, ok := p.prepared[query]; !ok {
		statement, err := p.executer.PrepareContext(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("could not prepare statement: %w", err)
		}

		p.prepared[query] = statement
	}

	return p.prepared[query].ExecContext(ctx, args...)
}

// QueryContext implements Executer.
func (p *preparedExecuter) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if _, ok := p.prepared[query]; !ok {
		statement, err := p.executer.PrepareContext(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("could not prepare statement: %w", err)
		}

		p.prepared[query] = statement
	}

	return p.prepared[query].QueryContext(ctx, args...)
}

// QueryRowContext implements Executer.
func (p *preparedExecuter) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if _, ok := p.prepared[query]; !ok {
		statement, err := p.executer.PrepareContext(ctx, query)
		if err != nil {
			return nil
		}

		p.prepared[query] = statement
	}

	return p.prepared[query].QueryRowContext(ctx, args...)
}

var _ Executer = &preparedExecuter{}
