package domain

import "time"

type BookRegistered struct {
	ISBN       ISBN
	Title      string
	OccurredAt time.Time
}

func (e BookRegistered) EventName() string { return "BookRegistered" }
