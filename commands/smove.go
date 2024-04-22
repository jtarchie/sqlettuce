package commands

import (
	"context"
	"log/slog"

	"github.com/jtarchie/sqlettus/sdk"
	"github.com/tidwall/redcon"
)

func smove(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	if len(args) < 4 {
		conn.WriteError("not enough arguments")

		return
	}

	source := string(args[1])
	destination := string(args[2])
	member := string(args[3])

	err := client.SetMove(context.TODO(), destination, source, member)
	if err != nil {
		slog.Error("smove", slog.String("error", err.Error()))
		conn.WriteInt64(0)
	} else {
		conn.WriteInt64(1)
	}
}
