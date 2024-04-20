package commands

import (
	"strings"
	"time"

	"github.com/jtarchie/sqlettus/sdk"
	"github.com/tidwall/redcon"
)

func Handle(
	client *sdk.Client,
	conn redcon.Conn,
	args [][]byte,
) {
	command := string(args[0])
	switch strings.ToLower(command) {
	case "command", "info":
		conn.WriteString("OK")
	case "del", "unlink":
		del(client, args, conn)
	case "echo":
		echo(args, conn)
	case "exists":
		exists(client, args, conn)
	case "expire":
		expire(client, args, conn)
	case "get":
		get(client, args, conn)
	case "pttl":
		ttl(client, args, conn, time.Millisecond)
	case "quit":
		quit(conn)
	case "randomkey":
		randomKey(client, conn)
	case "rename":
		rename(client, args, conn)
	case "renamenx":
		renamenx(client, args, conn)
	case "set":
		set(client, args, conn)
	case "ttl":
		ttl(client, args, conn, time.Second)
	default:
		conn.WriteError("ERR unknown command '" + command + "'")
	}
}
