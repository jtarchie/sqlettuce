package commands

import (
	"context"
	"log/slog"

	"github.com/jtarchie/sqlettuce/sdk"
	"github.com/tidwall/redcon"
)

func llen(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	if len(args) < 2 {
		conn.WriteError("not enough arguments")

		return
	}

	name := string(args[1])

	length, err := client.ListLength(context.TODO(), name)
	if err != nil {
		slog.Error("llen", slog.String("error", err.Error()))
		conn.WriteError("could not determine length")
	} else {
		conn.WriteInt64(length)
	}
}
