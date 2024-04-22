package sqlettus_test

import (
	"context"
	"testing"
	"time"

	"github.com/jtarchie/sqlettus/executers"
	"github.com/jtarchie/sqlettus/sdk"
)

func BenchmarkAddInt(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)
	_ = client.Set(context.TODO(), "a", "1", 0)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = client.AddInt(context.TODO(), "a", 1)
		}
	})
}

func BenchmarkDBSize(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = client.DBSize(context.TODO())
		}
	})
}

func BenchmarkDelete(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = client.Delete(context.TODO(), "a")
		}
	})
}

func BenchmarkExists(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)
	_ = client.Set(context.TODO(), "a", "1", 0)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = client.Exists(context.TODO(), "a")
		}
	})
}

func BenchmarkExpire(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)
	_ = client.Set(context.TODO(), "a", "1", 0)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = client.Expire(context.TODO(), "a", time.Second)
		}
	})
}

func BenchmarkFlushDB(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = client.FlushDB(context.TODO())
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

func BenchmarkKeys(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)
	_ = client.Set(context.TODO(), "aaa", "1", 0)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = client.Keys(context.TODO(), "a??")
		}
	})
}

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

func BenchmarkListAt(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)
	_, _ = client.ListPush(context.TODO(), "a", "a")

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = client.ListAt(context.TODO(), "a", 0)
		}
	})
}

func BenchmarkListLength(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)
	_, _ = client.ListPush(context.TODO(), "a", "a")

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = client.ListLength(context.TODO(), "a")
		}
	})
}

func BenchmarkListPush(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = client.ListPush(context.TODO(), "a", "a")
		}
	})
}

func BenchmarkListRange(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)
	_, _ = client.ListPush(context.TODO(), "a", "a")

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = client.ListRange(context.TODO(), "a", 0, 1)
		}
	})
}

func BenchmarkListRemove(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)
	_, _ = client.ListPush(context.TODO(), "a", "a")

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = client.ListRemove(context.TODO(), "a", 0, "a")
		}
	})
}

func BenchmarkListShift(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)
	_, _ = client.ListPush(context.TODO(), "a", "a")

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = client.ListShift(context.TODO(), "a", 0)
		}
	})
}

func BenchmarkListUnshift(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)
	_, _ = client.ListPush(context.TODO(), "a", "a")

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = client.ListUnshift(context.TODO(), "a", "0")
		}
	})
}

func BenchmarkRandomKey(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)
	_ = client.Set(context.TODO(), "aaa", "1", 0)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = client.RandomKey(context.TODO())
		}
	})
}

func BenchmarkRename(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)
	_ = client.Set(context.TODO(), "a", "1", 0)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = client.Rename(context.TODO(), "a", "b")
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

func BenchmarkSort(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)
	_, _ = client.ListPush(context.TODO(), "a", "1", "4", "3", "2")

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = client.Sort(context.TODO(), "a")
		}
	})
}

func BenchmarkSubstr(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)
	_ = client.Set(context.TODO(), "a", "some string", 0)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ , _= client.Substr(context.TODO(), "a", 0, 4)
		}
	})
}

func BenchmarkTTL(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)
	_ = client.Set(context.TODO(), "a", "some string", 0)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ , _= client.TTL(context.TODO(), "a")
		}
	})
}

func BenchmarkType(b *testing.B) {
	db, err := executers.FromDB(":memory:")
	if err != nil {
		b.Fatalf("could not start db: %s", err)
	}

	client := sdk.NewClient(db)
	_ = client.Set(context.TODO(), "a", "some string", 0)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ , _= client.Type(context.TODO(), "a")
		}
	})
}