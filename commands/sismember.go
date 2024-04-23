package commands

import (
	"context"
	"log/slog"

	"github.com/jtarchie/sqlettuce/sdk"
	"github.com/tidwall/redcon"
)

func sismember(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	if len(args) < 3 {
		conn.WriteError("not enough arguments")

		return
	}

	name := string(args[1])
	element := string(args[2])

	found, err := client.SetContains(context.TODO(), name, element)
	if err != nil {
		slog.Error("sismember", slog.String("error", err.Error()))
		conn.WriteError("could not find element")
	} else {
		if found {
			conn.WriteInt64(1)
		} else {
			conn.WriteInt64(0)
		}
	}
}
