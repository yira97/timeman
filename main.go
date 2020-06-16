package main

import (
	"flag"
	"github.com/yrfg/timeman/pkg/setting"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v4"
	"github.com/yrfg/timeman/pkg/server"
	"github.com/yrfg/timeman/pkg/store"
)

var (
	DefaultConn *pgx.Conn
	task = flag.String("task", "start-local", "")
)

func main() {
	flag.Parse()
	switch *task {
	case "test":
		os.Exit(0)
	}

	setting.Setup("configs.json")
	cfg := setting.GetGlobalConfig()

	store.Setup(store.DataBaseSetupConfig{
		Host:    cfg.DB.Host,
		Port:     cfg.DB.Port,
		Database: cfg.DB.Database,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
	})

	server.Setup(server.ServerSetupConfig{
		Mode: server.ServerMode(cfg.Server.Mode),
		Port: cfg.Server.Port,
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
