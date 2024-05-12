package telegram

import (
	"context"
	"fmt"
	"github.com/b3liv3r/tgbot-for-gym/modules"
	umodels "github.com/b3liv3r/tgbot-for-gym/modules/users/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
	"strings"
)

const (
	HelpCmd                        = "/help"
	StartCmd                       = "/start"
	UserCreateCmd                  = "/u_create"
	UserProfileCmd                 = "/u_profile"
	UserUpdateCmd                  = "/u_update"
	UserChangeSubscriptionCmd      = "/u_change_subscription"
	UserChangeCurrentGymCmd        = "/u_change_current_gym"
	WalletDepositCmd               = "/w_deposit"
	WalletBalanceCmd               = "/w_balance"
	WalletHistoryCmd               = "/w_history"
	TrainerListForGymCmd           = "/t_list_for_gym"
	TrainerAvailableBookingListCmd = "/t_available_booking_list"
	TrainerCurrentBookingListCmd   = "/t_current_booking_list"
	TrainerBookingCmd              = "/t_booking"
	TrainerUnBookingCmd            = "/t_unbooking"
	SubscriptionExtendCmd          = "/s_extend"
	SubscriptionDataCmd            = "/s_data"
	GymListCmd                     = "/g_list"
	GymSchedulesCmd                = "/g_schedules"
	SupportCmd                     = "/support"
)

type TelegramRouter struct {
	Bot              *tgbotapi.BotAPI
	Services         modules.Services
	TelegramHandlers modules.Handlers
	log              *zap.Logger
}

func NewTelegramRouter(services modules.Services, bot *tgbotapi.BotAPI, logger *zap.Logger) *TelegramRouter {
	return &TelegramRouter{
		Bot:      bot,
		Services: services,
		log:      logger,
	}
}

func (r *TelegramRouter) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := r.Bot.GetUpdatesChan(u)
	if err != nil {
		r.log.Fatal("Error getting updates", zap.Error(err))
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch {
		case update.Message.IsCommand():
			r.handleCommand(update.Message)
		default:
			r.handleMessage(update.Message)

		}
	}
}

func (r *TelegramRouter) sendMessage(chatID int64, text string) {
	message := tgbotapi.NewMessage(chatID, text)
	_, err := r.Bot.Send(message)
	if err != nil {
		r.log.Error("Error sending message", zap.Error(err))
	}
}

func (r *TelegramRouter) handleMessage(msg *tgbotapi.Message) {
	r.sendMessage(msg.Chat.ID, msgUnknownCommand)
}

func (r *TelegramRouter) handleCommand(msg *tgbotapi.Message) {
	if len(msg.Text) == 0 {
		return
	}
	words := strings.Fields(msg.Text)
	command := words[0]
	args := words[1:]

	switch command {
	case HelpCmd:
		r.handleHelpCommand(msg)
	case StartCmd:
		r.handleStartCommand(msg)
	case UserCreateCmd:
		r.handleUserCreateCommand(msg, args)
	case UserProfileCmd:
		r.handleUserProfileCommand(msg)
	case UserUpdateCmd:
		r.handleUserUpdateCommand(msg, args)
	case UserChangeSubscriptionCmd:
		r.handleUserChangeSubscriptionCmd(msg, args)
	case UserChangeCurrentGymCmd:
		r.handleUserChangeCurrentGymCmd(msg, args)
	case WalletDepositCmd:
		r.handleWalletDepositCmd(msg, args)
	case WalletBalanceCmd:
		r.handleWalletBalanceCmd(msg)
	case WalletHistoryCmd:
		r.handleWalletHistoryCmd(msg, args)
	case SubscriptionExtendCmd:
		r.handleSubscriptionExtendCmd(msg, args)
	case SubscriptionDataCmd:
		r.handleSubscriptionDataCmd(msg)
	case GymListCmd:
		r.handleGymListCmd(msg)
	case GymSchedulesCmd:
		r.handleGymSchedulesCmd(msg, args)
	case TrainerListForGymCmd:
		r.handleTrainerListForGymCmd(msg, args)
	case TrainerAvailableBookingListCmd:
		r.handleTrainerAvailableBookingListCmd(msg, args)
	case TrainerCurrentBookingListCmd:
		r.handleTrainerCurrentBookingListCmd(msg)
	case TrainerBookingCmd:
		r.handleBookingCmd(msg, args)
	case TrainerUnBookingCmd:
		r.handleUnBookingCmd(msg, args)
	case SupportCmd:
		r.handleSupportCmd(msg)

	default:
		r.sendMessage(msg.Chat.ID, msgUnknownCommand)
	}
}

