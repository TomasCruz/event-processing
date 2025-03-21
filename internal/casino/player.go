package casino

import "time"

type Player struct {
	Email          string    `json:"email,omitempty"`
	LastSignedInAt time.Time `json:"last_signed_in_at,omitempty"`
}

func (p Player) IsZero() bool {
	return p.Email == "" || p.LastSignedInAt.IsZero()
}
