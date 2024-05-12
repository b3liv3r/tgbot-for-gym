package service

import (
	"context"
	"github.com/b3liv3r/tgbot-for-gym/modules/subscription/models"
	client "github.com/b3liv3r/tgbot-for-gym/modules/subscription/srpc"
	"go.uber.org/zap"
)

type Subscription struct {
	logger              *zap.Logger
	subscriptionService client.RPCSubscriptioner
}

func NewSubscription(logger *zap.Logger, subscriptionService client.RPCSubscriptioner) Subscriptioner {
	return &Subscription{
		subscriptionService: subscriptionService,
		logger:              logger,
	}
}

func (s *Subscription) Create(ctx context.Context, userId int) (string, error) {
	result, err := s.subscriptionService.Create(ctx, userId)
	if err != nil {
		s.logger.Error("subscription.create", zap.Error(err))
		return "", err
	}

	return result, nil
}

func (s *Subscription) UpdateType(ctx context.Context, userId, Type, month int) (string, error) {
	result, err := s.subscriptionService.UpdateType(ctx, userId, Type, month)
	if err != nil {
		s.logger.Error("subscription.update", zap.Error(err))
		return "", err
	}

	return result, nil
}

func (s *Subscription) Extend(ctx context.Context, userId, month int) (string, error) {
	result, err := s.subscriptionService.Extend(ctx, userId, month)
	if err != nil {
		s.logger.Error("subscription.extend", zap.Error(err))
		return "", err
	}

	return result, nil
}

func (s *Subscription) GetData(ctx context.Context, userId int) (models.Subscription, error) {
	result, err := s.subscriptionService.GetData(ctx, userId)
	if err != nil {
		s.logger.Error("subscription.getData", zap.Error(err))
		return models.Subscription{}, err
	}
	return result, nil
}
