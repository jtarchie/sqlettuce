package commands

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jtarchie/sqlettus/sdk"
	"github.com/tidwall/redcon"
)

func keyType(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	keyType, err := client.Type(context.TODO(), string(args[1]))
	if errors.Is(err, sdk.ErrKeyDoesNotExist) {
		conn.WriteString("none")

		return
	}

	if err != nil {
		slog.Error("type", slog.String("error", err.Error()))
		conn.WriteString("none")

		return
	}

	conn.WriteString(keyType.String())
}
