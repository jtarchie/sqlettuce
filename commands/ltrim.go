package commands

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/jtarchie/sqlettuce/sdk"
	"github.com/tidwall/redcon"
)

func ltrim(client *sdk.Client, args [][]byte, conn redcon.Conn) {
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

	err = client.ListTrim(context.TODO(), name, start, end)
	if err != nil {
		slog.Error("ltrim", slog.String("error", err.Error()))
		conn.WriteError("could not trim the list")
	} else {
		conn.WriteString("OK")
	}
}
