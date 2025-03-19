package ports

import (
	"io"
	"time"

	"github.com/TomasCruz/event-processing/internal/casino"
)

type DB interface {
	io.Closer
	GetPlayerByID(id int) (casino.Player, error)
	CreateEvent(casino.Event) error
	GetEventsTotal() (int, error)
	GetEventsPerMinute() (float64, error)
	GetEventsPerSecondMovingAverage(since time.Time) (float64, error)
	GetTopPlayerBets() (int, int, error)
	GetTopPlayerWins() (int, int, error)
	GetTopPlayerDeposits() (int, int, error)
}
