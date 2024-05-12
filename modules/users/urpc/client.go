package client

import (
	"context"
	userv1 "github.com/b3liv3r/protos-for-gym/gen/go/user"
	"github.com/b3liv3r/tgbot-for-gym/modules/users/models"
	"google.golang.org/grpc"
	"log"
)

type RPCUserer interface {
	Create(ctx context.Context, user models.User) (string, error)
	Profile(ctx context.Context, userID int) (models.User, error)
	Update(ctx context.Context, user models.User) (string, error)
	ChangeSubscription(ctx context.Context, userID, subLvl int) (string, error)
	ChangeCurrentGym(ctx context.Context, userID, gymID int) (string, error)
}

type UserClient struct {
	rpc  userv1.UserClient
	addr string
}

func NewUserClient(addr string) RPCUserer {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := userv1.NewUserClient(conn)

	return &UserClient{
		rpc:  c,
		addr: addr,
	}
}

func (u *UserClient) Create(ctx context.Context, user models.User) (string, error) {
	resp, err := u.rpc.Create(ctx, &userv1.CreateRequest{
		UserId:   int64(user.Id),
		Username: user.Username,
		Phone:    user.Phone,
		Email:    user.Email,
	})
	if err != nil {
		return "", err
	}

	return resp.GetMessage(), nil
}

func (u *UserClient) Profile(ctx context.Context, userID int) (models.User, error) {
	resp, err := u.rpc.Profile(ctx, &userv1.ProfileRequest{
		UserId: int64(userID),
	})
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		Username:        resp.Username,
		Phone:           resp.Phone,
		Email:           resp.Email,
		SubscriptionLvl: int(resp.SubscriptionLvl),
		CurrentGymId:    int(resp.CurrentGymId),
	}, nil
}

func (u *UserClient) Update(ctx context.Context, user models.User) (string, error) {
	resp, err := u.rpc.Update(ctx, &userv1.UpdateRequest{
		UserId:   int64(user.Id),
		Username: user.Username,
		Phone:    user.Phone,
		Email:    user.Email,
	})
	if err != nil {
		return "", err
	}

	return resp.GetMessage(), nil
}

func (u *UserClient) ChangeCurrentGym(ctx context.Context, userID, gymID int) (string, error) {
	resp, err := u.rpc.ChangeCurrentGym(ctx, &userv1.ChangeCurrentGymRequest{
		UserId:       int64(userID),
		CurrentGymId: int64(gymID),
	})
	if err != nil {
		return "", err
	}

	return resp.GetMessage(), nil
}

func (u *UserClient) ChangeSubscription(ctx context.Context, userID, subLvl int) (string, error) {
	resp, err := u.rpc.ChangeSubscriptions(ctx, &userv1.ChangeSubscriptionsRequest{
		UserId:          int64(userID),
		SubscriptionLvl: int64(subLvl),
	})
	if err != nil {
		return "", err
	}

	return resp.GetMessage(), nil
}
