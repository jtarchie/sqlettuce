package commands

import (
	"context"
	"log/slog"

	"github.com/jtarchie/sqlettuce/sdk"
	"github.com/tidwall/redcon"
)

func rpush(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	if len(args) < 3 {
		conn.WriteError("not enough arguments")

		return
	}

	name := string(args[1])
	values := make([]string, 0, len(args)-2)

	for _, value := range args[2:] {
		values = append(values, string(value))
	}

	length, err := client.ListPush(context.TODO(), name, values...)
	if err != nil {
		slog.Error("lpush", slog.String("error", err.Error()))
		conn.WriteError("unexpected error")
	} else {
		conn.WriteInt64(length)
	}
}
