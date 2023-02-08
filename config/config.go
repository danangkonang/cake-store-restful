package config

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	// _ "time/tzdata"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	Mysql *sql.DB
}

func Connection() *DB {
	connection := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		"Asia%2FJakarta",
	)
	db, err := sql.Open(os.Getenv("DB_DRIVER"), connection)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	return &DB{Mysql: db}
}
