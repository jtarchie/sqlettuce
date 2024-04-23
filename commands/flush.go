package commands

import (
	"context"
	"log/slog"

	"github.com/jtarchie/sqlettuce/sdk"
	"github.com/tidwall/redcon"
)

func flush(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	err := client.FlushDB(context.TODO())
	if err != nil {
		slog.Error(string(args[0]), slog.String("error", err.Error()))
		conn.WriteError("could not flush")
	} else {
		conn.WriteString("OK")
	}
}
