package store

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var DefaultConn *pgx.Conn

type DataBaseSetupConfig struct {
	Host     string
	Port     int64
	Database string
	User     string
	Password string
}

func Setup(setupCfg DataBaseSetupConfig) {
	var err error
	connCfg, err := pgx.ParseConfig("")
	// connCfg := pgx.ConnConfig{
	// 	Config: pgconn.Config{
	// 		Host:     setupCfg.Host,
	// 		Port:     uint16(setupCfg.Port),
	// 		Database: setupCfg.Database,
	// 		User:     setupCfg.User,
	// 		Password: setupCfg.Password,
	// 	},
	// }
	DefaultConn, err = pgx.ConnectConfig(context.Background(), connCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
}
