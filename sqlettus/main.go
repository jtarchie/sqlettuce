package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/jtarchie/sqlettus"
	"github.com/jtarchie/sqlettus/sdk"
)

type cli struct {
	Path string `default:":memory:"       help:"path to database file" require:"" short:"p"`
	Addr string `default:"localhost:6379" help:"server address"        short:"a"`
}

func (c *cli) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client, err := sdk.NewClient(ctx, c.Path)
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
