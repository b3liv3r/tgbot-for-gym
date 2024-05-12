package service

import (
	"context"
	"errors"
	"github.com/b3liv3r/tgbot-for-gym/modules/users/models"
	client "github.com/b3liv3r/tgbot-for-gym/modules/users/urpc"
	"go.uber.org/zap"
)

type User struct {
	logger      *zap.Logger
	userService client.RPCUserer
}

func NewUser(logger *zap.Logger, userService client.RPCUserer) Userer {
	return &User{logger: logger, userService: userService}
}

func (u *User) Create(ctx context.Context, user models.User) (string, error) {
	result, err := u.userService.Create(ctx, user)
	if err != nil {
		u.logger.Error("user.create", zap.Error(err))
		return "", err
	}

	return result, nil
}

func (u *User) Profile(ctx context.Context, userID int) (models.User, error) {
	result, err := u.userService.Profile(ctx, userID)
	if err != nil {
		if err.Error() == "rpc error: code = Unknown desc = sql: no rows in result set" {
			return models.User{}, errors.New("user not found")
		}
		u.logger.Error("user.profile", zap.Error(err))
		return models.User{}, err
	}

	return result, nil
}

func (u *User) Update(ctx context.Context, user models.User) (string, error) {
	result, err := u.userService.Update(ctx, user)
	if err != nil {
		u.logger.Error("user.update", zap.Error(err))
		return "", err
	}

	return result, nil
}

func (u *User) ChangeSubscription(ctx context.Context, userID, subLvl int) (string, error) {
	result, err := u.userService.ChangeSubscription(ctx, userID, subLvl)
	if err != nil {
		u.logger.Error("user.changeSubscription", zap.Error(err))
		return "", err
	}

	return result, nil
}

func (u *User) ChangeCurrentGym(ctx context.Context, userID, gymID int) (string, error) {
	result, err := u.userService.ChangeCurrentGym(ctx, userID, gymID)
	if err != nil {
		u.logger.Error("user.changeCurrentGym", zap.Error(err))
		return "", err
	}

	return result, nil
}
