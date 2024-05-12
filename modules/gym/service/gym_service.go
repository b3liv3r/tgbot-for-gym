package service

import (
	"context"
	client "github.com/b3liv3r/tgbot-for-gym/modules/gym/grpc"
	"github.com/b3liv3r/tgbot-for-gym/modules/gym/models"
	"go.uber.org/zap"
)

type Gym struct {
	logger     *zap.Logger
	gymService client.RPCGymer
}

func NewGym(logger *zap.Logger, gymService client.RPCGymer) Gymer {
	return &Gym{logger: logger, gymService: gymService}
}

func (g *Gym) List(ctx context.Context) ([]models.Gym, error) {
	result, err := g.gymService.List(ctx)
	if err != nil {
		g.logger.Error("gym.list", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (g *Gym) GetSchedules(ctx context.Context, gymId int) ([]models.Schedules, error) {
	result, err := g.gymService.GetSchedules(ctx, gymId)
	if err != nil {
		g.logger.Error("gym.getSchedules", zap.Error(err))
		return nil, err
	}

	return result, nil
}
