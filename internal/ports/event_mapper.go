package ports

import (
	"github.com/TomasCruz/event-processing/internal/casino"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func EventToPB(ev casino.Event) *Event {
	var hasWonValid bool
	var hasWonValue bool

	if ev.HasWon != nil {
		hasWonValid = true
		hasWonValue = *ev.HasWon
	}

	return &Event{
		Id:       int32(ev.ID),
		PlayerId: int32(ev.PlayerID),
		GameId:   int32(ev.GameID),
		Typ:      ev.Type,
		Amount:   int32(ev.Amount),
		Currency: ev.Currency,
		HasWon: &BoolPtr{
			Valid: hasWonValid,
			Value: hasWonValue,
		},
		CreatedAt: timestamppb.New(ev.CreatedAt),
		AmountEur: int32(ev.AmountEUR),
		Player:    PlayerToPB(ev.Player),
	}
}

func PBToEvent(ev *Event) casino.Event {
	var hasWon *bool
	if ev.HasWon.Valid {
		hasWon = &ev.HasWon.Value
	}

	return casino.Event{
		ID:        int(ev.Id),
		PlayerID:  int(ev.PlayerId),
		GameID:    int(ev.GameId),
		Type:      ev.Typ,
		Amount:    int(ev.Amount),
		Currency:  ev.Currency,
		HasWon:    hasWon,
		CreatedAt: ev.CreatedAt.AsTime(),
		AmountEUR: int(ev.AmountEur),
		Player:    PBToPlayer(ev.Player),
	}
}

func PlayerToPB(player casino.Player) *Player {
	return &Player{
		Email:          player.Email,
		LastSignedInAt: timestamppb.New(player.LastSignedInAt),
	}
}

func PBToPlayer(player *Player) casino.Player {
	return casino.Player{
		Email:          player.Email,
		LastSignedInAt: player.LastSignedInAt.AsTime(),
	}
}
