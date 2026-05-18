package data

import (
	"database/sql"

	"github.com/Karik-ribasu/golang-todo-list-api/infra/config"
	"github.com/go-sql-driver/mysql"
)

// openSQL is sql.Open by default; tests may replace it (e.g. sqlmock).
var openSQL = sql.Open

type conn struct {
	db *sql.DB
}

func InitDB(cfg config.Config) (dbManager DbManager, err error) {
	mysqlConfig := mysql.NewConfig()
	mysqlConfig.Net = "tcp"
	mysqlConfig.User = cfg.Db.User
	mysqlConfig.Passwd = cfg.Db.Passwd
	mysqlConfig.Addr = cfg.Db.Addr + ":" + cfg.Db.Port
	mysqlConfig.DBName = cfg.Db.Name

	db, err := openSQL("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &conn{db: db}, nil
}

// NewManagerFromDB exposes a DbManager backed by an existing pool (for tests).
func NewManagerFromDB(db *sql.DB) DbManager {
	return &conn{db: db}
}
