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

	// Configure connection pool
	database.SetMaxOpenConns(25)                 // Maximum number of open connections
	database.SetMaxIdleConns(10)                 // Maximum number of idle connections
	database.SetConnMaxLifetime(5 * time.Minute) // Connection maximum lifetime
	database.SetConnMaxIdleTime(2 * time.Minute) // Connection maximum idle time

	//testa se o conexão está viva
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Increased timeout
	defer cancel()

	if err = database.PingContext(ctx); err != nil {
		slog.Error("error on MySql connection", "error", err)
		panic(err)
	}
	slog.Info("Database answered the ping. MySql connection successfuly!")

	// Run lightweight migrations (idempotent) for new tables
	if err := runMigrations(database); err != nil {
		slog.Error("database migration failed", "error", err)
		panic(err)
	}

	c.db = database

}

func (c *config) GetDatabase() *sql.DB {
	return c.db
}

// runMigrations executes idempotent DDL for required tables.
// NOTE: For production consider using a full migration tool, this is a minimal bootstrap.
func runMigrations(db *sql.DB) error {
	// sessions table for refresh token tracking
	ddl := `CREATE TABLE IF NOT EXISTS sessions (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		user_id BIGINT NOT NULL,
		refresh_hash CHAR(64) NOT NULL UNIQUE,
		token_jti CHAR(36) NULL,
		expires_at DATETIME(6) NOT NULL,
		absolute_expires_at DATETIME(6) NULL,
		created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
		rotated_at DATETIME(6) NULL,
		user_agent VARCHAR(255) NULL,
		ip VARCHAR(45) NULL,
		device_id VARCHAR(100) NULL,
		rotation_counter INT NOT NULL DEFAULT 0,
		last_refresh_at DATETIME(6) NULL,
		revoked TINYINT(1) NOT NULL DEFAULT 0,
		INDEX idx_sessions_user_id (user_id),
		INDEX idx_sessions_expires_at (expires_at),
		INDEX idx_sessions_revoked (revoked),
		INDEX idx_sessions_token_jti (token_jti)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

	// Attempt to add new columns if table existed (idempotent best-effort)
	alterStmts := []string{
		"ALTER TABLE sessions ADD COLUMN token_jti CHAR(36) NULL AFTER refresh_hash",
		"ALTER TABLE sessions ADD COLUMN absolute_expires_at DATETIME(6) NULL AFTER expires_at",
		"ALTER TABLE sessions ADD COLUMN device_id VARCHAR(100) NULL AFTER ip",
		"ALTER TABLE sessions ADD COLUMN rotation_counter INT NOT NULL DEFAULT 0 AFTER device_id",
		"ALTER TABLE sessions ADD COLUMN last_refresh_at DATETIME(6) NULL AFTER rotation_counter",
		"CREATE INDEX idx_sessions_token_jti ON sessions(token_jti)",
	}
	for _, stmt := range alterStmts {
		_, _ = db.Exec(stmt) // ignore errors if already applied
	}

	if _, err := db.Exec(ddl); err != nil {
		return err
	}
	return nil
}
