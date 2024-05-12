package models

import "time"

type Gym struct {
	Id      int
	Address string
	SubLvl  int
}

type Schedules struct {
	Id        int
	GymId     int
	DayOfWeek string
	StartTime time.Time
	EndTime   time.Time
}
