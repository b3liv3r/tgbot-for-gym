package models

import "time"

type Subscription struct {
	Id               int
	SubscriptionType int
	StartDate        time.Time
	EndDate          time.Time
}
