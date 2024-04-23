package commands

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/jtarchie/sqlettuce/sdk"
	"github.com/tidwall/redcon"
)

func lrem(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	if len(args) < 4 {
		conn.WriteError("not enough arguments")

		return
	}

	name := string(args[1])
	element := string(args[3])

	count, err := strconv.ParseInt(string(args[2]), 10, 64)
	if err != nil {
		conn.WriteError("expected integer value for start")

		return
	}

	final, err := client.ListRemove(context.TODO(), name, count, element)
	if err != nil {
		slog.Error("lrem", slog.String("error", err.Error()))
		conn.WriteError("could not lrem")
	} else {
		conn.WriteInt64(final)
	}
}
