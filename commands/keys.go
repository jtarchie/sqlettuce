package commands

import (
	"context"
	"log/slog"

	"github.com/jtarchie/sqlettuce/sdk"
	"github.com/tidwall/redcon"
)

func keys(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	if len(args) != 2 {
		conn.WriteError("need one argument for glob")

		return
	}

	keys, err := client.Keys(context.TODO(), string(args[1]))
	if err != nil {
		slog.Error("keys", slog.String("error", err.Error()))
		conn.WriteError("could not keys")
	} else {
		conn.WriteArray(len(keys))

		for _, key := range keys {
			conn.WriteBulkString(key)
		}
	}
}
