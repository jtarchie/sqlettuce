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
	case "decr":
		incr(client, args, conn, -1)
	case "decrby":
		incrBy(client, args, conn, -1)
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
	case "getset":
		getset(client, args, conn)
	case "keys":
		keys(client, args, conn)
	case "incr":
		incr(client, args, conn, 1)
	case "incrby":
		incrBy(client, args, conn, 1)
	case "lindex":
		lindex(client, args, conn)
	case "llen":
		llen(client, args, conn)
	case "lpop":
		lpop(client, args, conn)
	case "lpush":
		lpush(client, args, conn)
	case "lrange":
		lrange(client, args, conn)
	case "lrem":
		lrem(client, args, conn)
	case "lset":
		lset(client, args, conn)
	case "ltrim":
		ltrim(client, args, conn)
	case "mget":
		mget(client, args, conn)
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
	case "rpop":
		rpop(client, args, conn)
	case "rpush":
		rpush(client, args, conn)
	case "sadd":
		sadd(client, args, conn)
	case "scard":
		scard(client, args, conn)
	case "sdiff":
		sdiff(client, args, conn)
	case "sdiffstore":
		sdiffstore(client, args, conn)
	case "sinter":
		sinter(client, args, conn)
	case "sinterstore":
		sinterstore(client, args, conn)
	case "sismember":
		sismember(client, args, conn)
	case "smembers":
		smembers(client, args, conn)
	case "srandmember":
		srandmember(client, args, conn)
	case "srem":
		srem(client, args, conn)
	case "smove":
		smove(client, args, conn)
	case "set":
		set(client, args, conn)
	case "setnx":
		setnx(client, args, conn)
	case "sort":
		sort(client, args, conn)
	case "substr":
		substr(client, args, conn)
	case "ttl":
		ttl(client, args, conn, time.Second)
	case "type":
		keyType(client, args, conn)
	default:
		conn.WriteError("ERR unknown command '" + command + "'")
	}
}
