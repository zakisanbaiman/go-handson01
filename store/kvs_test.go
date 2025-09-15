package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/zakisanbaiman/go-handson01/entity"
	"github.com/zakisanbaiman/go-handson01/testutil"
)

func TestKVS_Save(t *testing.T) {
	t.Parallel()

	client := testutil.OpenRedisForTest(t)

	sut := &KVS{
		Cli: client,
	}

	key := "TestKVS_Save"
	userID := entity.UserID(1)
	ctx := context.Background()

	t.Cleanup(func() {
		client.Del(ctx, key)
	})

	if err := sut.Save(ctx, key, userID); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
}

func TestKVS_Load(t *testing.T) {
	t.Parallel()

	client := testutil.OpenRedisForTest(t)
	sut := &KVS{
		Cli: client,
	}

	key := "TestKVS_Load_ok"
	userID := entity.UserID(1234)
	ctx := context.Background()

	client.Set(ctx, key, userID, 30*time.Second)
	t.Cleanup(func() {
		client.Del(ctx, key)
	})

	got, err := sut.Load(ctx, key)
	if err != nil {
		t.Errorf("want no error, but got %v", err)
	}

	if got != userID {
		t.Errorf("want %d, but got %d", userID, got)
	}

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		key := "TestKVS_Load_not_found"
		ctx := context.Background()
		got, err := sut.Load(ctx, key)
		if err == nil || !errors.Is(err, ErrNotFound) {
			t.Errorf("want %v, but got %v(value: %d)", ErrNotFound, err, got)
		}
	})
}