func (r *TelegramRouter) handleHelpCommand(msg *tgbotapi.Message) {
	r.sendMessage(msg.Chat.ID, msgHelp)
}

func (r *TelegramRouter) handleStartCommand(msg *tgbotapi.Message) {
	r.sendMessage(msg.Chat.ID, msgStart)
}

func (r *TelegramRouter) handleUserCreateCommand(msg *tgbotapi.Message, args []string) {
	if r.UserIsExist(msg.From.ID) {
		r.sendMessage(msg.Chat.ID, msgUserAlreadyExist)
		return
	}
	if len(args) != 4 {
		r.sendMessage(msg.Chat.ID, msgNotEnoughArguments)
		return
	}
	firstName := args[0]
	lastName := args[1]
	phone := args[2]
	email := args[3]
	username := firstName + " " + lastName
	err := isValidUsername(username)
	if err != nil {
		r.sendMessage(msg.Chat.ID, err.Error())
		return
	}
	err = isValidPhone(phone)
	if err != nil {
		r.sendMessage(msg.Chat.ID, err.Error())
		return
	}
	err = isValidEmail(email)
	if err != nil {
		r.sendMessage(msg.Chat.ID, err.Error())
		return
	}

	res, err := r.Services.Userer.Create(context.Background(), umodels.User{
		Id:       msg.From.ID,
		Username: username,
		Phone:    phone,
		Email:    email,
	})
	if err != nil {
		r.log.Error("Error creating user", zap.Error(err))
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		return
	}
	_, err = r.Services.Walleter.Create(context.Background(), msg.From.ID)
	if err != nil {
		r.log.Error("Error creating wallet", zap.Error(err))
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		return
	}
	_, err = r.Services.Subscriptioner.Create(context.Background(), msg.From.ID)
	if err != nil {
		r.log.Error("Error creating subscription", zap.Error(err))
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		return
	}

	r.sendMessage(msg.Chat.ID, res)
}

func (r *TelegramRouter) handleUserProfileCommand(msg *tgbotapi.Message) {
	user, err := r.Services.Userer.Profile(context.Background(), msg.From.ID)
	if err != nil {
		if err.Error() == "user not found" {
			r.sendMessage(msg.Chat.ID, msgUserNotFound)
			return
		}
		r.log.Error("Error getting user profile", zap.Error(err))
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		return
	}

	resp := fmt.Sprintf(msgProfileResp, user.Username, user.Phone, user.Email, user.SubscriptionLvl, user.CurrentGymId)
	r.sendMessage(msg.Chat.ID, resp)
}

func (r *TelegramRouter) handleUserUpdateCommand(msg *tgbotapi.Message, args []string) {
	if !r.UserIsExist(msg.From.ID) {
		r.sendMessage(msg.Chat.ID, msgUserNotFound)
		return
	}

	if len(args) != 4 {
		r.sendMessage(msg.Chat.ID, msgNotEnoughUpdateArguments)
		return
	}

	firstName := args[0]
	lastName := args[1]
	phone := args[2]
	email := args[3]
	username := firstName + " " + lastName
	var err error

	if isDephis(firstName) != isDephis(lastName) {
		r.sendMessage(msg.Chat.ID, msgUserUpdateError)
	} else if isDephis(firstName) && isDephis(lastName) {
		username = ""
	} else {
		err = isValidUsername(username)
		if err != nil {
			r.sendMessage(msg.Chat.ID, err.Error())
			return
		}
	}

	if isDephis(phone) {
		phone = ""
	} else {
		err = isValidPhone(phone)
		if err != nil {
			r.sendMessage(msg.Chat.ID, err.Error())
			return
		}
	}

	if isDephis(email) {
		email = ""
	} else {
		err = isValidEmail(email)
		if err != nil {
			r.sendMessage(msg.Chat.ID, err.Error())
			return
		}
	}

	res, err := r.Services.Userer.Update(context.Background(), umodels.User{
		Id:       msg.From.ID,
		Username: username,
		Phone:    phone,
		Email:    email,
	})
	if err != nil {
		r.log.Error("Error update user", zap.Error(err))
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		return
	}
	r.sendMessage(msg.Chat.ID, res)
}

