package model

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// MysqlConnect creates a connection to to the nuclious db
func InitDB(dsn string) *sqlx.DB {
	var db *sqlx.DB
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	// Preventing connection timeout for idle connections or connections which take too much time
	db.SetMaxIdleConns(0)
	return db
}
