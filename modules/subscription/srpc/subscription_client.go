package client

import (
	"context"
	subv1 "github.com/b3liv3r/protos-for-gym/gen/go/subscription"
	"github.com/b3liv3r/tgbot-for-gym/modules/subscription/models"
	"google.golang.org/grpc"
	"log"
	"time"
)

type RPCSubscriptioner interface {
	Create(ctx context.Context, userId int) (string, error)
	UpdateType(ctx context.Context, userId, Type, month int) (string, error)
	Extend(ctx context.Context, userId, month int) (string, error)
	GetData(ctx context.Context, userId int) (models.Subscription, error)
}

type SubscriptionClient struct {
	rpc  subv1.SubscriptionClient
	addr string
}

func NewSubscriptionClient(addr string) RPCSubscriptioner {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := subv1.NewSubscriptionClient(conn)

	return &SubscriptionClient{
		rpc:  c,
		addr: addr,
	}
}

func (s *SubscriptionClient) Create(ctx context.Context, userId int) (string, error) {
	resp, err := s.rpc.Create(ctx, &subv1.CreateRequest{UserId: int64(userId)})
	if err != nil {
		return "", err
	}
	return resp.Message, nil
}

func (s *SubscriptionClient) UpdateType(ctx context.Context, userId, Type, month int) (string, error) {
	resp, err := s.rpc.UpdateType(ctx, &subv1.UpdateTypeRequest{UserId: int64(userId), Type: int64(Type), Month: int64(month)})
	if err != nil {
		return "", err
	}
	return resp.Message, nil
}

func (s *SubscriptionClient) Extend(ctx context.Context, userId, month int) (string, error) {
	resp, err := s.rpc.Extend(ctx, &subv1.ExtendRequest{UserId: int64(userId), Count: int64(month)})
	if err != nil {
		return "", err
	}
	return resp.Message, nil
}

func (s *SubscriptionClient) GetData(ctx context.Context, userId int) (models.Subscription, error) {
	resp, err := s.rpc.GetData(ctx, &subv1.GetDataRequest{UserId: int64(userId)})
	if err != nil {
		return models.Subscription{}, err
	}
	startTime := time.Unix(resp.StartTime.GetSeconds(), int64(resp.StartTime.GetNanos())).UTC()
	endTime := time.Unix(resp.EndTime.GetSeconds(), int64(resp.EndTime.GetNanos())).UTC()
	return models.Subscription{
		SubscriptionType: int(resp.GetType()),
		StartDate:        startTime,
		EndDate:          endTime,
	}, nil
}