func (r *TelegramRouter) handleUserChangeSubscriptionCmd(msg *tgbotapi.Message, args []string) {
	profile, err := r.Services.Userer.Profile(context.Background(), msg.From.ID)
	if err != nil {
		r.log.Error("Error getting user profile", zap.Error(err))
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		return
	}

	if len(args) != 1 {
		r.sendMessage(msg.Chat.ID, msgNotEnoughArguments)
		return
	}

	lvl, err := isSubLvl(args[0])
	if lvl == profile.SubscriptionLvl {
		r.sendMessage(msg.Chat.ID, msgLvlAlreadyGained)
		return
	}
	price := subLvlPrice(lvl)

	err = r.Payment(msg.From.ID, price, descChangeSubscription)
	if err != nil {
		return
	}

	res, err := r.Services.Subscriptioner.UpdateType(context.Background(), msg.From.ID, lvl, 1)
	if err != nil {
		r.log.Error("SUB.Error changing subscription", zap.Error(err))
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		err = r.Rollback(msg.From.ID, price, descChangeSubscription)
		return
	}

	res, err = r.Services.Userer.ChangeSubscription(context.Background(), msg.From.ID, lvl)
	if err != nil {
		r.log.Error("USR.Error changing subscription", zap.Error(err))
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		err = r.Rollback(msg.From.ID, price, descChangeSubscription)
		return
	}

	r.sendMessage(msg.Chat.ID, res)
}

func (r *TelegramRouter) handleUserChangeCurrentGymCmd(msg *tgbotapi.Message, args []string) {
	profile, err := r.Services.Userer.Profile(context.Background(), msg.From.ID)
	if err != nil {
		if err.Error() == "user not found" {
			r.sendMessage(msg.Chat.ID, msgUserNotFound)
			return
		}
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		return
	}

	if len(args) != 1 {
		r.sendMessage(msg.Chat.ID, msgNotEnoughArguments)
		return
	}

	gymID, err := isPositiveNumber(args[0])
	if err != nil {
		r.sendMessage(msg.Chat.ID, err.Error())
		return
	}

	res, err := r.Services.Gymer.List(context.Background())
	if err != nil {
		r.log.Error("Error getting user gym list", zap.Error(err))
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		return
	}
	lvl, stop := isFoundGym(res, gymID)
	if stop {
		r.sendMessage(msg.Chat.ID, msgUnknownGym)
		return
	}

	if profile.SubscriptionLvl < lvl {
		r.sendMessage(msg.Chat.ID, msgLowLvlSub)
		return
	}

	message, err := r.Services.Userer.ChangeCurrentGym(context.Background(), msg.From.ID, gymID)
	if err != nil {
		r.log.Error("Error changing current gym", zap.Error(err))
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		return
	}

	r.sendMessage(msg.Chat.ID, message)
}

func (r *TelegramRouter) handleWalletDepositCmd(msg *tgbotapi.Message, args []string) {
	if !r.UserIsExist(msg.From.ID) {
		r.sendMessage(msg.Chat.ID, msgUserNotFound)
		return
	}

	if len(args) != 1 {
		r.sendMessage(msg.Chat.ID, msgNotEnoughArguments)
		return
	}

	amount, err := isPositiveFloat(args[0])
	if err != nil {
		r.sendMessage(msg.Chat.ID, err.Error())
		return
	}

	err = r.Deposit(msg.From.ID, amount, descDeposit)
	if err != nil {
		r.log.Error("Error depositing wallet", zap.Error(err))
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		return
	}
}

func (r *TelegramRouter) handleWalletBalanceCmd(msg *tgbotapi.Message) {
	if !r.UserIsExist(msg.From.ID) {
		r.sendMessage(msg.Chat.ID, msgUserNotFound)
		return
	}

	res, err := r.Services.Walleter.GetBalance(context.Background(), msg.From.ID)
	if err != nil {
		r.log.Error("Error getting wallet balance", zap.Error(err))
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		return
	}
	resp := fmt.Sprintf(descCurrentBalance, res)
	r.sendMessage(msg.Chat.ID, resp)
}

