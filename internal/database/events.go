package database

import (
	"time"

	"github.com/TomasCruz/event-processing/internal/casino"
)

func (pDB postgresDB) CreateEvent(ev casino.Event) error {
	pTx, err := pDB.newTransaction()
	if err != nil {
		return err
	}
	defer pTx.commitOrRollbackOnError(&err)

	sqlStatement := `INSERT INTO events (id, player_id, game_id, typ, amount, currency, has_won, created_at, amount_eur)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	stmt, err := pTx.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(ev.ID, ev.PlayerID, ev.GameID, ev.Type, ev.Amount, ev.Currency, ev.HasWon, ev.CreatedAt, ev.AmountEUR); err != nil {
		return err
	}

	return nil
}

func (pDB postgresDB) GetEventsTotal() (int, error) {
	var result int

	query := `SELECT COUNT(1) FROM events`
	err := pDB.db.QueryRow(query).Scan(&result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (pDB postgresDB) GetEventsPerMinute() (float64, error) {
	var result float64

	// prevent division by zero with GREATEST
	query := `SELECT CAST(COUNT(*) AS FLOAT) / GREATEST(COUNT(DISTINCT DATE_TRUNC('minute', created_at)), 1)
		FROM events`

	err := pDB.db.QueryRow(query).Scan(&result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (pDB postgresDB) GetEventsPerSecondMovingAverage(since time.Time) (float64, error) {
	var result float64

	// prevent division by zero with GREATEST
	query := `SELECT CAST(COUNT(*) AS FLOAT) / GREATEST(COUNT(DISTINCT DATE_TRUNC('second', created_at)), 1)
		FROM events
		WHERE created_at > $1`

	err := pDB.db.QueryRow(query, since).Scan(&result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (pDB postgresDB) GetTopPlayerBets() (int, int, error) {
	var playerID, result int

	query := `SELECT player_id, SUM(1) s
		FROM events
		WHERE typ = 'bet'
		GROUP BY player_id
		ORDER BY s DESC`

	err := pDB.db.QueryRow(query).Scan(&playerID, &result)
	if err != nil {
		return 0, 0, err
	}

	return playerID, result, nil
}

func (pDB postgresDB) GetTopPlayerWins() (int, int, error) {
	var playerID, result int

	query := `SELECT player_id, SUM(1) s
		FROM events
		WHERE typ = 'bet' AND has_won
		GROUP BY player_id
		ORDER BY s DESC`

	err := pDB.db.QueryRow(query).Scan(&playerID, &result)
	if err != nil {
		return 0, 0, err
	}

	return playerID, result, nil
}

func (pDB postgresDB) GetTopPlayerDeposits() (int, int, error) {
	var playerID, result int

	query := `SELECT player_id, SUM(amount_eur) s
		FROM events
		WHERE typ = 'deposit'
		GROUP BY player_id
		ORDER BY s DESC`

	err := pDB.db.QueryRow(query).Scan(&playerID, &result)
	if err != nil {
		return 0, 0, err
	}

	return playerID, result, nil
}
