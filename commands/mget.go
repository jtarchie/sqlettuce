package commands

import (
	"context"
	"log/slog"

	"github.com/jtarchie/sqlettus/sdk"
	"github.com/tidwall/redcon"
)

func mget(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	if len(args) < 2 {
		conn.WriteError("not enough arguments")

		return
	}

	names := make([]string, 0, len(args)-1)

	for _, name := range args[1:] {
		names = append(names, string(name))
	}

	values, err := client.MGet(context.TODO(), names...)
	if err != nil {
		slog.Error("mget", slog.String("error", err.Error()))
		conn.WriteError("could not get the keys")
	} else {
		conn.WriteArray(len(values))

		for _, value := range values {
			if value == nil {
				conn.WriteNull()
			} else {
				conn.WriteBulkString(*value)
			}
		}
	}
}
