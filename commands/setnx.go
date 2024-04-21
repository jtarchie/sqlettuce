package commands

import (
	"context"
	"log/slog"

	"github.com/jtarchie/sqlettus/sdk"
	"github.com/tidwall/redcon"
)

func setnx(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	if len(args) < 3 {
		conn.WriteError("Not enough args")

		return
	}

	err := client.SetIfNotExists(context.TODO(), string(args[1]), string(args[2]))
	if err != nil {
		slog.Error("set", slog.String("error", err.Error()))
		conn.WriteInt64(0)
	} else {
		conn.WriteInt64(1)
	}
}
