package commands

import (
	"context"
	"log/slog"

	"github.com/jtarchie/sqlettus/sdk"
	"github.com/tidwall/redcon"
)

func srem(client *sdk.Client, args [][]byte, conn redcon.Conn) {
	if len(args) < 3 {
		conn.WriteError("not enough arguments")

		return
	}

	name := string(args[1])
	members := make([]string, 0, len(args)-2)

	for _, member := range args[2:] {
		members = append(members, string(member))
	}

	count, err := client.SetRemove(context.TODO(), name, members...)
	if err != nil {
		slog.Error("srem", slog.String("error", err.Error()))
		conn.WriteError("could not remove elements from set")
	} else {
		conn.WriteInt64(count)
	}
}
