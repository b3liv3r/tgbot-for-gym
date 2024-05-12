package client

import (
	"context"
	gymv1 "github.com/b3liv3r/protos-for-gym/gen/go/gym"
	"github.com/b3liv3r/tgbot-for-gym/modules/gym/models"
	"google.golang.org/grpc"
	"log"
	"time"
)

type RPCGymer interface {
	List(ctx context.Context) ([]models.Gym, error)
	GetSchedules(ctx context.Context, gymId int) ([]models.Schedules, error)
}

type GymClient struct {
	rpc  gymv1.GymClient
	addr string
}

func NewGymClient(addr string) RPCGymer {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := gymv1.NewGymClient(conn)

	return &GymClient{
		rpc:  c,
		addr: addr,
	}
}

func (g *GymClient) List(ctx context.Context) ([]models.Gym, error) {
	resp, err := g.rpc.List(ctx, &gymv1.ListRequest{})
	if err != nil {
		return nil, err
	}

	var result []models.Gym
	for _, val := range resp.GymList {
		result = append(result, models.Gym{
			Id:      int(val.GymId),
			Address: val.Address,
			SubLvl:  int(val.SubLvl),
		})
	}
	return result, nil
}

func (g *GymClient) GetSchedules(ctx context.Context, gymId int) ([]models.Schedules, error) {
	resp, err := g.rpc.GetSchedules(ctx, &gymv1.GetSchedulesRequest{GymId: int64(gymId)})
	if err != nil {
		return nil, err
	}

	var result []models.Schedules
	for _, val := range resp.ScheduleList {
		startTime := time.Unix(val.StartTime.GetSeconds(), int64(val.StartTime.GetNanos()))
		endTime := time.Unix(val.EndTime.GetSeconds(), int64(val.EndTime.GetNanos()))
		result = append(result, models.Schedules{
			Id:        int(val.ScheduleId),
			GymId:     int(val.GymId),
			DayOfWeek: val.DayOfWeek,
			StartTime: startTime,
			EndTime:   endTime,
		})
	}
	return result, nil
}
