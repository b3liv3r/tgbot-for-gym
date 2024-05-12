package service

import (
	"context"
	"errors"
	"github.com/b3liv3r/tgbot-for-gym/modules/trainers/models"
	client "github.com/b3liv3r/tgbot-for-gym/modules/trainers/trpc"
	"go.uber.org/zap"
)

type Trainer struct {
	logger         *zap.Logger
	trainerService client.RPCTrainerer
}

func NewTrainer(logger *zap.Logger, trainerService client.RPCTrainerer) Trainerer {
	return &Trainer{
		trainerService: trainerService,
		logger:         logger,
	}
}

func (ts *Trainer) ListForGym(ctx context.Context, gymID int) ([]models.Trainer, error) {
	result, err := ts.trainerService.ListForGym(ctx, gymID)
	if err != nil {
		ts.logger.Error("trainerService.ListForGym", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (ts *Trainer) AvailableBookingList(ctx context.Context, trainerID int) ([]models.Booking, error) {
	result, err := ts.trainerService.AvailableBookingList(ctx, trainerID)
	if err != nil {
		ts.logger.Error("trainerService.AvailableBookingList", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (ts *Trainer) CurrentBookingList(ctx context.Context, userID int) ([]models.Booking, error) {
	result, err := ts.trainerService.CurrentBookingList(ctx, userID)
	if err != nil {
		ts.logger.Error("trainerService.CurrentBookingList", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (ts *Trainer) Booking(ctx context.Context, userId, bookingId int) (string, error) {
	result, err := ts.trainerService.Booking(ctx, userId, bookingId)
	if err != nil {
		if err.Error() == "rpc error: code = Unknown desc = невозможно забронировать прошедший слот" {
			return "", errors.New("unable to book past slot")
		} else if err.Error() == "rpc error: code = Unknown desc = sql: no rows in result set" {
			return "", errors.New("Unknown booking, use /t_available_booking_list <trainer_id> to see available booking")
		}
		ts.logger.Error("trainerService.Booking", zap.Error(err))
		return "", err
	}

	return result, nil
}

func (ts *Trainer) UnBooking(ctx context.Context, currentBookingID int) (string, error) {
	result, err := ts.trainerService.UnBooking(ctx, currentBookingID)
	if err != nil {
		if err.Error() == "rpc error: code = Unknown desc = невозможно отменить бронирование, так как до начала слота остается 24 часа или менее" {
			return "", errors.New("it is not possible to cancel the reservation as there are 24 hours or less left before the slot starts")
		} else if err.Error() == "rpc error: code = Unknown desc = sql: no rows in result set" {
			return "", errors.New("You don't have such a booking, use /t_current_booking_list to see your current booking")
		}
		ts.logger.Error("trainerService.UnBooking", zap.Error(err))
		return "", err
	}

	return result, nil
}
