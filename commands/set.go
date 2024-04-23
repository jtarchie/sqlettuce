package commands

import (
	"context"
	"log/slog"

	"github.com/jtarchie/sqlettuce/sdk"
	"github.com/tidwall/redcon"
)

func set(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	if len(args) < 3 {
		conn.WriteError("Not enough args")

		return
	}

	err := client.Set(context.TODO(), string(args[1]), string(args[2]), 0)
	if err != nil {
		slog.Error("set", slog.String("error", err.Error()))
		conn.WriteError("could not set the key")
	} else {
		conn.WriteString("OK")
	}
}
