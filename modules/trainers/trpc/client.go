package client

import (
	"context"
	trainerv1 "github.com/b3liv3r/protos-for-gym/gen/go/trainer"
	"github.com/b3liv3r/tgbot-for-gym/modules/trainers/models"
	"google.golang.org/grpc"
	"log"
	"time"
)

type RPCTrainerer interface {
	ListForGym(ctx context.Context, gymID int) ([]models.Trainer, error)
	AvailableBookingList(ctx context.Context, trainerID int) ([]models.Booking, error)
	CurrentBookingList(ctx context.Context, userID int) ([]models.Booking, error)
	Booking(ctx context.Context, userId, bookingId int) (string, error)
	UnBooking(ctx context.Context, currentBookingID int) (string, error)
}

type TrainerClient struct {
	rpc  trainerv1.TrainerClient
	addr string
}

func NewTrainerClient(addr string) RPCTrainerer {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := trainerv1.NewTrainerClient(conn)

	return &TrainerClient{
		rpc:  c,
		addr: addr,
	}
}

func (tc *TrainerClient) ListForGym(ctx context.Context, gymID int) ([]models.Trainer, error) {
	resp, err := tc.rpc.ListForGym(ctx, &trainerv1.ListForGymRequest{GymId: int64(gymID)})
	if err != nil {
		return nil, err
	}
	trainers := make([]models.Trainer, len(resp.Trainers))
	for i, t := range resp.Trainers {
		trainers[i] = models.Trainer{
			ID:         int(t.TrainerId),
			GymId:      int(t.GymId),
			Name:       t.Name,
			Speciality: t.Speciality,
		}
	}
	return trainers, nil
}

func (tc *TrainerClient) AvailableBookingList(ctx context.Context, trainerID int) ([]models.Booking, error) {
	resp, err := tc.rpc.AvailableBookingList(ctx, &trainerv1.AvailableBookingListRequest{TrainerId: int64(trainerID)})
	if err != nil {
		return nil, err
	}
	bookings := make([]models.Booking, len(resp.Bookings))
	for i, b := range resp.Bookings {
		bookings[i] = models.Booking{
			ID:        int(b.BookingId),
			Activity:  b.Activity,
			StartTime: b.StartTime.AsTime(),
			EndTime:   b.EndTime.AsTime(),
		}
	}
	return bookings, nil
}

func (tc *TrainerClient) CurrentBookingList(ctx context.Context, userID int) ([]models.Booking, error) {
	resp, err := tc.rpc.CurrentBookingList(ctx, &trainerv1.CurrentBookingListRequest{UserId: int64(userID)})
	if err != nil {
		return nil, err
	}
	bookings := make([]models.Booking, len(resp.Bookings))
	for i, b := range resp.Bookings {
		bookings[i] = models.Booking{
			ID:        int(b.BookingId),
			Activity:  b.Activity,
			StartTime: b.StartTime.AsTime(),
			EndTime:   b.EndTime.AsTime(),
		}
	}
	return bookings, nil
}

func (tc *TrainerClient) Booking(ctx context.Context, userId, bookingId int) (string, error) {
	resp, err := tc.rpc.Booking(ctx, &trainerv1.BookingRequest{
		AvailableBookingId: int64(bookingId),
		UserId:             int64(userId),
	})
	if err != nil {
		return "", err
	}
	return resp.Message, nil
}

func (tc *TrainerClient) UnBooking(ctx context.Context, currentBookingID int) (string, error) {
	resp, err := tc.rpc.UnBooking(ctx, &trainerv1.UnBookingRequest{CurrentBookingId: int64(currentBookingID)})
	if err != nil {
		return "", err
	}
	return resp.Message, nil
}
