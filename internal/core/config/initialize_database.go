package config

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func (c *config) InitializeDatabase() {

	//abre o banco no URI
	database, err := sql.Open("mysql", c.env.DB.URI)
	if err != nil {
		slog.Error("error trying to open mysql", "error", err)
		panic(err)
	}

	//testa se o conexão está viva
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err = database.PingContext(ctx); err != nil {
		slog.Error("error on MySql connection", "error", err)
		panic(err)
	}
	slog.Info("Database answered the ping. MySql connection successfuly!")

	c.db = database

}

func (c *config) GetDatabase() *sql.DB {
	return c.db
}
