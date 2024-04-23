package sqlettus_test

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jtarchie/sqlettus"
	"github.com/jtarchie/sqlettus/executers"
	"github.com/jtarchie/sqlettus/sdk"
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

	db, err := executers.FromDB(":memory:")
	if err != nil {
		t.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)
	server := sqlettus.NewServer(":6464", client)

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

	total, failed := 0, 0

	for _, test := range payload {
		total++

		if test.Skipped || test.Since != "1.0.0" || test.Tags == "cluster" {
			continue
		}

		err = client.FlushDB(context.TODO())
		if err != nil {
			t.Fatalf("could not flush db: %s", err)
		}

		for index, command := range test.Command {
			var args []interface{}
			for _, arg := range strings.Split(command, " ") {
				args = append(args, arg)
			}

			result, err := rdb.Do(context.TODO(), args...).Result()

			if err != nil && err.Error() != "redis: nil" {
				t.Logf("could not run test %q command %q: %s", test.Name, command, err)

				failed++

				break
			}

			contents, err := json.Marshal(result)
			if err != nil {
				t.Fatalf("could not marshal: %s", err)
			}

			var actual interface{}

			err = json.Unmarshal(contents, &actual)
			if err != nil {
				t.Fatalf("could not unmarshal: %s", err)
			}

			if diff := cmp.Diff(test.Result[index], actual); diff != "" {
				t.Logf("%q (-want +got):\n%s", test.Name, diff)

				failed++
			}
		}
	}

	if 0 < failed {
		t.Fatalf("failed = %d, success = %d", failed, total-failed)
	}
}
