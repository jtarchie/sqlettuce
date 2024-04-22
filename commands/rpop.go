package commands

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/jtarchie/sqlettus/sdk"
	"github.com/tidwall/redcon"
)

func rpop(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	if len(args) < 2 {
		conn.WriteError("not enough arguments")

		return
	}

	name := string(args[1])

	count := int64(1)
	if len(args) == 3 {
		var err error

		count, err = strconv.ParseInt(string(args[2]), 10, 64)
		if err != nil {
			conn.WriteError("not a valid number for count")

			return
		}
	}

	values, err := client.ListPop(context.TODO(), name, count)
	if err != nil {
		slog.Error("get", slog.String("error", err.Error()), slog.String("name", name))
		conn.WriteError("could not rpop the key")
	} else {
		if len(values) == 0 {
			conn.WriteNull()
		} else if len(values) == 1 {
			conn.WriteBulkString(values[0])
		} else {

			conn.WriteArray(len(values))

			for _, value := range values {
				conn.WriteBulkString(value)
			}
		}
	}
}
