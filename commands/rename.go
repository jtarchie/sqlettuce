package commands

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jtarchie/sqlettus/sdk"
	"github.com/tidwall/redcon"
)

func rename(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	err := client.Rename(context.TODO(), string(args[1]), string(args[2]))
	if errors.Is(err, sdk.ErrKeyDoesNotExist) {
		conn.WriteError("key does not exist")

		return
	}

	if err != nil {
		slog.Error("rename", slog.String("error", err.Error()))
		conn.WriteError("could not rename")

		return
	}

	conn.WriteString("OK")
}
