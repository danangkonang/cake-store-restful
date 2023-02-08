package service

import "database/sql"

type connection struct {
	Mysql *sql.DB
}
