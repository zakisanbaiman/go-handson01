package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx)
	})

	in := "message"
	rsp, err := http.Get("http://localhost:18080/" + in)
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}

	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body);
	if err != nil {
		t.Errorf("failed to read body: %+v", err)
	}

	// HTTPサーバの戻り値を検証する
	expected := fmt.Sprintf("Hello, %s!", in)
	if string(got) != expected {
		t.Errorf("want %q, but got %q", expected, string(got))
	}

	// run関数に終了通知を送信する。
	cancel()

	// 終了通知を送信したら、run関数が終了するのを待つ。
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}