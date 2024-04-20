package commands

import (
	"context"
	"log/slog"

	"github.com/jtarchie/sqlettus/sdk"
	"github.com/tidwall/redcon"
)

func exists(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	count := 0

	for _, name := range args[1:] {
		ok, err := client.Exists(context.TODO(), string(name))
		if err != nil {
			slog.Error("del", slog.String("error", err.Error()))
		}

		if ok {
			count++
		}
	}

	conn.WriteInt(count)
}
