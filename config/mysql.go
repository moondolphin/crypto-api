package config

import (
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

func MySQLDSN() (string, error) {
	cfg := mysql.NewConfig()
	cfg.User = Getenv("MYSQL_USER", "root")
	cfg.Passwd = os.Getenv("MYSQL_PASSWORD")
	cfg.Net = "tcp"
	cfg.Addr = Getenv("MYSQL_HOST", "127.0.0.1") + ":" + Getenv("MYSQL_PORT", "3308")
	cfg.DBName = Getenv("MYSQL_DB", "crypto")
	cfg.ParseTime = true

	if cfg.Passwd == "" {
		return "", fmt.Errorf("MYSQL_PASSWORD required")
	}

	return cfg.FormatDSN(), nil
}
