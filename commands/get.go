package commands

import (
	"context"
	"log/slog"

	"github.com/jtarchie/sqlettus/sdk"
	"github.com/tidwall/redcon"
)

func get(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	name := string(args[1])

	value, err := client.Get(context.TODO(), name)
	if err != nil {
		slog.Error("get", slog.String("error", err.Error()), slog.String("name", name))
		conn.WriteError("could not get the key")
	} else {
		if value == nil {
			conn.WriteNull()
		} else {
			conn.WriteBulkString(*value)
		}
	}
}
