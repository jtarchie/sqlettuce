package commands

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/jtarchie/sqlettuce/sdk"
	"github.com/tidwall/redcon"
)

func substr(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	if len(args) < 4 {
		conn.WriteError("not enough arguments")

		return
	}

	name := string(args[1])

	start, err := strconv.ParseInt(string(args[2]), 10, 64)
	if err != nil {
		conn.WriteError("expected integer value for start")

		return
	}

	end, err := strconv.ParseInt(string(args[3]), 10, 64)
	if err != nil {
		conn.WriteError("expected integer value for start")

		return
	}

	value, err := client.Substr(context.TODO(), name, start, end)
	if err != nil {
		slog.Error("substr", slog.String("error", err.Error()))
		conn.WriteError("could not substr")
	} else {
		conn.WriteBulkString(value)
	}
}
