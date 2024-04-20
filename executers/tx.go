package executers

import (
	"context"
	"database/sql"
)

type TxExecuter struct {
	*sql.Tx
}

func (t *TxExecuter) WithTX(_ context.Context, fun func(Executer) error) error {
	return fun(t)
}

func (t *TxExecuter) Close() error {
	return nil
}

var _ Executer = &TxExecuter{}
