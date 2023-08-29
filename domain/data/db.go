package data

import (
	"database/sql"

	"github.com/Karik-ribasu/golang-todo-list-api/infra/config"
	"github.com/go-sql-driver/mysql"
)

type conn struct {
	db *sql.DB
	tx sql.Tx
}

func InitDB(cfg config.Config) (dbManager DbManager, err error) {
	mysqlConfig := mysql.NewConfig()
	mysqlConfig.User = cfg.Db.User
	mysqlConfig.Passwd = cfg.Db.Passwd
	mysqlConfig.Addr = cfg.Db.Addr + "::" + cfg.Db.Port
	mysqlConfig.DBName = cfg.Db.Name

	conn := conn{}
	db, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		return &conn, err
	}

	conn.db = db
	return &conn, err
}
