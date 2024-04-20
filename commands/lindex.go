package commands

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/jtarchie/sqlettus/sdk"
	"github.com/tidwall/redcon"
)

func lindex(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	if len(args) < 3 {
		conn.WriteError("not enough arguments")

		return
	}

	name := string(args[1])

	index, err := strconv.ParseInt(string(args[2]), 10, 64)
	if err != nil {
		conn.WriteError("expected integer value for index")

		return
	}

	value, err := client.ListAt(context.TODO(), name, index)
	if err != nil {
		slog.Error("lindex", slog.String("error", err.Error()))
		conn.WriteError("could not determine value")
	} else {
		conn.WriteBulkString(value)
	}
}
