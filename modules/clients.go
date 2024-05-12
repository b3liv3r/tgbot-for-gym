package modules

import (
	gclient "github.com/b3liv3r/tgbot-for-gym/modules/gym/grpc"
	sclient "github.com/b3liv3r/tgbot-for-gym/modules/subscription/srpc"
	tclient "github.com/b3liv3r/tgbot-for-gym/modules/trainers/trpc"
	uclient "github.com/b3liv3r/tgbot-for-gym/modules/users/urpc"
	wclient "github.com/b3liv3r/tgbot-for-gym/modules/wallet/wrpc"
	"os"
)

type Clients struct {
	wclient.RPCWalleter
	gclient.RPCGymer
	sclient.RPCSubscriptioner
	tclient.RPCTrainerer
	uclient.RPCUserer
}

func NewClients() Clients {
	return Clients{
		RPCWalleter:       wclient.NewWalletClient(os.Getenv("WALLET_CLIENT_ADDR")),
		RPCGymer:          gclient.NewGymClient(os.Getenv("GYM_CLIENT_ADDR")),
		RPCTrainerer:      tclient.NewTrainerClient(os.Getenv("TRAINER_CLIENT_ADDR")),
		RPCUserer:         uclient.NewUserClient(os.Getenv("USER_CLIENT_ADDR")),
		RPCSubscriptioner: sclient.NewSubscriptionClient(os.Getenv("SUBSCRIPTION_CLIENT_ADDR")),
	}
}
