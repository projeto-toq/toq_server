package mysqladapter

import (
	"database/sql"
)

type Database struct {
	DB *sql.DB
}

func NewDB(db *sql.DB) *Database {
	return &Database{
		DB: db,
	}
}
