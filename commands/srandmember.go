package commands

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/jtarchie/sqlettuce/sdk"
	"github.com/tidwall/redcon"
)

func srandmember(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	var err error

	if len(args) < 2 {
		conn.WriteError("not enough arguments")

		return
	}

	name := string(args[1])
	length := int64(1)

	if len(args) > 2 {
		length, err = strconv.ParseInt(string(args[2]), 10, 64)
		if err != nil {
			conn.WriteError("please provide a valid integer")

			return
		}
	}

	values, err := client.SetRandomMember(context.TODO(), name, length)
	if err != nil {
		slog.Error("spop", slog.String("error", err.Error()))
		conn.WriteNull()
	}

	if len(values) == 1 {
		conn.WriteBulkString(values[0])

		return
	}

	conn.WriteArray(len(values))

	for _, value := range values {
		conn.WriteBulkString(value)
	}
}
