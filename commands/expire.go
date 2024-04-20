package commands

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"github.com/jtarchie/sqlettus/sdk"
	"github.com/tidwall/redcon"
)

func expire(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	seconds, err := strconv.Atoi(string(args[2]))
	if err != nil {
		conn.WriteError("not a valid number for seconds")

		return
	}

	duration := time.Second * time.Duration(seconds)

	err = client.Expire(context.TODO(), string(args[1]), duration)

	if errors.Is(err, sdk.ErrKeyDoesNotExist) {
		conn.WriteInt(0)

		return
	}

	if err != nil {
		slog.Error("ttl", slog.String("error", err.Error()))
		conn.WriteError("could not find ttl")

		return
	}

	conn.WriteInt(1)
}
