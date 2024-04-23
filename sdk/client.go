package sdk

import (
	_ "embed"
	"errors"
	"fmt"

	"github.com/jtarchie/sqlettuce/executers"
	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrKeyAlreadyExists = errors.New("key already exists")
	ErrKeyDoesNotExist  = errors.New("key does not exist")
)

type Client struct {
	db executers.Executer
}

func NewClient(db executers.Executer) *Client {
	return &Client{
		db: db,
	}
}

func (c *Client) Close() error {
	err := c.db.Close()
	if err != nil {
		return fmt.Errorf("could not close db: %w", err)
	}

	return nil
}
