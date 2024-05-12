package models

type User struct {
	Id              int
	Username        string
	Phone           string
	Email           string
	SubscriptionLvl int
	CurrentGymId    int
}
