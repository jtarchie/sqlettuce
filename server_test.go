package sqlettuce_test

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jtarchie/sqlettuce"
	"github.com/jtarchie/sqlettuce/executers"
	"github.com/jtarchie/sqlettuce/sdk"
	"github.com/redis/go-redis/v9"
)

type Compatibility []struct {
	Name          string   `json:"name"`
	Command       []string `json:"command"`
	Result        []any    `json:"result"`
	Since         string   `json:"since"`
	Tags          string   `json:"tags,omitempty"`
	CommandBinary bool     `json:"command_binary,omitempty"`
	Skipped       bool     `json:"skipped,omitempty"`
	SortResult    bool     `json:"sort_result,omitempty"`
}

func TestCompatibility(t *testing.T) {
	var payload Compatibility

	db, err := executers.FromDB("file:/data.db?vfs=memdb")
	if err != nil {
		t.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)
	server := sqlettuce.NewServer(":6464", client)

	server.Start()
	defer server.Close()

	contents, err := os.ReadFile("cts.json")
	if err != nil {
		t.Fatalf("could not open cts.json: %s", err)
	}

	err = json.Unmarshal(contents, &payload)
	if err != nil {
		t.Fatalf("could not unmarshal JSON: %s", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     ":6464",
		Protocol: 2,
	})

	for _, test := range payload {
		if test.Skipped || test.Since != "1.0.0" || test.Tags == "cluster" {
			continue
		}

		err = client.FlushDB(context.TODO())
		if err != nil {
			t.Fatalf("could not flush db: %s", err)
		}

		for index, command := range test.Command {
			t.Run(test.Name, func(t *testing.T) {
				var args []interface{}
				for _, arg := range strings.Split(command, " ") {
					args = append(args, arg)
				}

				result, err := rdb.Do(context.TODO(), args...).Result()

				if err != nil && err.Error() != "redis: nil" {
					t.Errorf("command=%q err=%s", command, err)
				}

				contents, err := json.Marshal(result)
				if err != nil {
					t.Errorf("could not marshal: %s", err)
				}

				var actual interface{}

				err = json.Unmarshal(contents, &actual)
				if err != nil {
					t.Errorf("could not unmarshal: %s", err)
				}

				if diff := cmp.Diff(test.Result[index], actual); diff != "" {
					t.Errorf("%q (-want +got):\n%s", command, diff)
				}
			})
		}
	}
}

func XBenchmarkCompatibility(b *testing.B) {
	var payload Compatibility

	db, err := executers.FromDB("file:/data.db?vfs=memdb")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)
	server := sqlettuce.NewServer(":6464", client)

	server.Start()
	defer server.Close()

	contents, err := os.ReadFile("cts.json")
	if err != nil {
		b.Fatalf("could not open cts.json: %s", err)
	}

	err = json.Unmarshal(contents, &payload)
	if err != nil {
		b.Fatalf("could not unmarshal JSON: %s", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     ":6464",
		Protocol: 2,
	})

	for _, test := range payload {
		if test.Skipped || test.Since != "1.0.0" || test.Tags == "cluster" {
			continue
		}

		err = client.FlushDB(context.TODO())
		if err != nil {
			b.Fatalf("could not flush db: %s", err)
		}

		for index, command := range test.Command {
			b.Run(test.Name, func(b *testing.B) {
				var args []interface{}
				for _, arg := range strings.Split(command, " ") {
					args = append(args, arg)
				}

				result, err := rdb.Do(context.TODO(), args...).Result()

				if err != nil && err.Error() != "redis: nil" {
					b.Fatalf("command=%q err=%s", command, err)
				}

				contents, err := json.Marshal(result)
				if err != nil {
					b.Fatalf("could not marshal: %s", err)
				}

				var actual interface{}

				err = json.Unmarshal(contents, &actual)
				if err != nil {
					b.Fatalf("could not unmarshal: %s", err)
				}

				if diff := cmp.Diff(test.Result[index], actual); diff != "" {
					b.Fatalf("%q (-want +got):\n%s", command, diff)
				}
			})
		}
	}
}
