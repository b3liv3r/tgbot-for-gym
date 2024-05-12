package service

import (
	"context"
	"github.com/b3liv3r/tgbot-for-gym/modules/trainers/models"
)

type Trainerer interface {
	ListForGym(ctx context.Context, gymID int) ([]models.Trainer, error)
	AvailableBookingList(ctx context.Context, trainerID int) ([]models.Booking, error)
	CurrentBookingList(ctx context.Context, userID int) ([]models.Booking, error)
	Booking(ctx context.Context, userId, bookingId int) (string, error)
	UnBooking(ctx context.Context, currentBookingID int) (string, error)
}
