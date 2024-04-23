package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/jtarchie/sqlettuce"
	"github.com/jtarchie/sqlettuce/executers"
	"github.com/jtarchie/sqlettuce/sdk"
)

type cli struct {
	Path string `default:":memory:"       help:"path to database file" require:"" short:"p"`
	Addr string `default:"localhost:6379" help:"server address"        short:"a"`
}

func (c *cli) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := executers.FromDB(c.Path)
	if err != nil {
		return fmt.Errorf("could not start db: %w", err)
	}

	client := sdk.NewClient(db)
	server := sqlettuce.NewServer(
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
