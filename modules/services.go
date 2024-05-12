package modules

import (
	gservice "github.com/b3liv3r/tgbot-for-gym/modules/gym/service"
	sservice "github.com/b3liv3r/tgbot-for-gym/modules/subscription/service"
	tservice "github.com/b3liv3r/tgbot-for-gym/modules/trainers/service"
	uservice "github.com/b3liv3r/tgbot-for-gym/modules/users/service"
	wservice "github.com/b3liv3r/tgbot-for-gym/modules/wallet/service"
	"go.uber.org/zap"
)

type Services struct {
	Walleter       wservice.Walleter
	Gymer          gservice.Gymer
	Subscriptioner sservice.Subscriptioner
	Trainerer      tservice.Trainerer
	Userer         uservice.Userer
}

func NewServices(logger *zap.Logger, clients Clients) Services {
	return Services{
		Walleter:       wservice.NewWallet(logger, clients.RPCWalleter),
		Gymer:          gservice.NewGym(logger, clients.RPCGymer),
		Subscriptioner: sservice.NewSubscription(logger, clients.RPCSubscriptioner),
		Trainerer:      tservice.NewTrainer(logger, clients.RPCTrainerer),
		Userer:         uservice.NewUser(logger, clients.RPCUserer),
	}
}
