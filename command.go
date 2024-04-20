package sqlettus

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/tidwall/redcon"
)

func handleCommand(
	client *Client,
	conn redcon.Conn,
	args [][]byte,
) {
	switch strings.ToLower(string(args[0])) {
	default:
		conn.WriteError("ERR unknown command '" + string(args[0]) + "'")
	case "command", "info":
		conn.WriteString("OK")
	case "echo":
		for _, arg := range args[1:] {
			conn.WriteBulk(arg)
			conn.WriteString(" ")
		}
	case "quit":
		conn.WriteString("OK")
		conn.Close()
	case "rename":
		err := client.Rename(string(args[1]), string(args[2]))
		if errors.Is(err, ErrKeyDoesNotExist) {
			conn.WriteError("key does not exist")
			return
		}

		if err != nil {
			slog.Error("rename", slog.String("error", err.Error()))
			conn.WriteError("could not rename")
			return
		}

		conn.WriteString("OK")

	case "renamenx":
		err := client.RenameIfNotExists(string(args[1]), string(args[2]))
		if errors.Is(err, ErrKeyAlreadyExists) {
			conn.WriteInt(0)
			return
		}

		if err != nil {
			slog.Error("renamenx", slog.String("error", err.Error()))
			conn.WriteError("could not renamenx")
			return
		}
		
		conn.WriteInt(1)
	case "del", "unlink":
		count := 0
		for _, name := range args[1:] {
			ok, err := client.Delete(string(name))
			if err != nil {
				slog.Error("del", slog.String("error", err.Error()))
			}
			if ok {
				count++
			}
		}
		conn.WriteInt(count)
	case "exists":
		count := 0
		for _, name := range args[1:] {
			ok, err := client.Exists(string(name))
			if err != nil {
				slog.Error("del", slog.String("error", err.Error()))
			}
			if ok {
				count++
			}
		}
		conn.WriteInt(count)
	case "set":
		err := client.Set(string(args[1]), string(args[2]), 0)
		if err != nil {
			slog.Error("set", slog.String("error", err.Error()))
			conn.WriteError("could not set the key")
		} else {
			conn.WriteString("OK")
		}
	case "get":
		name := string(args[1])
		value, err := client.Get(name)
		if err != nil {
			slog.Error("get", slog.String("error", err.Error()), slog.String("name", name))
			conn.WriteError("could not get the key")
		} else {
			conn.WriteBulkString(value)
		}
	case "randomkey":
		name, err := client.RandomKey()
		if errors.Is(err, ErrKeyDoesNotExist) {
			conn.WriteNull()
			return
		}

		if err != nil {
			slog.Error("randomkey", slog.String("error", err.Error()))
			conn.WriteError("could not get random key")
			return
		}

		conn.WriteBulkString(name)
	}
}