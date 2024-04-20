package commands

import (
	"log/slog"

	"github.com/jtarchie/sqlettus/sdk"
	"github.com/tidwall/redcon"
)

func set(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	err := client.Set(string(args[1]), string(args[2]), 0)
	if err != nil {
		slog.Error("set", slog.String("error", err.Error()))
		conn.WriteError("could not set the key")
	} else {
		conn.WriteString("OK")
	}
}
