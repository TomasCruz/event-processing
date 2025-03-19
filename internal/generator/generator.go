package generator

import (
	"context"
	"math/rand"
	"time"

	"github.com/TomasCruz/event-processing/internal/casino"
)

func Generate(ctx context.Context) <-chan casino.Event {
	eventCh := make(chan casino.Event)
	var id int

	go func() {
		defer close(eventCh)

		for {
			id++

			select {
			case <-ctx.Done():
				return
			default:
				eventCh <- generate(id)
			}

			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		}
	}()

	return eventCh
}

func generate(id int) casino.Event {
	amount, currency := randomAmountCurrency()

	// has_won field only make sense for bet type
	typ := randomType()
	var hasWon *bool
	if typ == "bet" {
		won := randomHasWon()
		hasWon = &won
	}

	return casino.Event{
		ID:        id,
		PlayerID:  10 + rand.Intn(10),
		GameID:    100 + rand.Intn(10),
		Type:      typ,
		Amount:    amount,
		Currency:  currency,
		HasWon:    hasWon,
		CreatedAt: time.Now(),
	}
}

func randomType() string {
	return casino.EventTypes[rand.Intn(len(casino.EventTypes))]
}

func randomAmountCurrency() (amount int, currency string) {
	currency = casino.Currencies[rand.Intn(len(casino.Currencies))]

	switch currency {
	case "BTC":
		amount = rand.Intn(1e5)
	default:
		amount = rand.Intn(2000)
	}

	return
}

func randomHasWon() bool {
	return rand.Intn(100) < 5
}
