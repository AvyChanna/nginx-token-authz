package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AvyChanna/nginx-token-authz/internal/app"
	"github.com/AvyChanna/nginx-token-authz/internal/rbac/reader"
	"github.com/AvyChanna/nginx-token-authz/internal/server"
)

var (
	configPath   string
	debugEnabled bool
	port         int
)

const PollDuration = 5 * time.Second

func init() {
	flag.StringVar(&configPath, "c", "rbac.yaml", "path to yaml config for rbac")
	flag.BoolVar(&debugEnabled, "d", false, "Enable debug mode")
	flag.IntVar(&port, "p", 8080, "server port")
}

func main() {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app.Init(ctx, debugEnabled)

	fmt.Println(configPath)
	rbacReader := reader.New(configPath)
	autherPtr := rbacReader.WatchConfig(PollDuration)
	if autherPtr.Load() == nil {
		app.Get().Log().Error("error loading config on first run. Exiting...")
		os.Exit(1)
	}

	mux := server.NewMux(autherPtr)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go server.ListenAndServe()

	<-ctx.Done()
	err := server.Shutdown(ctx)
	if err != nil {
		app.Get().Log().Error("error shutting down http server:", err)
	}
}
