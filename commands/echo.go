package commands

import "github.com/tidwall/redcon"

func echo(args [][]byte, conn redcon.Conn) {
	for _, arg := range args[1:] {
		conn.WriteBulk(arg)
		conn.WriteString(" ")
	}
}
