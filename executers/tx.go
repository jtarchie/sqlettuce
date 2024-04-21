//nolint: sqlclosecheck, wrapcheck
package executers

import (
	"context"
	"database/sql"
)

type TxExecuter struct {
	tx       *sql.Tx
	prepared map[string]*sql.Stmt
}

func NewTX(tx *sql.Tx) *TxExecuter {
	return &TxExecuter{
		tx:       tx,
		prepared: map[string]*sql.Stmt{},
	}
}

func (t *TxExecuter) WithTX(_ context.Context, fun func(Executer) error) error {
	return fun(t)
}

func (t *TxExecuter) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return t.tx.PrepareContext(ctx, query)
}

func (t *TxExecuter) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if _, ok := t.prepared[query]; !ok {
		statement, err := t.tx.PrepareContext(ctx, query)
		if err != nil {
			return t.tx.ExecContext(ctx, query, args...)
		}

		t.prepared[query] = statement
	}

	return t.prepared[query].ExecContext(ctx, args...)
}

func (t *TxExecuter) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if _, ok := t.prepared[query]; !ok {
		statement, err := t.tx.PrepareContext(ctx, query)
		if err != nil {
			return t.tx.QueryContext(ctx, query, args...)
		}

		t.prepared[query] = statement
	}

	return t.prepared[query].QueryContext(ctx, args...)
}

func (t *TxExecuter) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if _, ok := t.prepared[query]; !ok {
		statement, err := t.tx.PrepareContext(ctx, query)
		if err != nil {
			return t.tx.QueryRowContext(ctx, query, args...)
		}

		t.prepared[query] = statement
	}

	return t.prepared[query].QueryRowContext(ctx, args...)
}

func (t *TxExecuter) Close() error {
	for _, prepared := range t.prepared {
		_ = prepared.Close()
	}

	return nil
}

var _ Executer = &TxExecuter{}
