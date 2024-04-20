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
	case "dbsize":
		dbsize(client, conn)
	case "del", "unlink":
		del(client, args, conn)
	case "echo":
		echo(args, conn)
	case "exists":
		exists(client, args, conn)
	case "expire":
		expire(client, args, conn)
	case "flushdb", "flushall":
		flush(client, args, conn)
	case "get":
		get(client, args, conn)
	case "keys":
		keys(client, args, conn)
	case "lpush":
		lpush(client, args, conn)
	case "mset":
		mset(client, args, conn)
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
	case "sort":
		sort(client, args, conn)
	case "ttl":
		ttl(client, args, conn, time.Second)
	case "type":
		keyType(client, args, conn)
	default:
		conn.WriteError("ERR unknown command '" + command + "'")
	}
}
