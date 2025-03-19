package service

import (
	"time"

	"github.com/TomasCruz/event-processing/internal/config"
	"github.com/TomasCruz/event-processing/internal/ports"
)

type StatsSvc struct {
	Config config.Config
	DB     ports.DB
}

func (svc StatsSvc) Materialized() (ports.MaterializedData, error) {
	total, err := svc.DB.GetEventsTotal()
	if err != nil {
		return ports.MaterializedData{}, err
	}

	eventsPerMinute, err := svc.DB.GetEventsPerMinute()
	if err != nil {
		return ports.MaterializedData{}, err
	}

	movingAverage, err := svc.DB.GetEventsPerSecondMovingAverage(time.Now().Add(-1 * time.Minute))
	if err != nil {
		return ports.MaterializedData{}, err
	}

	topPlayerBetsID, topPlayerBetsCount, err := svc.DB.GetTopPlayerBets()
	if err != nil {
		return ports.MaterializedData{}, err
	}

	topPlayerWinsID, topPlayerWinsCount, err := svc.DB.GetTopPlayerWins()
	if err != nil {
		return ports.MaterializedData{}, err
	}

	topPlayerDepositsID, topPlayerDepositsCount, err := svc.DB.GetTopPlayerDeposits()
	if err != nil {
		return ports.MaterializedData{}, err
	}

	return ports.MaterializedData{
		EventsTotal:                  total,
		EventsPerMinute:              ports.Float2DecFromFloat64(eventsPerMinute),
		EventsPerSecondMovingAverage: ports.Float2DecFromFloat64(movingAverage),
		TopPlayerBets: ports.TopPlayerStats{
			ID:    topPlayerBetsID,
			Count: topPlayerBetsCount,
		},
		TopPlayerWins: ports.TopPlayerStats{
			ID:    topPlayerWinsID,
			Count: topPlayerWinsCount,
		},
		TopPlayerDeposits: ports.TopPlayerStats{
			ID:    topPlayerDepositsID,
			Count: topPlayerDepositsCount,
		},
	}, nil
}
