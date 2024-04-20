package commands

import (
	"context"
	"log/slog"

	"github.com/jtarchie/sqlettus/sdk"
	"github.com/tidwall/redcon"
)

func mset(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	if len(args) < 3 || len(args[1:])%2 != 0 {
		conn.WriteError("even number of args required (key value pairs)")

		return
	}

	pairs := make([][2]string, 0, len(args[1:])/2)
	for i := 1; i < len(args); i += 2 {
		pairs = append(pairs, [2]string{string(args[i]), string(args[i+1])})
	}

	err := client.MSet(context.TODO(), pairs...)
	if err != nil {
		slog.Error("mset", slog.String("error", err.Error()))
		conn.WriteError("could not mset the key(s)")
	} else {
		conn.WriteString("OK")
	}
}
