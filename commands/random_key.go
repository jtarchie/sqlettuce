package commands

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jtarchie/sqlettuce/sdk"
	"github.com/tidwall/redcon"
)

func randomKey(client *sdk.Client, conn redcon.Conn) {
	name, err := client.RandomKey(context.TODO())
	if errors.Is(err, sdk.ErrKeyDoesNotExist) {
		conn.WriteNull()

		return
	}

	if err != nil {
		slog.Error("randomkey", slog.String("error", err.Error()))
		conn.WriteError("could not get random key")

		return
	}

	conn.WriteBulkString(name)
}
