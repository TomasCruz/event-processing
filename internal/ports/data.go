package ports

// I usually put all the biz logic definitions in entities package;
// here I wanted to minimize changes to existing casino package, so I added this ports package
// Idea is to have all the biz logic definitions in one place which has no code deps on the rest of the solution,
// as it is the innermost layer of hexagonal architecture
type TopPlayerStats struct {
	ID    int `json:"id"`
	Count int `json:"count"`
}

type Float2Dec float64

// this trick rounds input to 2 decimal places
func Float2DecFromFloat64(x float64) Float2Dec {
	return Float2Dec(float64(int(x*100)) / 100)
}

type MaterializedData struct {
	EventsTotal                  int            `json:"events_total"`
	EventsPerMinute              Float2Dec      `json:"events_per_minute"`
	EventsPerSecondMovingAverage Float2Dec      `json:"events_per_second_moving_average"`
	TopPlayerBets                TopPlayerStats `json:"top_player_bets"`
	TopPlayerWins                TopPlayerStats `json:"top_player_wins"`
	TopPlayerDeposits            TopPlayerStats `json:"top_player_deposits"`
}
