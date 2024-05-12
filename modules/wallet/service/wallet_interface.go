package service

import (
	"context"
	"github.com/b3liv3r/tgbot-for-gym/modules/wallet/models"
	"time"
)

type Walleter interface {
	Create(ctx context.Context, userID int) (string, error)
	GetBalance(ctx context.Context, userID int) (float64, error)
	Transaction(ctx context.Context, transaction models.Transaction) (string, float64, error)
	History(ctx context.Context, userID int, Type string, startTime, endTime time.Time) ([]models.Transaction, error)
}
