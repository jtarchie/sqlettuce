package commands

import (
	"context"
	"log/slog"

	"github.com/jtarchie/sqlettus/sdk"
	"github.com/tidwall/redcon"
)

func sinter(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	if len(args) < 2 {
		conn.WriteError("not enough arguments")

		return
	}

	names := make([]string, 0, len(args)-1)

	for _, name := range args[1:] {
		names = append(names, string(name))
	}

	values, err := client.SetIntersect(context.TODO(), names...)
	if err != nil {
		slog.Error("sinter", slog.String("error", err.Error()))
		conn.WriteError("unexpected error")
	} else {
		conn.WriteArray(len(values))

		for _, value := range values {
			conn.WriteBulkString(value)
		}
	}
}
