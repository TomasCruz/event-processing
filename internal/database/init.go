package database

import (
	"database/sql"
	"fmt"

	"github.com/TomasCruz/event-processing/internal/config"
	"github.com/TomasCruz/event-processing/internal/ports"
	_ "github.com/lib/pq"
)

type postgresDB struct {
	db     *sql.DB
	config config.Config
}

type postgresTx struct {
	*sql.Tx
}

func New(c config.Config) (ports.DB, error) {
	db, err := sql.Open("postgres", c.DBURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging the database: %w", err)
	}

	return postgresDB{
		db:     db,
		config: c,
	}, nil
}
