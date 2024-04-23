package commands

import (
	"context"
	"log/slog"

	"github.com/jtarchie/sqlettuce/sdk"
	"github.com/tidwall/redcon"
)

func del(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	count := 0

	for _, name := range args[1:] {
		ok, err := client.Delete(context.TODO(), string(name))
		if err != nil {
			slog.Error("del", slog.String("error", err.Error()))
		}

		if ok {
			count++
		}
	}

	conn.WriteInt(count)
}
