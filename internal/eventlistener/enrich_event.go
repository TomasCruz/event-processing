package eventlistener

import (
	"log"

	"github.com/TomasCruz/event-processing/internal/ports"
)

func (eSvc enricherSvc) enrichEvent(ev *ports.Event) {
	playerCh := make(chan *ports.Event)
	currencyCh := make(chan *ports.Event)

	go eSvc.enrichEventPlayerData(ev, playerCh)
	go eSvc.enrichEventCurrencyData(ev, currencyCh)

	<-playerCh
	<-currencyCh
}

func (eSvc enricherSvc) enrichEventPlayerData(ev *ports.Event, evCh chan<- *ports.Event) {
	player, err := eSvc.db.GetPlayerByID(int(ev.PlayerId))
	if err != nil {
		// process event further without enrichment
		log.Printf("couldn't get data for player id: %d, err: %v\n", ev.PlayerId, err)
		evCh <- ev
		return
	}

	ev.Player = ports.PlayerToPB(player)
	evCh <- ev
}

func (eSvc enricherSvc) enrichEventCurrencyData(ev *ports.Event, evCh chan<- *ports.Event) {
	// https://api.exchangerate.host/ redirects to https://apilayer.com/. Free subscription there is for up to a 100 requests per months,
	// so it will stop working before the end of work on this challenge.
	// as a fallback, https://api.freecurrencyapi.com/ will be used, as
	// switching between data providers is a normal business practice anyway.
	// As freecurrencyapi doesn't support BTC rates,
	// hardcoded 2025-03-15 rate of 1 BTC -> EUR 77070.53064 will be used.
	// Apart from the described hack, I will also cut corners with regards to rate caching.
	// In the professinal setting, I would use Redis with 1 min expiration entries, while function loading the rate
	// would first try fetching it from Redis, and if it's not there would get the rate from the API and store new entry to Redis.

	// enrich with converted currency data
	if ev.Currency == "EUR" {
		ev.AmountEur = ev.Amount
		evCh <- ev
		return
	}

	rate, err := eSvc.doAPILayerReq(ev.Currency)
	if err != nil {
		if ev.Currency == "BTC" {
			rate = 77070.53064
		} else {
			rate, err = eSvc.doFreecurrencyAPIReq(ev.Currency)
			if err != nil {
				// send the event without enriching for further processing
				log.Printf("Event %d no enrichment: %v\n", ev.Id, err)
				evCh <- ev
				return
			}
		}
	}

	if rate != 0.0 {
		// Examples: 300 = 3.00 EUR, 1 = 0.00000001 BTC.
		amountEUR := 0.0
		if ev.Currency == "BTC" {
			amountEUR = (float64(ev.Amount) * rate) / float64(1e6)
		} else {
			amountEUR = float64(ev.Amount) * rate
		}

		// I'm not sure what rounding method is used here, so using int conversion
		ev.AmountEur = int32(amountEUR)
	}

	evCh <- ev
}
