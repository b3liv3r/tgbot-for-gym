package service

import (
	"context"
	"github.com/b3liv3r/tgbot-for-gym/modules/subscription/models"
)

type Subscriptioner interface {
	Create(ctx context.Context, userId int) (string, error)
	UpdateType(ctx context.Context, userId, Type, month int) (string, error)
	Extend(ctx context.Context, userId, month int) (string, error)
	GetData(ctx context.Context, userId int) (models.Subscription, error)
}
