package commands

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/jtarchie/sqlettus/sdk"
	"github.com/tidwall/redcon"
)

func incr(client *sdk.Client, args [][]byte, conn redcon.Conn, direction int64) {
	var err error

	name := string(args[1])
	add := int64(1)

	value, err := client.AddInt(context.TODO(), name, add*direction)
	if err != nil {
		slog.Error("del", slog.String("error", err.Error()))
	}

	conn.WriteInt64(value)
}

func incrBy(client *sdk.Client, args [][]byte, conn redcon.Conn, direction int64) {
	var err error

	name := string(args[1])
	add := int64(1)

	if len(args) > 2 {
		add, err = strconv.ParseInt(string(args[2]), 10, 64)
		if err != nil {
			conn.WriteError("please provide a valid integer")

			return
		}
	}

	value, err := client.AddInt(context.TODO(), name, add*direction)
	if err != nil {
		slog.Error("del", slog.String("error", err.Error()))
	}

	conn.WriteInt64(value)
}
