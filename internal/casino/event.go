package casino

import "time"

var EventTypes = []string{
	"game_start",
	"bet",
	"deposit",
	"game_stop",
}

// Event field HasWon has to be *bool, as there is a need to distinguish between undefined and false values,
// because of clearness of stored data (field is only relevant for 'bet' type),
// and in JSON representation omitempty for bool would omit false values
// description is not used, so it is kicked out

type Event struct {
	ID       int `json:"id"`
	PlayerID int `json:"player_id"`

	// Except for `deposit`.
	GameID int `json:"game_id,omitempty"`

	Type string `json:"type"`

	// Smallest possible unit for the given currency.
	// Examples: 300 = 3.00 EUR, 1 = 0.00000001 BTC.
	// Only for types `bet` and `deposit`.
	Amount int `json:"amount,omitempty"`

	// Only for types `bet` and `deposit`.
	Currency string `json:"currency,omitempty"`

	// Only for type `bet`.
	HasWon *bool `json:"has_won,omitempty"`

	CreatedAt time.Time `json:"created_at"`

	AmountEUR int    `json:"amount_eur,omitempty"`
	Player    Player `json:"player,omitempty"`
}
