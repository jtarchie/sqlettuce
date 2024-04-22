//nolint: sqlclosecheck, wrapcheck
package executers

import (
	"context"
	"database/sql"

	csmap "github.com/mhmtszr/concurrent-swiss-map"
)

type TxExecuter struct {
	tx       *sql.Tx
	prepared *csmap.CsMap[string, *sql.Stmt]
}

func NewTX(tx *sql.Tx) *TxExecuter {
	return &TxExecuter{
		tx:       tx,
		prepared: csmap.Create[string, *sql.Stmt](),
	}
}

func (t *TxExecuter) WithTX(_ context.Context, fun func(Executer) error) error {
	return fun(t)
}

func (t *TxExecuter) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return t.tx.PrepareContext(ctx, query)
}

func (t *TxExecuter) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if ok := t.prepared.Has(query); !ok {
		statement, err := t.tx.PrepareContext(ctx, query)
		if err != nil {
			return t.tx.ExecContext(ctx, query, args...)
		}

		t.prepared.Store(query, statement)
	}

	v, _ := t.prepared.Load(query)

	return v.ExecContext(ctx, args...)
}

func (t *TxExecuter) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if ok := t.prepared.Has(query); !ok {
		statement, err := t.tx.PrepareContext(ctx, query)
		if err != nil {
			return t.tx.QueryContext(ctx, query, args...)
		}

		t.prepared.Store(query, statement)
	}

	v, _ := t.prepared.Load(query)

	return v.QueryContext(ctx, args...)
}

func (t *TxExecuter) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if ok := t.prepared.Has(query); !ok {
		statement, err := t.tx.PrepareContext(ctx, query)
		if err != nil {
			return t.tx.QueryRowContext(ctx, query, args...)
		}

		t.prepared.Store(query, statement)
	}

	v, _ := t.prepared.Load(query)

	return v.QueryRowContext(ctx, args...)
}

func (t *TxExecuter) Close() error {
	t.prepared.Range(func(_ string, value *sql.Stmt) bool {
		_ = value.Close()

		return false
	})

	return nil
}

var _ Executer = &TxExecuter{}
