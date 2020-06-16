package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v4"
	"github.com/yrfg/timeman/pkg/server"
	"github.com/yrfg/timeman/pkg/store"
)

var (
	DefaultConn *pgx.Conn
)

func main() {
	store.Setup(store.DataBaseSetupConfig{
		Host:     "localhost",
		Port:     5532,
		Database: "nakama",
		User:     "nakama",
		Password: "nakama",
	})

	server.Setup(server.ServerSetupConfig{
		Mode: server.ServerSetupConfigModeDebug,
		Port: 54431,
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
