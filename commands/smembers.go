package commands

import (
	"context"
	"log/slog"

	"github.com/jtarchie/sqlettuce/sdk"
	"github.com/tidwall/redcon"
)

func smembers(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	if len(args) < 2 {
		conn.WriteError("not enough arguments")

		return
	}

	name := string(args[1])

	values, err := client.SetMembers(context.TODO(), name)
	if err != nil {
		slog.Error("smembers", slog.String("error", err.Error()))
		conn.WriteError("could not get members")
	} else {
		conn.WriteArray(len(values))

		for _, value := range values {
			conn.WriteBulkString(value)
		}
	}
}
