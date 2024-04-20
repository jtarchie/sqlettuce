package commands

import "github.com/tidwall/redcon"

func quit(conn redcon.Conn) {
	conn.WriteString("OK")
	conn.Close()
}
