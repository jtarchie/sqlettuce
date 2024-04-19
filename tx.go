package sqlettus

import (
	"context"
	"database/sql"
)

type txExecuter struct {
	*sql.Tx
}

func (t *txExecuter) WithTX(_ context.Context, fun func(Executer) error) error {
	return fun(t)
}

func (t *txExecuter) Close() error {
	return nil
}

var _ Executer = &txExecuter{}