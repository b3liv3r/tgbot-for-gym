package service

import (
	"context"
	"github.com/b3liv3r/tgbot-for-gym/modules/gym/models"
)

type Gymer interface {
	List(ctx context.Context) ([]models.Gym, error)
	GetSchedules(ctx context.Context, gymId int) ([]models.Schedules, error)
}
