package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	sqlettus "github.com/jtarchie/sqlettuce"
)

type cli struct {
	Path string `short:"p" help:"path to database file" default:":memory:" require:""`
	Addr string `short:"a" help:"server address" default:"localhost:6379"`
}

func (c *cli) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client, err := sqlettus.NewClient(ctx, c.Path)
	if err != nil {
		return fmt.Errorf("could not start client: %w", err)
	}

	server := sqlettus.NewServer(
		c.Addr,
		client,
	)

	server.Start()

	<-ctx.Done()

	if err := server.Close(); err != nil {
		return fmt.Errorf("server stop errored: %w", err)
	}

	return nil
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))

	cli := &cli{}
	ctx := kong.Parse(cli)
	// Call the Run() method of the selected parsed command.
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
