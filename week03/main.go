package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	g, errCtx := errgroup.WithContext(ctx)
	srv := &http.Server{Addr: ":9090"}
	// serve http
	g.Go(func() error {
		return srv.ListenAndServe()
	})
	// 关闭http server
	g.Go(func() error {
		<-errCtx.Done()
		return srv.Shutdown(errCtx)
	})
	// 处理signal
	g.Go(func() error {
		select {
		case <-errCtx.Done():
			return errCtx.Err()
		case <-c:
			cancel()
		}
		return nil
	})
	err := g.Wait()
	fmt.Println(err)
	fmt.Println(ctx.Err())
}
