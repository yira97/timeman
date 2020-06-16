package store

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
)

var DefaultConn *pgx.Conn

type DataBaseSetupConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

func Setup(setupCfg DataBaseSetupConfig) {
	var err error
	connCfg, err := pgx.ParseConfig(
		fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", setupCfg.User, setupCfg.Password, setupCfg.Host, setupCfg.Port, setupCfg.Database),
	)
	if err != nil {
		log.Fatalf("Unable to connect to database(cfg): %v\n", err)
	}
	DefaultConn, err = pgx.ConnectConfig(context.Background(), connCfg)
	if err != nil {
		log.Fatalf("Unable to connect to database(conn): %v\n", err)
	}
	err = TableConstruct(context.Background(),DefaultConn)
	if err != nil {
		log.Fatalf("Unable to init database: %v\n", err)
	}
}
