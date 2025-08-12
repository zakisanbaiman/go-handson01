package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/zakisanbaiman/go-handson01/config"
)

// func run(ctx context.Context) error {
// 	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
// 	defer stop()

// 	cfg, err := config.New()
// 	if err != nil {
// 		return err
// 	}

// 	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
// 	if err != nil {
// 		return err
// 	}
// 	url := fmt.Sprintf("http://%s", l.Addr().String())
// 	log.Printf("start with %s", url)

// 	s := &http.Server{
// 		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			time.Sleep(5 * time.Second)
// 			fmt.Fprintf(w, "Hello, %s! 🚀 LIVE HOT RELOAD TEST!5 🚀", r.URL.Path[1:])
// 		}),
// 	}
// 	eg, ctx := errgroup.WithContext(ctx)

// 	// 別ゴルーチンでHTTPサーバを起動する
// 	eg.Go(func() error {
// 		// http.ErrServerClosedは
// 		// http.Server.Shutdown()が正常に終了したことを示すので以上ではない.
// 		if err := s.Serve(l); err != nil &&
// 			err != http.ErrServerClosed {
// 			log.Printf("failed to close: %+v", err)
// 			return err
// 		}
// 		return nil
// 	})

// 	// チャネルからの終了通知を待つ
// 	<-ctx.Done()

// 	// サーバーを停止する
// 	if err := s.Shutdown(context.Background()); err != nil {
// 		log.Printf("failed to shutdown: %+v", err)
// 	}

// 	return eg.Wait()
// }

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Printf("failed to read config: %v", err)
		os.Exit(1)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Printf("failed to listen port %d: %v", cfg.Port, err)
		os.Exit(1)
	}

	mux := NewMux()
	server := NewServer(l, mux)

	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		return server.Run(ctx)
	})

	if err := eg.Wait(); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}
