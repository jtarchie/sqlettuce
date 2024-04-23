package commands

import (
	"context"
	"log/slog"

	"github.com/jtarchie/sqlettuce/sdk"
	"github.com/tidwall/redcon"
)

func scard(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	if len(args) < 2 {
		conn.WriteError("not enough arguments")

		return
	}

	name := string(args[1])

	length, err := client.SetLength(context.TODO(), name)
	if err != nil {
		slog.Error("scard", slog.String("error", err.Error()))
		conn.WriteError("could not determine length")
	} else {
		conn.WriteInt64(length)
	}
}
