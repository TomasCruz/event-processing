package database

import (
	"time"

	"github.com/TomasCruz/event-processing/internal/casino"
)

func (pDB postgresDB) GetPlayerByID(id int) (casino.Player, error) {
	var (
		email          string
		lastSignedInAt time.Time
	)

	query := `SELECT email, last_signed_in_at
		FROM players
		WHERE id=$1`
	err := pDB.db.QueryRow(query, int64(id)).Scan(&email, &lastSignedInAt)
	if err != nil {
		return casino.Player{}, err
	}

	return casino.Player{
		Email:          email,
		LastSignedInAt: lastSignedInAt,
	}, nil
}
