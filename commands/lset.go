package commands

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/jtarchie/sqlettuce/sdk"
	"github.com/tidwall/redcon"
)

func lset(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	if len(args) < 4 {
		conn.WriteError("not enough arguments")

		return
	}

	name := string(args[1])
	element := string(args[3])

	index, err := strconv.ParseInt(string(args[2]), 10, 64)
	if err != nil {
		conn.WriteError("expected integer value for start")

		return
	}

	err = client.ListSet(context.TODO(), name, index, element)
	if err != nil {
		slog.Error("lset", slog.String("error", err.Error()))
		conn.WriteError("could not set element in list")
	} else {
		conn.WriteString("OK")
	}
}
