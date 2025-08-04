package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen port %v", err)
	}
	defer func() {
		l.Close()
	}()
	
	// cancelableなcontextを作成
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx, l)
	})

	// 動的ポートを取得
	url := fmt.Sprintf("http://%s/", l.Addr().String())
	
	in := "message"
	rsp, err := http.Get(url + in)
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