func (r *TelegramRouter) handleWalletHistoryCmd(msg *tgbotapi.Message, args []string) {
	if !r.UserIsExist(msg.From.ID) {
		r.sendMessage(msg.Chat.ID, msgUserNotFound)
		return
	}

	if len(args) != 2 {
		r.sendMessage(msg.Chat.ID, msgNotEnoughArguments)
		return
	}

	startTime, err := parseTimeString(args[0])
	if err != nil {
		r.sendMessage(msg.Chat.ID, err.Error())
		return
	}
	endTime, err := parseTimeString(args[1])
	if err != nil {
		r.sendMessage(msg.Chat.ID, err.Error())
		return
	}
	res, err := r.Services.Walleter.History(context.Background(), msg.From.ID, "", startTime, endTime)
	if err != nil {
		r.log.Error("Error getting wallet history", zap.Error(err))
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		return
	}

	r.sendMessage(msg.Chat.ID, formatTransactionList(res))
}

func (r *TelegramRouter) handleSubscriptionExtendCmd(msg *tgbotapi.Message, args []string) {
	profile, err := r.Services.Userer.Profile(context.Background(), msg.From.ID)
	if err != nil {
		r.log.Error("Error getting user profile", zap.Error(err))
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		return
	}

	if profile.SubscriptionLvl == 0 {
		r.sendMessage(msg.Chat.ID, "You have no subscription level")
	}

	if len(args) != 1 {
		r.sendMessage(msg.Chat.ID, msgNotEnoughArguments)
		return
	}
	month, err := isPositiveNumber(args[0])
	if err != nil {
		r.sendMessage(msg.Chat.ID, err.Error())
		return
	}

	amount := subLvlPrice(profile.SubscriptionLvl) * float64(month)
	desc := fmt.Sprintf(descExtend, profile.SubscriptionLvl, month)
	err = r.Payment(msg.From.ID, amount, desc)
	if err != nil {
		return
	}

	res, err := r.Services.Subscriptioner.Extend(context.Background(), msg.From.ID, month)
	if err != nil {
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		r.Rollback(msg.From.ID, amount, desc)
		return
	}

	r.sendMessage(msg.Chat.ID, res)
}

func (r *TelegramRouter) handleSubscriptionDataCmd(msg *tgbotapi.Message) {
	if !r.UserIsExist(msg.From.ID) {
		r.sendMessage(msg.Chat.ID, msgUserNotFound)
		return
	}

	res, err := r.Services.Subscriptioner.GetData(context.Background(), msg.From.ID)
	if err != nil {
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		return
	}
	r.sendMessage(msg.Chat.ID, formatSubscriptionMessage(res))
}

func (r *TelegramRouter) handleGymListCmd(msg *tgbotapi.Message) {
	res, err := r.Services.Gymer.List(context.Background())
	if err != nil {
		r.log.Error("Error getting user gym list", zap.Error(err))
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		return
	}

	r.sendMessage(msg.Chat.ID, formatGymList(res))
}

func (r *TelegramRouter) handleGymSchedulesCmd(msg *tgbotapi.Message, args []string) {
	if len(args) != 1 {
		r.sendMessage(msg.Chat.ID, msgNotEnoughArguments)
		return
	}

	gymID, err := isPositiveNumber(args[0])
	if err != nil {
		r.sendMessage(msg.Chat.ID, err.Error())
		return
	}

	res, err := r.Services.Gymer.GetSchedules(context.Background(), gymID)
	if err != nil {
		r.log.Error("Error getting user gym list", zap.Error(err))
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		return
	}
	resp := formatSchedulesList(res)
	if resp == "" {
		r.sendMessage(msg.Chat.ID, msgUnknownGym)
		return
	}
	r.sendMessage(msg.Chat.ID, resp)
}

func (r *TelegramRouter) handleTrainerListForGymCmd(msg *tgbotapi.Message, args []string) {
	if len(args) != 1 {
		r.sendMessage(msg.Chat.ID, msgNotEnoughArguments)
		return
	}

	gymId, err := isPositiveNumber(args[0])
	if err != nil {
		r.sendMessage(msg.Chat.ID, err.Error())
		return
	}

	res, err := r.Services.Trainerer.ListForGym(context.Background(), gymId)
	if err != nil {
		r.log.Error("Error getting user gym list", zap.Error(err))
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		return
	}

	if len(res) == 0 {
		r.sendMessage(msg.Chat.ID, msgUnknownGym)
		return
	}

	r.sendMessage(msg.Chat.ID, formatTrainersList(res))
}

