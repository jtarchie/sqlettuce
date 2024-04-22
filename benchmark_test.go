package sqlettus_test

import (
	"context"
	"testing"

	"github.com/jtarchie/sqlettus/executers"
	"github.com/jtarchie/sqlettus/sdk"
)

func BenchmarkMSet(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = client.MSet(context.TODO(), [][2]string{{"a", "1"}, {"b", "2"}}...)
		}
	})
}

func BenchmarkSet(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = client.Set(context.TODO(), "a", "1", 0)
		}
	})
}

func BenchmarkGet(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)
	_ = client.Set(context.TODO(), "a", "1", 0)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = client.Get(context.TODO(), "a")
		}
	})
}
