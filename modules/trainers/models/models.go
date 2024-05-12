package models

import "time"

type Trainer struct {
	ID         int
	GymId      int
	Name       string
	Speciality string
}

type Booking struct {
	ID        int
	UserID    int
	TrainerID int
	Activity  string
	StartTime time.Time
	EndTime   time.Time
}
