package commands

import (
	"errors"
	"log/slog"
	"time"

	"github.com/jtarchie/sqlettus/sdk"
	"github.com/tidwall/redcon"
)

func ttl(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	ttl, err := client.TTL(string(args[1]))

	if errors.Is(err, sdk.ErrKeyDoesNotExist) {
		conn.WriteInt(-2)

		return
	}

	if err != nil {
		slog.Error("ttl", slog.String("error", err.Error()))
		conn.WriteError("could not find ttl")

		return
	}

	if ttl == nil {
		conn.WriteInt(-1)

		return
	}

	conn.WriteInt64(*ttl / int64(time.Second))
}
