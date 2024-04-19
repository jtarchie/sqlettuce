package sqlettus

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/tidwall/redcon"
)

type Server struct {
	addr   string
	client *Client
	wg     *sync.WaitGroup
	server *redcon.Server
}

func NewServer(addr string, client *Client) *Server {
	return &Server{
		addr:   addr,
		client: client,
		server: redcon.NewServer(addr,
			func(conn redcon.Conn, cmd redcon.Command) {
				switch strings.ToLower(string(cmd.Args[0])) {
				default:
					conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
				case "command", "info":
					conn.WriteString("OK")
				case "echo":
					for _, arg := range cmd.Args[1:] {
						conn.WriteBulk(arg)
						conn.WriteString(" ")
					}
				case "quit":
					conn.WriteString("OK")
					conn.Close()
				case "set":
					err := client.Set(string(cmd.Args[1]), string(cmd.Args[2]), 0)
					if err != nil {
						slog.Error("set", slog.String("error", err.Error()))
						conn.WriteError("could not set the key")
					} else {
						conn.WriteString("OK")
					}
				case "get":
					name := string(cmd.Args[1])
					value, err := client.Get(name)
					if err != nil {
						slog.Error("get", slog.String("error", err.Error()), slog.String("name", name))
						conn.WriteError("could not get the key")
					} else {
						conn.WriteBulkString(value)
					}
				}
			}, func(conn redcon.Conn) bool {
				slog.Debug("connection.accept", slog.String("client", conn.RemoteAddr()))
				return true
			}, func(conn redcon.Conn, err error) {
				if err != nil {
					slog.Debug("connection.close", slog.String("client", conn.RemoteAddr()), slog.String("error", err.Error()))
				} else {
					slog.Debug("connection.close", slog.String("client", conn.RemoteAddr()))
				}
			}),
		wg: &sync.WaitGroup{},
	}
}

func (s *Server) Start() {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		slog.Info("server.starting", slog.String("addr", s.addr))

		err := s.server.ListenAndServe()
		if err != nil {
			slog.Error("server.stop", slog.String("error", err.Error()))
		}
	}()
}

func (s *Server) Stop() error {
	slog.Info("server.stopping", slog.String("addr", s.server.Addr().String()))

	err := s.server.Close()
	if err != nil {
		return fmt.Errorf("could not stop server: %w", err)
	}

	slog.Info("server.stopped", slog.String("addr", s.server.Addr().String()))
	slog.Info("client.stopping")

	err = s.client.Close()
	if err != nil {
		return fmt.Errorf("could not stop client: %w", err)
	}

	slog.Info("client.stopped")

	return nil
}
