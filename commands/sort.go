package commands

import (
	"context"
	"log/slog"

	"github.com/jtarchie/sqlettus/sdk"
	"github.com/tidwall/redcon"
)

func sort(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	name := string(args[1])

	values, err := client.Sort(context.TODO(), name)
	if err != nil {
		slog.Error("sort", slog.String("error", err.Error()))
		conn.WriteError("could not sort elements")

		return
	}

	conn.WriteArray(len(values))

	for _, value := range values {
		conn.WriteBulkString(value)
	}
}
