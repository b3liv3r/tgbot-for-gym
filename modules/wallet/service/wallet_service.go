package service

import (
	"context"
	"errors"
	"github.com/b3liv3r/tgbot-for-gym/modules/wallet/models"
	client "github.com/b3liv3r/tgbot-for-gym/modules/wallet/wrpc"
	"go.uber.org/zap"
	"time"
)

type Wallet struct {
	logger        *zap.Logger
	walletService client.RPCWalleter
}

func NewWallet(logger *zap.Logger, walletService client.RPCWalleter) Walleter {
	return &Wallet{logger: logger, walletService: walletService}
}

func (w *Wallet) Create(ctx context.Context, userID int) (string, error) {
	result, err := w.walletService.Create(ctx, userID)
	if err != nil {
		w.logger.Error("walletService.Create", zap.Error(err))
		return "", err
	}

	return result, nil
}

func (w *Wallet) GetBalance(ctx context.Context, userID int) (float64, error) {
	result, err := w.walletService.GetBalance(ctx, userID)
	if err != nil {
		w.logger.Error("walletService.GetBalance", zap.Error(err))
		return 0, err
	}

	return result, nil
}

func (w *Wallet) Transaction(ctx context.Context, transaction models.Transaction) (string, float64, error) {
	result, currentBalance, err := w.walletService.Transaction(ctx, transaction)
	if err != nil {
		if err.Error() == "rpc error: code = Unknown desc = insufficient funds" {
			return "", 0, errors.New("insufficient funds")
		}
		w.logger.Error("walletService.Transaction", zap.Error(err))
		return "", 0, err
	}

	return result, currentBalance, nil
}

func (w *Wallet) History(ctx context.Context, userID int, Type string, startTime, endTime time.Time) ([]models.Transaction, error) {
	result, err := w.walletService.History(ctx, userID, Type, startTime, endTime)
	if err != nil {
		w.logger.Error("walletService.History", zap.Error(err))
		return nil, err
	}

	return result, nil
}
