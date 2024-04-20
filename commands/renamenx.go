package commands

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jtarchie/sqlettus/sdk"
	"github.com/tidwall/redcon"
)

func renamenx(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	err := client.RenameIfNotExists(context.TODO(), string(args[1]), string(args[2]))
	if errors.Is(err, sdk.ErrKeyAlreadyExists) {
		conn.WriteInt(0)

		return
	}

	if err != nil {
		slog.Error("renamenx", slog.String("error", err.Error()))
		conn.WriteError("could not renamenx")

		return
	}

	conn.WriteInt(1)
}
