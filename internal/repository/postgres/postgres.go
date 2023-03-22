package postgres

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// PostgresSettings are settings for postgreSQL
type PostgresSettings struct {
	Addr     string
	Database string
	Username string
	Password string
}

// PostgresClient is client for postgreSQL
type PostgresClient struct {
	db *sqlx.DB
}

// NewPostgresClient connects to postgresSQL and return instance of client
func NewPostgresClient(settings PostgresSettings) (*PostgresClient, error) {
	conn := fmt.Sprintf("postgres://%s:%s@%s/%s", settings.Username, settings.Password, settings.Addr, settings.Database)
	db, err := sqlx.Connect("pgx", conn)
	if err != nil {
		return nil, err
	}
	return &PostgresClient{
		db: db,
	}, nil
}

// Ping pings database
func (c *PostgresClient) Ping() error {
	return c.db.Ping()
}

// Close closes database
func (c *PostgresClient) Close() error {
	return c.db.Close()
}
