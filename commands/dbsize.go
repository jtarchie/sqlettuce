package commands

import (
	"context"
	"log/slog"

	"github.com/jtarchie/sqlettus/sdk"
	"github.com/tidwall/redcon"
)

func dbsize(client *sdk.Client, conn redcon.Conn) {
	count, err := client.DBSize(context.TODO())
	if err != nil {
		slog.Error("dbsize", slog.String("error", err.Error()))
		conn.WriteError("could not determine dbsize")
	} else {
		conn.WriteInt64(count)
	}
}
