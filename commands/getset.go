package commands

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jtarchie/sqlettus/sdk"
	"github.com/tidwall/redcon"
)

func getset(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	if len(args) < 3 {
		conn.WriteError("Not enough args")

		return
	}

	name := string(args[1])
	value := string(args[2])

	value, err := client.GetSet(context.TODO(), name, value)

	switch {
	case errors.Is(err, sdk.ErrKeyDoesNotExist):
		conn.WriteNull()
	case err != nil:
		slog.Error("getset", slog.String("error", err.Error()))
		conn.WriteError("could not set the key")
	default:
		conn.WriteBulkString(value)
	}
}