func (r *TelegramRouter) handleTrainerAvailableBookingListCmd(msg *tgbotapi.Message, args []string) {
	if len(args) != 1 {
		r.sendMessage(msg.Chat.ID, msgNotEnoughArguments)
		return
	}

	trainerId, err := isPositiveNumber(args[0])
	if err != nil {
		r.sendMessage(msg.Chat.ID, err.Error())
		return
	}

	res, err := r.Services.Trainerer.AvailableBookingList(context.Background(), trainerId)
	if err != nil {
		r.log.Error("Error getting user trainer list", zap.Error(err))
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		return
	}

	if len(res) == 0 {
		r.sendMessage(msg.Chat.ID, msgUnknownTrainerOrEmptyBookings)
		return
	}

	r.sendMessage(msg.Chat.ID, formatABookingList(res))
}

func (r *TelegramRouter) handleTrainerCurrentBookingListCmd(msg *tgbotapi.Message) {
	if !r.UserIsExist(msg.From.ID) {
		r.sendMessage(msg.Chat.ID, msgUserNotFound)
		return
	}

	res, err := r.Services.Trainerer.CurrentBookingList(context.Background(), msg.From.ID)
	if err != nil {
		r.log.Error("Error getting user trainer list", zap.Error(err))
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		return
	}

	if len(res) == 0 {
		r.sendMessage(msg.Chat.ID, msgCurrentBookingsEmpty)
	}

	r.sendMessage(msg.Chat.ID, formatABookingList(res))
}

func (r *TelegramRouter) handleBookingCmd(msg *tgbotapi.Message, args []string) {
	if !r.UserIsExist(msg.From.ID) {
		r.sendMessage(msg.Chat.ID, msgUserNotFound)
		return
	}

	if len(args) != 1 {
		r.sendMessage(msg.Chat.ID, msgNotEnoughArguments)
		return
	}

	bookingId, err := isPositiveNumber(args[0])
	if err != nil {
		r.sendMessage(msg.Chat.ID, err.Error())
		return
	}

	err = r.Payment(msg.From.ID, costDefaultTrainer, descBookTrainer)
	if err != nil {
		r.sendMessage(msg.Chat.ID, err.Error())
		return
	}

	res, err := r.Services.Trainerer.Booking(context.Background(), msg.From.ID, bookingId)
	if err != nil {
		r.Rollback(msg.From.ID, costDefaultTrainer, descBookTrainer)
		if err.Error() == msgPastSlotBooking || err.Error() == msgUnknownBooking {
			r.sendMessage(msg.Chat.ID, err.Error())
			return
		}
		r.log.Error("Error getting user trainer list", zap.Error(err))
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		return
	}

	r.sendMessage(msg.Chat.ID, res)
}

func (r *TelegramRouter) handleUnBookingCmd(msg *tgbotapi.Message, args []string) {
	if !r.UserIsExist(msg.From.ID) {
		r.sendMessage(msg.Chat.ID, msgUserNotFound)
		return
	}

	if len(args) != 1 {
		r.sendMessage(msg.Chat.ID, msgNotEnoughArguments)
		return
	}

	bookingId, err := isPositiveNumber(args[0])
	if err != nil {
		r.sendMessage(msg.Chat.ID, err.Error())
		return
	}

	err = r.Payment(msg.From.ID, -costDefaultTrainer, descUnBookTrainer)
	if err != nil {
		r.sendMessage(msg.Chat.ID, err.Error())
		return
	}

	res, err := r.Services.Trainerer.UnBooking(context.Background(), bookingId)
	if err != nil {
		r.Rollback(msg.From.ID, -costDefaultTrainer, descUnBookTrainer)
		if err.Error() == msgLateUnBooking || err.Error() == msgUnknownUnBooking {
			r.sendMessage(msg.Chat.ID, err.Error())
			return
		}
		r.log.Error("Error getting user trainer list", zap.Error(err))
		r.sendMessage(msg.Chat.ID, msgUnknownError)
		return
	}

	r.sendMessage(msg.Chat.ID, res)
}

func (r *TelegramRouter) handleSupportCmd(msg *tgbotapi.Message) {
	r.sendMessage(msg.Chat.ID, msgSupport)
}
