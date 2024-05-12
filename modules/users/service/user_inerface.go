package service

import (
	"context"
	"github.com/b3liv3r/tgbot-for-gym/modules/users/models"
)

type Userer interface {
	Create(ctx context.Context, user models.User) (string, error)
	Profile(ctx context.Context, userID int) (models.User, error)
	Update(ctx context.Context, user models.User) (string, error)
	ChangeSubscription(ctx context.Context, userID, subLvl int) (string, error)
	ChangeCurrentGym(ctx context.Context, userID, gymID int) (string, error)
}
