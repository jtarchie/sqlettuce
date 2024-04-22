package commands

import (
	"context"
	"log/slog"

	"github.com/jtarchie/sqlettus/sdk"
	"github.com/tidwall/redcon"
)

func sinterstore(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	if len(args) < 3 {
		conn.WriteError("not enough arguments")

		return
	}

	name := string(args[1])
	names := make([]string, 0, len(args)-2)

	for _, name := range args[2:] {
		names = append(names, string(name))
	}

	count, err := client.SetIntersectAndStore(context.TODO(), name, names...)
	if err != nil {
		slog.Error("sinterstore", slog.String("error", err.Error()))
		conn.WriteError("unexpected error")
	} else {
		conn.WriteInt64(count)
	}
}
