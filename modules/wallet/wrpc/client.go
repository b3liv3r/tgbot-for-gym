package client

import (
	"context"
	walletv1 "github.com/b3liv3r/protos-for-gym/gen/go/wallet"
	"github.com/b3liv3r/tgbot-for-gym/modules/wallet/models"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"log"
	"time"
)

type RPCWalleter interface {
	Create(ctx context.Context, userID int) (string, error)
	GetBalance(ctx context.Context, userID int) (float64, error)
	Transaction(ctx context.Context, transaction models.Transaction) (string, float64, error)
	History(ctx context.Context, userID int, Type string, startTime, endTime time.Time) ([]models.Transaction, error)
}

type WalletClient struct {
	rpc  walletv1.WalletClient
	addr string
}

func NewWalletClient(addr string) RPCWalleter {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := walletv1.NewWalletClient(conn)

	return &WalletClient{
		rpc:  c,
		addr: addr,
	}
}

func (w *WalletClient) Create(ctx context.Context, userID int) (string, error) {
	resp, err := w.rpc.Create(ctx, &walletv1.CreateRequest{UserId: int64(userID)})
	if err != nil {
		return "", err
	}

	return resp.GetMessage(), nil
}

func (w *WalletClient) GetBalance(ctx context.Context, userID int) (float64, error) {
	resp, err := w.rpc.GetBalance(ctx, &walletv1.GetBalanceRequest{UserId: int64(userID)})
	if err != nil {
		return 0, err
	}
	return resp.GetBalance(), nil
}

func (w *WalletClient) Transaction(ctx context.Context, transaction models.Transaction) (string, float64, error) {
	resp, err := w.rpc.Transaction(ctx, &walletv1.TransactionRequest{
		UserId:      int64(transaction.UserId),
		Amount:      transaction.Amount,
		Type:        transaction.Type,
		Description: transaction.Description,
	})
	if err != nil {
		return "", 0, err
	}

	return resp.GetMessage(), resp.GetBalance(), nil
}

func (w *WalletClient) History(ctx context.Context, userID int, Type string, startTime, endTime time.Time) ([]models.Transaction, error) {
	startTimestamp, err := ptypes.TimestampProto(startTime)
	if err != nil {
		return nil, err
	}

	endTimestamp, err := ptypes.TimestampProto(endTime)
	if err != nil {
		return nil, err
	}

	resp, err := w.rpc.History(ctx, &walletv1.HistoryRequest{
		UserId:    int64(userID),
		Type:      Type,
		StartTime: startTimestamp,
		EndTime:   endTimestamp,
	})
	if err != nil {
		return nil, err
	}

	transactions := make([]models.Transaction, 0)
	for _, transaction := range resp.GetTransactions() {
		date := time.Unix(transaction.Timestamp.GetSeconds(), int64(transaction.Timestamp.GetNanos()))
		transactions = append(transactions, models.Transaction{
			Id:          int(transaction.Id),
			Amount:      transaction.Amount,
			Type:        transaction.Type,
			Description: transaction.Description,
			Date:        date,
		})
	}

	return transactions, nil
}
